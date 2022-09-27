/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cloudprober

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ns1labs/orb/agent/otel/cloudprobereceiver"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-cmd/cmd"
	"github.com/go-co-op/gocron"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/agent/otel/otlpmqttexporter"
	"github.com/ns1labs/orb/agent/policies"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

var _ backend.Backend = (*cloudproberBackend)(nil)

const (
	BackendName         = "cloudprober"
	DefaultBinary       = "/usr/local/sbin/cloudprober"
	ReadinessBackoff    = 10
	ReadinessTimeout    = 10
	ApplyPolicyTimeout  = 10
	RemovePolicyTimeout = 15
	VersionTimeout      = 5
	ScrapeTimeout       = 5
	TapsTimeout         = 5
)

type cloudproberBackend struct {
	logger             *zap.Logger
	binary             string
	configFile         string
	cloudproberVersion string
	proc               *cmd.Cmd
	statusChan         <-chan cmd.Status
	startTime          time.Time
	cancelFunc         context.CancelFunc

	// MQTT Config for OTEL MQTT
	mqttConfig config.MQTTConfig

	mqttClient   mqtt.Client
	metricsTopic string
	otlpMetricsTopic string
	scraper      *gocron.Scheduler
	policyRepo   policies.PolicyRepo

	receiver component.MetricsReceiver
	exporter component.MetricsExporter

	adminAPIHost     string
	adminAPIPort     string
	adminAPIProtocol string

	scrapeOtel bool
}

func (c *cloudproberBackend) GetStartTime() time.Time {
	return c.startTime
}

func (c *cloudproberBackend) SetCommsClient(agentID string, client mqtt.Client, baseTopic string) {
	c.mqttClient = client
	c.metricsTopic = fmt.Sprintf("%s/m/%c", baseTopic, agentID[0])
}

func (c *cloudproberBackend) GetState() (backend.BackendState, string, error) {
	_, err := c.checkAlive()
	if err != nil {
		return backend.Unknown, "", err
	}
	return backend.Running, "", nil
}

// AppMetrics represents server application information
type AppMetrics struct {
	App struct {
		Version   string  `json:"version"`
		UpTimeMin float64 `json:"up_time_min"`
	} `json:"app"`
}

// note this needs to be stateless because it is called for multiple go routines
func (c *cloudproberBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string, timeout int32) error {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	alive, err := c.checkAlive()
	if !alive {
		return err
	}

	URL := fmt.Sprintf("%s://%s:%s/%s", c.adminAPIProtocol, c.adminAPIHost, c.adminAPIPort, url)

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		c.logger.Error("received error from payload", zap.Error(err))
		return err
	}

	req.Header.Add("Content-Type", contentType)
	res, getErr := client.Do(req)

	if getErr != nil {
		c.logger.Error("received error from payload", zap.Error(getErr))
		return getErr
	}

	if (res.StatusCode < 200) || (res.StatusCode > 299) {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("non 2xx HTTP error code from cloudprober, no or invalid body: %d", res.StatusCode))
		}
		if len(body) == 0 {
			return errors.New(fmt.Sprintf("%d empty body", res.StatusCode))
		} else if body[0] == '{' {
			var jsonBody map[string]interface{}
			err := json.Unmarshal(body, &jsonBody)
			if err == nil {
				if errMsg, ok := jsonBody["error"]; ok {
					return errors.New(fmt.Sprintf("%d %s", res.StatusCode, errMsg))
				}
			}
		}
	}

	if res.Body != nil {
		err = json.NewDecoder(res.Body).Decode(&payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *cloudproberBackend) checkAlive() (bool, error) {
	policyData, err := c.policyRepo.GetAllByBackend(BackendName)
	if err != nil {
		if err.Error() == "not found" {
			// if doesn't have any policies then returns true because is on standby
			return true, nil
		}
	}
	if len(policyData) > 0 && c.proc.Status().StopTs > 0 {
		return false, errors.New("cloudprobe contains policies but is stopped")
	}

	return true, nil
}

// ApplyPolicy in cloudprober, when we receive a new probe, we will stop the binary, recreate the config with all policies bundled and then
//
//	start the cloudprober again
func (c *cloudproberBackend) ApplyPolicy(_ policies.PolicyData, _ bool) error {

	c.logger.Debug("recreate config file based on policies")

	policiesData, err := c.policyRepo.GetAllByBackend(BackendName)
	if err != nil {
		c.logger.Error("Received error during retrieving cloudprober policies")
		return err
	}

	if len(policiesData) == 0 {
		c.logger.Warn("No policies to apply, waiting")
		return nil
	}

	fileContent, err := c.buildConfigFile(policiesData)
	if err != nil {
		c.logger.Error("error on building cloudprober file")
		return err
	}

	err = c.overrideConfigFile(fileContent, c.configFile)
	if err != nil {
		c.logger.Error("failure to override configuration")
		return err
	}
	err = c.stopBinary()
	err = c.startBinary()
	if err != nil {
		return err
	}

	return nil

}

// RemovePolicy in here we will only call the ApplyPolicy since the same concept applies, the policyManager will update the repo and re-generate the config file
func (c *cloudproberBackend) RemovePolicy(data policies.PolicyData) error {
	return c.ApplyPolicy(data, false)

}

func (c *cloudproberBackend) Version() (string, error) {
	var appMetrics AppMetrics
	err := c.request("metrics/app", &appMetrics, http.MethodGet, http.NoBody, "application/json", VersionTimeout)
	if err != nil {
		return "", err
	}
	c.cloudproberVersion = appMetrics.App.Version
	return appMetrics.App.Version, nil
}

func (c *cloudproberBackend) Write(payload []byte) (n int, err error) {
	c.logger.Info("cloudprober", zap.ByteString("log", payload))
	return len(payload), nil
}

func (c *cloudproberBackend) Start(ctx context.Context, cancelFunc context.CancelFunc) error {

	// this should record the start time whether it's successful or not
	// because it is used by the automatic restart system for last attempt
	c.startTime = time.Now()
	c.cancelFunc = cancelFunc

	// log STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		for c.proc.Stdout != nil || c.proc.Stderr != nil {
			select {
			case line, open := <-c.proc.Stdout:
				if !open {
					c.proc.Stdout = nil
					continue
				}
				c.logger.Info("cloudprober stdout", zap.String("log", line))
			case line, open := <-c.proc.Stderr:
				if !open {
					c.proc.Stderr = nil
					continue
				}
				c.logger.Info("cloudprober stderr", zap.String("log", line))
			}
		}
	}()

	c.logger.Info("cloudprober waiting for policies")
	c.scraper = gocron.NewScheduler(time.UTC)
	c.scraper.StartAsync()

	// only one scrape mechanism
	c.scrapeOpenTelemetry(ctx)

	err := c.ApplyPolicy(policies.PolicyData{}, false)
	if err != nil {
		c.logger.Error("error during applying policies")
		return err
	}

	return nil
}

func (c *cloudproberBackend) startBinary() error {
	_, err := exec.LookPath(c.binary)
	if err != nil {
		if _, err := exec.LookPath(DefaultBinary); err == nil {
			c.binary = DefaultBinary
		} else {
			c.logger.Error("cloudprober startup error: binary not found", zap.Error(err))
			return err
		}
	}

	pvOptions := []string{
		"-config_file",
	}
	if len(c.configFile) > 0 {
		pvOptions = append(pvOptions, c.configFile)
	}
	c.logger.Info("cloudprober startup", zap.Strings("arguments", pvOptions))

	c.proc = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, c.binary, pvOptions...)
	c.statusChan = c.proc.Start()

	// wait for simple startup errors
	time.Sleep(time.Second)

	status := c.proc.Status()

	if status.Error != nil {
		c.logger.Error("cloudprober startup error", zap.Error(status.Error))
		return status.Error
	}

	if status.Complete {
		err := c.stopBinary()
		if err != nil {
			c.logger.Error("proc.Stop error", zap.Error(err))
		}
		return errors.New("cloudprober startup error, check log")
	}
	return err
}

func (c *cloudproberBackend) stopBinary() error {
	err := c.proc.Stop()
	if err != nil {
		return err
	}
	return nil
}

func (c *cloudproberBackend) scrapeOpenTelemetry(ctx context.Context) {
	go func() {
		startExpCtx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()
		var ok bool
		var err error
		for i := 1; i < 10; i++ {
			select {
			case <-startExpCtx.Done():
				return
			default:
				if c.mqttClient != nil {
					var errStartExp error
					c.exporter, errStartExp = c.createOtlpMqttExporter(ctx)
					if errStartExp != nil {
						c.logger.Error("failed to create a exporter", zap.Error(err))
						return
					}

					c.receiver, err = createReceiver(ctx, c.exporter, c.logger)
					if err != nil {
						c.logger.Error("failed to create a receiver", zap.Error(err))
						return
					}

					err = c.exporter.Start(ctx, nil)
					if err != nil {
						c.logger.Error("otel mqtt exporter startup error", zap.Error(err))
						return
					}

					err = c.receiver.Start(ctx, nil)
					if err != nil {
						c.logger.Error("otel receiver startup error", zap.Error(err))
						return
					}

					ok = true
					return
				} else {
					c.logger.Info("waiting until mqtt client is connected", zap.String("wait time", (time.Duration(i)*time.Second).String()))
					time.Sleep(time.Duration(i) * time.Second)
					continue
				}
			}
		}
		if !ok {
			c.logger.Error("mqtt did not established a connection, stopping agent")
			err := c.Stop(startExpCtx)
			if err != nil {
				return
			}
		}
		return
	}()
}

func (c *cloudproberBackend) Stop(ctx context.Context) error {
	c.logger.Info("routine call to stop cloudprober", zap.Any("routine", ctx.Value("routine")))
	defer c.cancelFunc()
	err := c.stopBinary()
	finalStatus := <-c.statusChan
	if err != nil {
		c.logger.Error("cloudprober shutdown error", zap.Error(err))
	}
	c.scraper.Stop()

	if c.scrapeOtel {
		if c.exporter != nil {
			_ = c.exporter.Shutdown(context.Background())
		}
		if c.receiver != nil {
			_ = c.receiver.Shutdown(context.Background())
		}
	}

	c.logger.Info("cloudprober process stopped", zap.Int("pid", finalStatus.PID), zap.Int("exit_code", finalStatus.Exit))
	return nil
}

func (c *cloudproberBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo, config map[string]string, otelConfig map[string]interface{}) error {
	c.logger = logger
	c.policyRepo = repo

	var prs bool
	if c.binary, prs = config["binary"]; !prs {
		return errors.New("you must specify cloudprober binary")
	}
	if c.configFile, prs = config["config_file"]; !prs {
		c.configFile = ""
	}
	if c.adminAPIHost, prs = config["api_host"]; !prs {
		return errors.New("you must specify cloudprober admin API host")
	}
	if c.adminAPIPort, prs = config["api_port"]; !prs {
		return errors.New("you must specify cloudprober admin API port")
	}

	for k, v := range otelConfig {
		switch k {
		case "Enable":
			c.scrapeOtel = v.(bool)
		}
	}

	return nil
}

func (c *cloudproberBackend) scrapeMetrics() (map[string]interface{}, error) {
	var metrics map[string]interface{}
	err := c.request("metrics", &metrics, http.MethodGet, http.NoBody, "application/json", ScrapeTimeout)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (c *cloudproberBackend) GetCapabilities() (map[string]interface{}, error) {
	var taps interface{}
	err := c.request("taps", &taps, http.MethodGet, http.NoBody, "application/json", TapsTimeout)
	if err != nil {
		return nil, err
	}
	jsonBody := make(map[string]interface{})
	jsonBody["taps"] = taps
	return jsonBody, nil
}

func Register() bool {
	backend.Register("cloudprober", &cloudproberBackend{
		adminAPIHost:     "localhost", // using default by their doc
		adminAPIPort:     "9313",
		adminAPIProtocol: "http",
	})
	//p.logger.Error("trying to Register backend", zap.Error(err))
	return true
}

func (c *cloudproberBackend) createOtlpMqttExporter(ctx context.Context) (component.MetricsExporter, error) {

	if c.mqttClient != nil {
		cfg := otlpmqttexporter.CreateConfigClient(c.mqttClient, c.metricsTopic, c.cloudproberVersion)
		set := otlpmqttexporter.CreateDefaultSettings(c.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	} else {
		cfg := otlpmqttexporter.CreateConfig(c.mqttConfig.Address, c.mqttConfig.Id, c.mqttConfig.Key, c.mqttConfig.ChannelID, c.cloudproberVersion, c.otlpMetricsTopic)
		set := otlpmqttexporter.CreateDefaultSettings(c.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	}

}

func createReceiver(ctx context.Context, exporter component.MetricsExporter, logger *zap.Logger) (component.MetricsReceiver, error) {
	set := cloudprobereceiver.CreateDefaultSettings(logger)
	cfg := cloudprobereceiver.CreateDefaultConfig()
	// Create the Prometheus receiver and pass in the previously created Prometheus exporter.
	receiver, err := cloudprobereceiver.CreateMetricsReceiver(ctx, set, cfg, exporter)
	if err != nil {
		return nil, err
	}
	return receiver, nil
}

func (c *cloudproberBackend) FullReset(ctx context.Context) error {

	// force a stop, which stops scrape as well. if proc is dead, it no ops.
	if state, _, _ := c.GetState(); state == backend.Running {
		if err := c.Stop(ctx); err != nil {
			c.logger.Error("failed to stop backend on restart procedure", zap.Error(err))
			return err
		}
	}

	backendCtx, cancelFunc := context.WithCancel(context.WithValue(ctx, "routine", "cloudprober"))
	// start it
	if err := c.Start(backendCtx, cancelFunc); err != nil {
		c.logger.Error("failed to start backend on restart procedure", zap.Error(err))
		return err
	}

	return nil
}

func (c *cloudproberBackend) buildConfigFile(policyYaml []policies.PolicyData) ([]byte, error) {
	//hardcoded for now, we will use the proto from cloudprober dependencies to parse policy data
	hardCodedConfigs := "probe {\n  name: \"google_homepage\"\n  type: HTTP\n  targets {\n    host_names: \"www.google.com\"\n  }\n  interval_msec: 5000 \n  timeout_msec: 1000 \n}\nserver {\n  type: HTTP\n  http_server {\n    port: 8099\n  }\n}\n\nsurfacer {\n  type: PROMETHEUS\n\n  prometheus_surfacer {\n    # Following option adds a prefix to exported metrics, for example,\n    # \"total\" metric is exported as \"cloudprober_total\".\n    metrics_prefix: \"cloudprober_\"\n  }\n}"
	return []byte(hardCodedConfigs), nil
}

func (c *cloudproberBackend) overrideConfigFile(content []byte, location string) error {
	err := os.Remove(location)
	if err != nil {
		return err
	}

	err = os.WriteFile(location, content, os.ModeType)
	if err != nil {
		return err
	}

	return nil
}
