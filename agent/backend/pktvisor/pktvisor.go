/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/go-cmd/cmd"
	"github.com/go-co-op/gocron"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/policies"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

const DefaultBinary = "/usr/local/sbin/pktvisord"

type pktvisorBackend struct {
	logger          *zap.Logger
	binary          string
	configFile      string
	pktvisorVersion string
	proc            *cmd.Cmd
	statusChan      <-chan cmd.Status
	args            []string

	mqttClient   mqtt.Client
	metricsTopic string
	scraper      *gocron.Scheduler

	adminAPIHost     string
	adminAPIPort     string
	adminAPIProtocol string
}

func (p *pktvisorBackend) SetCommsClient(client mqtt.Client, baseTopic string) {
	p.mqttClient = client
	p.metricsTopic = fmt.Sprintf("%s/m", baseTopic)
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
func (p *pktvisorBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string) error {
	client := http.Client{
		Timeout: time.Second * 5,
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
		if body[0] == '{' {
			var jsonBody map[string]interface{}
			err := json.Unmarshal(body, &jsonBody)
			if err == nil {
				if errMsg, ok := jsonBody["error"]; ok {
					return errors.New(fmt.Sprintf("%d %s", res.StatusCode, errMsg))
				}
			}
		}
		return errors.New(fmt.Sprintf("%d %s", res.StatusCode, body))
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
		p.proc.Stop()
		// TODO auto restart
		return false, errors.New("pktvisor process ended")
	}

	return true, nil
}

func (p *pktvisorBackend) ApplyPolicy(data policies.PolicyData) error {

	p.logger.Debug("pktvisor policy apply", zap.String("policy_id", data.ID), zap.Any("data", data.Data))

	// add header with "p_data.ID" as policy name
	// note pktvisor doesn't allow policy names to start with a number, thus the prefix
	pktvisorPolicyName := fmt.Sprintf("p_%s", data.ID)
	fullPolicy := map[string]interface{}{
		"version": "1.0",
		"visor": map[string]interface{}{
			"policies": map[string]interface{}{
				pktvisorPolicyName: data.Data,
			},
		},
	}

	pyaml, err := yaml.Marshal(fullPolicy)
	if err != nil {
		p.logger.Warn("yaml policy marshal failure", zap.String("policy_id", data.ID), zap.Any("policy", fullPolicy))
		return err
	}

	var resp map[string]interface{}
	err = p.request("policies", &resp, http.MethodPost, bytes.NewBuffer(pyaml), "application/x-yaml")
	if err != nil {
		p.logger.Warn("yaml policy application failure", zap.String("policy_id", data.ID), zap.ByteString("policy", pyaml))
		return err
	}

	job, err := p.scraper.Every(1).Minute().WaitForSchedule().Tag(data.ID).Do(func() {
		metrics, err := p.scrapeMetrics(data.ID, 1)
		if err != nil {
			p.logger.Error("scrape failed", zap.String("policy_id", data.ID), zap.Error(err))
			return
		}
		payloadData, err := json.Marshal(metrics)
		if err != nil {
			p.logger.Error("error marshalling scraped metric json", zap.String("policy_id", data.ID), zap.Error(err))
			return
		}
		metricPayload := fleet.AgentMetricsRPCPayload{
			PolicyID:  data.ID,
			Datasets:  data.GetDatasetIDs(),
			Format:    "json",
			BEVersion: p.pktvisorVersion,
			Data:      payloadData,
		}
		rpc := fleet.AgentMetricsRPC{
			SchemaVersion: fleet.CurrentRPCSchemaVersion,
			Func:          fleet.AgentMetricsRPCFunc,
			Payload:       []fleet.AgentMetricsRPCPayload{metricPayload},
		}

		topic := fmt.Sprintf("%s/%c", p.metricsTopic, data.ID[0])
		body, err := json.Marshal(rpc)
		if err != nil {
			p.logger.Error("error marshalling metric rpc payload", zap.String("policy_id", data.ID), zap.Error(err))
			return
		}

		if token := p.mqttClient.Publish(topic, 1, false, body); token.Wait() && token.Error() != nil {
			p.logger.Error("error sending metrics RPC", zap.String("topic", topic), zap.Error(token.Error()))
		}
		p.logger.Info("scrapped and published metrics", zap.String("policy_id", data.ID), zap.String("topic", topic), zap.Int("payload_size", len(payloadData)))

	})
	job.SingletonMode()

	if err != nil {
		p.logger.Warn("application succeeded but scraper creation failed, attempting remove policy", zap.String("policy_id", data.ID))
		rerr := p.RemovePolicy(data.ID)
		if rerr != nil {
			p.logger.Error("policy removal failed", zap.String("policy_id", data.ID), zap.Error(rerr))
		}
		return err
	}

	return nil

}

func (p *pktvisorBackend) RemovePolicy(policyID string) error {
	var resp interface{}
	err := p.request(fmt.Sprintf("policies/%s", policyID), &resp, http.MethodDelete, nil, "")
	if err != nil {
		return err
	}

	err = p.scraper.RemoveByTag(policyID)
	if err != nil {
		return err
	}

	return nil

}

func (p *pktvisorBackend) Version() (string, error) {
	var appMetrics AppMetrics
	err := p.request("metrics/app", &appMetrics, http.MethodGet, nil, "")
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

	_, err := exec.LookPath(p.binary)
	if err != nil {
		p.logger.Error("pktvisor startup error: binary not found", zap.Error(err))
		return err
	}

	pvOptions := []string{
		strings.Join(p.args, " "),
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
		p.proc.Stop()
		return errors.New("pktvisor startup error, check log")
	}

	p.logger.Info("pktvisor process started", zap.Int("pid", status.PID))

	p.scraper = gocron.NewScheduler(time.UTC)
	p.scraper.StartAsync()

	return nil
}

func (p *pktvisorBackend) Stop() error {
	p.logger.Info("pktvisor stopping")
	err := p.proc.Stop()
	finalStatus := <-p.statusChan
	if err != nil {
		p.logger.Error("pktvisor shutdown error", zap.Error(err))
		return err
	}
	p.scraper.Stop()
	p.logger.Info("pktvisor process stopped", zap.Int("pid", finalStatus.PID), zap.Int("exit_code", finalStatus.Exit))
	return nil
}

func (p *pktvisorBackend) Configure(logger *zap.Logger, config map[string]interface{}) error {
	p.logger = logger
	var prs bool
	binary, prs := config["binary"]
	if !prs {
		return errors.New("you must specify pktvisor binary")
	} else {
		p.binary = fmt.Sprint(binary)
	}
	args, prs := config["binary_args"]
	if !prs {
		return errors.New("you must specify binary args")
	} else {
		s := reflect.ValueOf(args)
		if s.Kind() != reflect.Slice {
			panic("Interface griven a non-slice type")
		}
		for i := 0; i < s.Len(); i++ {
			p.args = append(p.args, fmt.Sprint(s.Index(i).Interface()))
		}
	}
	configFile, prs := config["config_file"]
	if !prs {
		p.configFile = ""
	} else {
		p.configFile = fmt.Sprint(configFile)
	}

	adminAPIHost, prs := config["api_host"]
	if !prs {
		return errors.New("you must specify pktvisor admin API host")
	} else {
		p.adminAPIHost = fmt.Sprint(adminAPIHost)
	}

	adminAPIPort, prs := config["api_port"]
	if !prs {
		return errors.New("you must specify pktvisor admin API port")
	} else {
		p.adminAPIPort = fmt.Sprint(adminAPIPort)
	}
	return nil
}

func (p *pktvisorBackend) scrapeMetrics(policyID string, period uint) (map[string]interface{}, error) {
	var metrics map[string]interface{}
	err := p.request(fmt.Sprintf("policies/p_%s/metrics/bucket/%d", policyID, period), &metrics, http.MethodGet, nil, "")
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (p *pktvisorBackend) GetCapabilities() (map[string]interface{}, error) {
	var taps interface{}
	err := p.request("taps", &taps, http.MethodGet, nil, "")
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
