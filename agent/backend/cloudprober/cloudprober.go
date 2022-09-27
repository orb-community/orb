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
	"gopkg.in/yaml.v3"
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
	scraper      *gocron.Scheduler
	policyRepo   policies.PolicyRepo

	receiver component.MetricsReceiver
	exporter component.MetricsExporter

	adminAPIHost     string
	adminAPIPort     string
	adminAPIProtocol string

	scrapeOtel bool
}

func (cpbe *cloudproberBackend) GetStartTime() time.Time {
	return cpbe.startTime
}

func (cpbe *cloudproberBackend) SetCommsClient(agentID string, client mqtt.Client, baseTopic string) {
	cpbe.mqttClient = client
	cpbe.metricsTopic = fmt.Sprintf("%s/m/%c", baseTopic, agentID[0])
}

func (cpbe *cloudproberBackend) GetState() (backend.BackendState, string, error) {
	_, err := cpbe.checkAlive()
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
func (cpbe *cloudproberBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string, timeout int32) error {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	alive, err := cpbe.checkAlive()
	if !alive {
		return err
	}

	URL := fmt.Sprintf("%s://%s:%s/%s", cpbe.adminAPIProtocol, cpbe.adminAPIHost, cpbe.adminAPIPort, url)

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		cpbe.logger.Error("received error from payload", zap.Error(err))
		return err
	}

	req.Header.Add("Content-Type", contentType)
	res, getErr := client.Do(req)

	if getErr != nil {
		cpbe.logger.Error("received error from payload", zap.Error(getErr))
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

func (cpbe *cloudproberBackend) checkAlive() (bool, error) {
	policyData, err := cpbe.policyRepo.GetAllByBackend(BackendName)
	if err != nil {
		if err.Error() == "not found" {
			// if doesn't have any policies then returns true because is on standby
			return true, nil
		}
	}
	if len(policyData) > 0 && cpbe.proc.Status().StopTs > 0 {
		return false, errors.New("cloudprobe contains policies but is stopped")
	}

	return true, nil
}

// ApplyPolicy in cloudprober, when we receive a new probe, we will stop the binary, recreate the config with all policies bundled and then
//
//	start the cloudprober again
func (cpbe *cloudproberBackend) ApplyPolicy(data policies.PolicyData, _ bool) error {

	cpbe.logger.Debug("recreate config file based on policies")

	cpbe.logger.Debug("cloudprober policy apply", zap.String("policy_id", data.ID), zap.Any("data", data.Data))

	fullPolicy := map[string]interface{}{
		"version": "1.0",
		"prober": map[string]interface{}{
			"policies": map[string]interface{}{
				data.Name: data.Data,
			},
		},
	}

	policyYaml, err := yaml.Marshal(fullPolicy)
	if err != nil {
		cpbe.logger.Warn("yaml policy marshal failure", zap.String("policy_id", data.ID), zap.Any("policy", fullPolicy))
		return err
	}

	defaultLocation := "/etc/cloudprober.cfg"
	fileContent, err := cpbe.buildConfigFile(policyYaml)
	if err != nil {
		cpbe.logger.Warn("policy convertion to cloudprober configuration failure", zap.String("policy_id", data.ID), zap.Any("policy", fullPolicy))
		return err
	}
	err = cpbe.overrideConfigFile(fileContent, defaultLocation)
	if err != nil {
		cpbe.logger.Error("failure to override configuration", zap.String("policy_id", data.ID), zap.Any("policy", fullPolicy))
		return err
	}
	err = cpbe.stopBinary()
	err = cpbe.startBinary()
	if err != nil {
		return err
	}

	return nil

}

func (cpbe *cloudproberBackend) RemovePolicy(data policies.PolicyData) error {
	var resp interface{}
	err := cpbe.request(fmt.Sprintf("policies/%s", data.Name), &resp, http.MethodDelete, http.NoBody, "application/json", RemovePolicyTimeout)
	if err != nil {
		cpbe.logger.Error("received error", zap.Error(err))
		return err
	}
	return nil
}

func (cpbe *cloudproberBackend) Version() (string, error) {
	var appMetrics AppMetrics
	err := cpbe.request("metrics/app", &appMetrics, http.MethodGet, http.NoBody, "application/json", VersionTimeout)
	if err != nil {
		return "", err
	}
	cpbe.cloudproberVersion = appMetrics.App.Version
	return appMetrics.App.Version, nil
}

func (cpbe *cloudproberBackend) Write(payload []byte) (n int, err error) {
	cpbe.logger.Info("cloudprober", zap.ByteString("log", payload))
	return len(payload), nil
}

func (cpbe *cloudproberBackend) Start(ctx context.Context, cancelFunc context.CancelFunc) error {

	// this should record the start time whether it's successful or not
	// because it is used by the automatic restart system for last attempt
	cpbe.startTime = time.Now()
	cpbe.cancelFunc = cancelFunc

	// log STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		for cpbe.proc.Stdout != nil || cpbe.proc.Stderr != nil {
			select {
			case line, open := <-cpbe.proc.Stdout:
				if !open {
					cpbe.proc.Stdout = nil
					continue
				}
				cpbe.logger.Info("cloudprober stdout", zap.String("log", line))
			case line, open := <-cpbe.proc.Stderr:
				if !open {
					cpbe.proc.Stderr = nil
					continue
				}
				cpbe.logger.Info("cloudprober stderr", zap.String("log", line))
			}
		}
	}()

	cpbe.logger.Info("cloudprober process started", zap.Int("pid", status.PID))

	cpbe.scraper = gocron.NewScheduler(time.UTC)
	cpbe.scraper.StartAsync()

	// only one scrape mechanism
	cpbe.scrapeOpenTelemetry(ctx)

	return nil
}

func (cpbe *cloudproberBackend) startBinary() error {
	_, err := exec.LookPath(cpbe.binary)
	if err != nil {
		cpbe.logger.Error("cloudprober startup error: binary not found", zap.Error(err))
		return err
	}

	pvOptions := []string{
		"-config_file",
	}
	if len(cpbe.configFile) > 0 {
		pvOptions = append(pvOptions, cpbe.configFile)
	}
	cpbe.logger.Info("cloudprober startup", zap.Strings("arguments", pvOptions))

	cpbe.proc = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, cpbe.binary, pvOptions...)
	cpbe.statusChan = cpbe.proc.Start()

	// wait for simple startup errors
	time.Sleep(time.Second)

	status := cpbe.proc.Status()

	if status.Error != nil {
		cpbe.logger.Error("cloudprober startup error", zap.Error(status.Error))
		return status.Error
	}

	if status.Complete {
		err := cpbe.stopBinary()
		if err != nil {
			cpbe.logger.Error("proc.Stop error", zap.Error(err))
		}
		return errors.New("cloudprober startup error, check log")
	}
	return err
}

func (cpbe *cloudproberBackend) stopBinary() error {
	err := cpbe.proc.Stop()
	if err != nil {
		return err
	}
	return nil
}

func (cpbe *cloudproberBackend) scrapeOpenTelemetry(ctx context.Context) {
	go func() {
		startExpCtx, cancelFunc := context.WithCancel(ctx)
		var ok bool
		var err error
		for i := 1; i < 10; i++ {
			select {
			case <-startExpCtx.Done():
				return
			default:
				if cpbe.mqttClient != nil {
					var errStartExp error
					cpbe.exporter, errStartExp = cpbe.createOtlpMqttExporter(ctx)
					if errStartExp != nil {
						cpbe.logger.Error("failed to create a exporter", zap.Error(err))
						return
					}

					cpbe.receiver, err = createReceiver(ctx, cpbe.exporter, cpbe.logger)
					if err != nil {
						cpbe.logger.Error("failed to create a receiver", zap.Error(err))
						return
					}

					err = cpbe.exporter.Start(ctx, nil)
					if err != nil {
						cpbe.logger.Error("otel mqtt exporter startup error", zap.Error(err))
						return
					}

					err = cpbe.receiver.Start(ctx, nil)
					if err != nil {
						cpbe.logger.Error("otel receiver startup error", zap.Error(err))
						return
					}

					ok = true
					return
				} else {
					cpbe.logger.Info("waiting until mqtt client is connected", zap.String("wait time", (time.Duration(i)*time.Second).String()))
					time.Sleep(time.Duration(i) * time.Second)
					continue
				}
			}
		}
		if !ok {
			cpbe.logger.Error("mqtt did not established a connection, stopping agent")
			cpbe.Stop(startExpCtx)
		}
		cancelFunc()
		return
	}()
}

func (cpbe *cloudproberBackend) Stop(ctx context.Context) error {
	cpbe.logger.Info("routine call to stop cloudprober", zap.Any("routine", ctx.Value("routine")))
	defer cpbe.cancelFunc()
	err := cpbe.stopBinary()
	finalStatus := <-cpbe.statusChan
	if err != nil {
		cpbe.logger.Error("cloudprober shutdown error", zap.Error(err))
	}
	cpbe.scraper.Stop()

	if cpbe.scrapeOtel {
		if cpbe.exporter != nil {
			_ = cpbe.exporter.Shutdown(context.Background())
		}
		if cpbe.receiver != nil {
			_ = cpbe.receiver.Shutdown(context.Background())
		}
	}

	cpbe.logger.Info("cloudprober process stopped", zap.Int("pid", finalStatus.PID), zap.Int("exit_code", finalStatus.Exit))
	return nil
}

func (cpbe *cloudproberBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo, config map[string]string, otelConfig map[string]interface{}) error {
	cpbe.logger = logger
	cpbe.policyRepo = repo

	var prs bool
	if cpbe.binary, prs = config["binary"]; !prs {
		return errors.New("you must specify cloudprober binary")
	}
	if cpbe.configFile, prs = config["config_file"]; !prs {
		cpbe.configFile = ""
	}
	if cpbe.adminAPIHost, prs = config["api_host"]; !prs {
		return errors.New("you must specify cloudprober admin API host")
	}
	if cpbe.adminAPIPort, prs = config["api_port"]; !prs {
		return errors.New("you must specify cloudprober admin API port")
	}

	for k, v := range otelConfig {
		switch k {
		case "Enable":
			cpbe.scrapeOtel = v.(bool)
		}
	}

	return nil
}

func (cpbe *cloudproberBackend) scrapeMetrics() (map[string]interface{}, error) {
	var metrics map[string]interface{}
	err := cpbe.request("metrics", &metrics, http.MethodGet, http.NoBody, "application/json", ScrapeTimeout)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (cpbe *cloudproberBackend) GetCapabilities() (map[string]interface{}, error) {
	var taps interface{}
	err := cpbe.request("taps", &taps, http.MethodGet, http.NoBody, "application/json", TapsTimeout)
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

func (cpbe *cloudproberBackend) createOtlpMqttExporter(ctx context.Context) (component.MetricsExporter, error) {

	if cpbe.mqttClient != nil {
		cfg := otlpmqttexporter.CreateConfigClient(cpbe.mqttClient, cpbe.metricsTopic, cpbe.cloudproberVersion)
		set := otlpmqttexporter.CreateDefaultSettings(cpbe.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	} else {
		cfg := otlpmqttexporter.CreateConfig(cpbe.mqttConfig.Address, cpbe.mqttConfig.Id, cpbe.mqttConfig.Key, cpbe.mqttConfig.ChannelID, cpbe.cloudproberVersion)
		set := otlpmqttexporter.CreateDefaultSettings(cpbe.logger)
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

func (cpbe *cloudproberBackend) FullReset(ctx context.Context) error {

	// force a stop, which stops scrape as well. if proc is dead, it no ops.
	if state, _, _ := cpbe.GetState(); state == backend.Running {
		if err := cpbe.Stop(ctx); err != nil {
			cpbe.logger.Error("failed to stop backend on restart procedure", zap.Error(err))
			return err
		}
	}

	backendCtx, cancelFunc := context.WithCancel(context.WithValue(ctx, "routine", "cloudprober"))
	// start it
	if err := cpbe.Start(backendCtx, cancelFunc); err != nil {
		cpbe.logger.Error("failed to start backend on restart procedure", zap.Error(err))
		return err
	}

	return nil
}

func (cpbe *cloudproberBackend) buildConfigFile(policyYaml []byte) ([]byte, error) {

	return nil, nil
}

func (cpbe *cloudproberBackend) overrideConfigFile(content []byte, location string) error {
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
