/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-cmd/cmd"
	"github.com/go-co-op/gocron"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/otel/otlpexporter"
	"github.com/ns1labs/orb/agent/otel/pktvisorreceiver"
	"github.com/ns1labs/orb/agent/policies"
	"github.com/ns1labs/orb/fleet"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

const DefaultBinary = "/usr/local/sbin/pktvisord"
const ReadinessBackoff = 10
const ReadinessTimeout = 10

type pktvisorBackend struct {
	logger          *zap.Logger
	binary          string
	configFile      string
	pktvisorVersion string
	proc            *cmd.Cmd
	statusChan      <-chan cmd.Status
	startTime       time.Time

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

func (p *pktvisorBackend) GetStartTime() time.Time {
	return p.startTime
}

func (p *pktvisorBackend) SetCommsClient(agentID string, client mqtt.Client, baseTopic string) {
	p.mqttClient = client
	p.metricsTopic = fmt.Sprintf("%s/m/%c", baseTopic, agentID[0])
}

func (p *pktvisorBackend) GetState() (backend.BackendState, string, error) {
	_, err := p.checkAlive()
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

// note this needs to be stateless because it is calledfor multiple go routines
func (p *pktvisorBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string, timeout int32) error {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	alive, err := p.checkAlive()
	if !alive {
		return err
	}

	URL := fmt.Sprintf("%s://%s:%s/api/v1/%s", p.adminAPIProtocol, p.adminAPIHost, p.adminAPIPort, url)

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return err
	}
	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Add("Content-Type", contentType)

	res, getErr := client.Do(req)

	if getErr != nil {
		return getErr
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("non 200 HTTP error code from pktvisord, no or invalid body: %d", res.StatusCode))
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

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return err
	}
	return nil
}

func (p *pktvisorBackend) checkAlive() (bool, error) {
	status := p.proc.Status()

	if status.Error != nil {
		p.logger.Error("pktvisor process error", zap.Error(status.Error))
		return false, status.Error
	}

	if status.Complete {
		err := p.proc.Stop()
		if err != nil {
			p.logger.Error("proc.Stop error", zap.Error(err))
		}
		return false, errors.New("pktvisor process ended")
	}

	return true, nil
}

func (p *pktvisorBackend) ApplyPolicy(data policies.PolicyData, updatePolicy bool) error {

	if updatePolicy {
		// To update a policy it's necessary first remove it and then apply a new version
		err := p.RemovePolicy(data)
		if err != nil {
			p.logger.Warn("policy failed to remove", zap.String("policy_id", data.ID), zap.String("policy_name", data.Name), zap.Error(err))
		}
	}

	p.logger.Debug("pktvisor policy apply", zap.String("policy_id", data.ID), zap.Any("data", data.Data))

	fullPolicy := map[string]interface{}{
		"version": "1.0",
		"visor": map[string]interface{}{
			"policies": map[string]interface{}{
				data.Name: data.Data,
			},
		},
	}

	pyaml, err := yaml.Marshal(fullPolicy)
	if err != nil {
		p.logger.Warn("yaml policy marshal failure", zap.String("policy_id", data.ID), zap.Any("policy", fullPolicy))
		return err
	}

	var resp map[string]interface{}
	err = p.request("policies", &resp, http.MethodPost, bytes.NewBuffer(pyaml), "application/x-yaml", 5)
	if err != nil {
		p.logger.Warn("yaml policy application failure", zap.String("policy_id", data.ID), zap.ByteString("policy", pyaml))
		return err
	}

	return nil

}

func (p *pktvisorBackend) RemovePolicy(data policies.PolicyData) error {
	var resp interface{}
	err := p.request(fmt.Sprintf("policies/%s", data.Name), &resp, http.MethodDelete, nil, "", 10)
	if err != nil {
		return err
	}
	return nil
}

func (p *pktvisorBackend) Version() (string, error) {
	var appMetrics AppMetrics
	err := p.request("metrics/app", &appMetrics, http.MethodGet, nil, "", 5)
	if err != nil {
		return "", err
	}
	p.pktvisorVersion = appMetrics.App.Version
	return appMetrics.App.Version, nil
}

func (p *pktvisorBackend) Write(payload []byte) (n int, err error) {
	p.logger.Info("pktvisord", zap.ByteString("log", payload))
	return len(payload), nil
}

func (p *pktvisorBackend) Start() error {

	// this should record the start time whether it's successful or not
	// because it is used by the automatic restart system for last attempt
	p.startTime = time.Now()

	_, err := exec.LookPath(p.binary)
	if err != nil {
		p.logger.Error("pktvisor startup error: binary not found", zap.Error(err))
		return err
	}

	pvOptions := []string{
		"--admin-api",
		"-l",
		p.adminAPIHost,
		"-p",
		p.adminAPIPort,
	}
	if len(p.configFile) > 0 {
		pvOptions = append(pvOptions, "--config", p.configFile)
	}
	p.logger.Info("pktvisor startup", zap.Strings("arguments", pvOptions))
	p.proc = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, p.binary, pvOptions...)
	p.statusChan = p.proc.Start()

	// log STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		for p.proc.Stdout != nil || p.proc.Stderr != nil {
			select {
			case line, open := <-p.proc.Stdout:
				if !open {
					p.proc.Stdout = nil
					continue
				}
				p.logger.Info("pktvisor stdout", zap.String("log", line))
			case line, open := <-p.proc.Stderr:
				if !open {
					p.proc.Stderr = nil
					continue
				}
				p.logger.Info("pktvisor stderr", zap.String("log", line))
			}
		}
	}()

	// wait for simple startup errors
	time.Sleep(time.Second)

	status := p.proc.Status()

	if status.Error != nil {
		p.logger.Error("pktvisor startup error", zap.Error(status.Error))
		return status.Error
	}

	if status.Complete {
		err = p.proc.Stop()
		if err != nil {
			p.logger.Error("proc.Stop error", zap.Error(err))
		}
		return errors.New("pktvisor startup error, check log")
	}

	p.logger.Info("pktvisor process started", zap.Int("pid", status.PID))

	var readinessError error
	for backoff := 0; backoff < ReadinessBackoff; backoff++ {
		var appMetrics AppMetrics
		readinessError = p.request("metrics/app", &appMetrics, http.MethodGet, nil, "", ReadinessTimeout)
		if readinessError == nil {
			p.logger.Info(fmt.Sprintf("pktvisor readiness ok, got version %s", appMetrics.App.Version))
			break
		}
		p.logger.Info("pktvisor is not ready, trying again")
	}

	if readinessError != nil {
		p.logger.Error("pktvisor error on readiness", zap.Error(readinessError))
		err = p.proc.Stop()
		if err != nil {
			p.logger.Error("proc.Stop error", zap.Error(err))
		}
		return readinessError
	}

	p.scraper = gocron.NewScheduler(time.UTC)
	p.scraper.StartAsync()

	if p.scrapeOtel {
		if err := p.scrapeOpentelemetry(); err != nil {
			return err
		}
	} else {
		if err := p.scrapeDefault(); err != nil {
			return err
		}
	}

	return nil
}

func (p *pktvisorBackend) scrapeDefault() error {
	// scrape all policy json output with one call every minute.
	// TODO support policies with custom bucket times
	job, err := p.scraper.Every(1).Minute().WaitForSchedule().Do(func() {
		metrics, err := p.scrapeMetrics(1)
		if err != nil {
			p.logger.Error("scrape failed", zap.Error(err))
			return
		}
		if len(metrics) == 0 {
			p.logger.Warn("scrape: no policies found, skipping")
			return
		}

		var batchPayload []fleet.AgentMetricsRPCPayload
		totalSize := 0
		for pName, pMetrics := range metrics {
			data, err := p.policyRepo.GetByName(pName)
			if err != nil {
				p.logger.Warn("skipping pktvisor policy not managed by orb", zap.String("policy", pName), zap.String("error message", err.Error()))
				continue
			}
			payloadData, err := json.Marshal(pMetrics)
			if err != nil {
				p.logger.Error("error marshalling scraped metric json", zap.String("policy", pName), zap.Error(err))
				continue
			}
			metricPayload := fleet.AgentMetricsRPCPayload{
				PolicyID:   data.ID,
				PolicyName: data.Name,
				Datasets:   data.GetDatasetIDs(),
				Format:     "json",
				BEVersion:  p.pktvisorVersion,
				Data:       payloadData,
			}
			batchPayload = append(batchPayload, metricPayload)
			totalSize += len(payloadData)
			p.logger.Info("scraped metrics for policy", zap.String("policy", pName), zap.String("policy_id", data.ID), zap.Int("payload_size_b", len(payloadData)))
		}

		rpc := fleet.AgentMetricsRPC{
			SchemaVersion: fleet.CurrentRPCSchemaVersion,
			Func:          fleet.AgentMetricsRPCFunc,
			Payload:       batchPayload,
		}

		body, err := json.Marshal(rpc)
		if err != nil {
			p.logger.Error("error marshalling metric rpc payload", zap.Error(err))
			return
		}

		if token := p.mqttClient.Publish(p.metricsTopic, 1, false, body); token.Wait() && token.Error() != nil {
			p.logger.Error("error sending metrics RPC", zap.String("topic", p.metricsTopic), zap.Error(token.Error()))
			return
		}
		p.logger.Info("scraped and published metrics", zap.String("topic", p.metricsTopic), zap.Int("payload_size_b", totalSize), zap.Int("batch_count", len(batchPayload)))

	})

	if err != nil {
		return err
	}

	job.SingletonMode()
	return nil
}

func (p *pktvisorBackend) scrapeOpentelemetry() (err error) {
	ctx := context.Background()
	p.exporter, err = createOtlpExporter(ctx, p.logger)
	if err != nil {
		p.logger.Error("failed to create a exporter", zap.Error(err))
	}
	p.receiver, err = createReceiver(ctx, p.exporter, p.logger)
	if err != nil {
		p.logger.Error("failed to create a receiver", zap.Error(err))
	}

	err = p.exporter.Start(ctx, nil)
	if err != nil {
		p.logger.Error("otel exporter startup error", zap.Error(err))
		return err
	}

	err = p.receiver.Start(ctx, nil)
	if err != nil {
		p.logger.Error("otel receiver startup error", zap.Error(err))
		return err
	}
	return nil
}

func (p *pktvisorBackend) Stop() error {
	p.logger.Info("pktvisor stopping")
	err := p.proc.Stop()
	finalStatus := <-p.statusChan
	if err != nil {
		p.logger.Error("pktvisor shutdown error", zap.Error(err))
	}
	p.scraper.Stop()

	if p.scrapeOtel {
		p.exporter.Shutdown(context.Background())
		p.receiver.Shutdown(context.Background())
	}

	p.logger.Info("pktvisor process stopped", zap.Int("pid", finalStatus.PID), zap.Int("exit_code", finalStatus.Exit))
	return nil
}

func (p *pktvisorBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo, config map[string]string, otelConfig map[string]interface{}) error {
	p.logger = logger
	p.policyRepo = repo

	var prs bool
	if p.binary, prs = config["binary"]; !prs {
		return errors.New("you must specify pktvisor binary")
	}
	if p.configFile, prs = config["config_file"]; !prs {
		p.configFile = ""
	}
	if p.adminAPIHost, prs = config["api_host"]; !prs {
		return errors.New("you must specify pktvisor admin API host")
	}
	if p.adminAPIPort, prs = config["api_port"]; !prs {
		return errors.New("you must specify pktvisor admin API port")
	}

	for k, v := range otelConfig {
		switch k {
		case "Enable":
			p.scrapeOtel = v.(bool)
		}
	}

	return nil
}

func (p *pktvisorBackend) scrapeMetrics(period uint) (map[string]interface{}, error) {
	var metrics map[string]interface{}
	err := p.request(fmt.Sprintf("policies/__all/metrics/bucket/%d", period), &metrics, http.MethodGet, nil, "", 5)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (p *pktvisorBackend) GetCapabilities() (map[string]interface{}, error) {
	var taps interface{}
	err := p.request("taps", &taps, http.MethodGet, nil, "", 5)
	if err != nil {
		return nil, err
	}
	jsonBody := make(map[string]interface{})
	jsonBody["taps"] = taps
	return jsonBody, nil
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{
		adminAPIProtocol: "http",
	})
	return true
}

func createOtlpExporter(ctx context.Context, logger *zap.Logger) (component.MetricsExporter, error) {
	cfg := otlpexporter.CreateDefaultConfig()
	set := otlpexporter.CreateDefaultSettings(logger)
	// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
	exporter, err := otlpexporter.CreateMetricsExporter(ctx, set, cfg)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func createReceiver(ctx context.Context, exporter component.MetricsExporter, logger *zap.Logger) (component.MetricsReceiver, error) {
	set := pktvisorreceiver.CreateDefaultSettings(logger)
	cfg := pktvisorreceiver.CreateDefaultConfig()
	// Create the Prometheus receiver and pass in the preivously created Prometheus exporter.
	receiver, err := pktvisorreceiver.CreateMetricsReceiver(ctx, set, cfg, exporter)
	if err != nil {
		return nil, err
	}
	return receiver, nil
}

func (p *pktvisorBackend) FullReset() error {

	// force a stop, which stops scrape as well. if proc is dead, it no ops.
	if err := p.Stop(); err != nil {
		p.logger.Error("failed to stop backend on restart procedure", zap.Error(err))
		return err
	}

	// start it
	if err := p.Start(); err != nil {
		p.logger.Error("failed to start backend on restart procedure", zap.Error(err))
		return err
	}

	return nil
}
