/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package diode

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-cmd/cmd"
	dconf "github.com/orb-community/diode/agent/config"
	"github.com/orb-community/orb/agent/backend"
	"github.com/orb-community/orb/agent/config"
	"github.com/orb-community/orb/agent/policies"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var _ backend.Backend = (*diodeBackend)(nil)

const (
	DefaultBinary       = "/usr/local/bin/diode-agent"
	DefaultHost         = "localhost"
	DefaultPort         = "10911"
	ReadinessBackoff    = 10
	ReadinessTimeout    = 10
	ApplyPolicyTimeout  = 10
	RemovePolicyTimeout = 20
	VersionTimeout      = 2
	ScrapeTimeout       = 5
	TapsTimeout         = 5
)

type diodeBackend struct {
	logger     *zap.Logger
	binary     string
	configFile string
	version    string
	proc       *cmd.Cmd
	statusChan <-chan cmd.Status
	startTime  time.Time
	cancelFunc context.CancelFunc
	ctx        context.Context

	// MQTT Config for OTEL MQTT Exporter
	mqttConfig config.MQTTConfig

	mqttClient *mqtt.Client
	logTopic   string
	policyRepo policies.PolicyRepo

	adminAPIHost     string
	adminAPIPort     string
	adminAPIProtocol string

	// added for Strings
	agentTags map[string]string

	// OpenTelemetry management
	otelReceiverType string
	otelReceiverHost string
	otelReceiverPort int
	receiver         receiver.Logs
	exporter         exporter.Logs
}

func Register() bool {
	backend.Register("diode", &diodeBackend{
		adminAPIProtocol: "http",
	})
	return true
}

func (d *diodeBackend) getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func (d *diodeBackend) GetStartTime() time.Time {
	return d.startTime
}

func (d *diodeBackend) GetCapabilities() (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

func (d *diodeBackend) GetRunningStatus() (backend.RunningStatus, string, error) {
	// first check process status
	runningStatus, errMsg, err := d.getProcRunningStatus()
	// if it's not running, we're done
	if runningStatus != backend.Running {
		return runningStatus, errMsg, err
	}
	// if it's running, check REST API availability too
	var dConf dconf.Status
	if err := d.request("status", &dConf, http.MethodGet, http.NoBody, "application/json", VersionTimeout); err != nil {
		return backend.BackendError, "process running, REST API unavailable", err
	}
	return runningStatus, "", nil
}

func (d *diodeBackend) Version() (string, error) {
	var dConf dconf.Status
	if err := d.request("status", &dConf, http.MethodGet, http.NoBody, "application/json", VersionTimeout); err != nil {
		return "", err
	}
	return dConf.Version, nil
}

func (d *diodeBackend) SetCommsClient(agentID string, client *mqtt.Client, baseTopic string) {
	d.mqttClient = client
	logTopic := strings.Replace(baseTopic, "?", "otlp", 1)
	d.logTopic = fmt.Sprintf("%s/m/%c", logTopic, agentID[0])
}

func (d *diodeBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo, config map[string]string, otelConfig map[string]interface{}) error {
	d.logger = logger
	d.policyRepo = repo

	var err bool
	if d.binary, err = config["binary"]; !err {
		d.binary = DefaultBinary
	}
	if d.configFile, err = config["config_file"]; !err {
		d.configFile = ""
	}
	if d.adminAPIHost, err = config["api_host"]; !err {
		d.adminAPIHost = DefaultHost
	}
	if d.adminAPIPort, err = config["api_port"]; !err {
		d.adminAPIPort = DefaultPort
	}
	if agentTags, ok := otelConfig["agent_tags"]; ok {
		d.agentTags = agentTags.(map[string]string)
	}
	return nil
}

func (d *diodeBackend) Start(ctx context.Context, cancelFunc context.CancelFunc) error {

	// this should record the start time whether it's successful or not
	// because it is used by the automatic restart system for last attempt
	d.startTime = time.Now()
	d.cancelFunc = cancelFunc
	d.ctx = ctx

	_, err := exec.LookPath(d.binary)
	if err != nil {
		d.logger.Error("diode-agent startup error: binary not found", zap.Error(err))
		return err
	}

	pvOptions := []string{
		"run",
		"-i",
		d.adminAPIHost,
		"-p",
		d.adminAPIPort,
	}
	if len(d.configFile) > 0 {
		pvOptions = append(pvOptions, "--config", d.configFile)
	}

	d.otelReceiverPort, err = d.getFreePort()

	d.logger.Info("diode-agent startup", zap.Strings("arguments", pvOptions))

	d.proc = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, d.binary, pvOptions...)
	d.statusChan = d.proc.Start()

	// log STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer func() {
			if doneChan != nil {
				close(doneChan)
			}
		}()
		for d.proc.Stdout != nil || d.proc.Stderr != nil {
			select {
			case line, open := <-d.proc.Stdout:
				if !open {
					d.proc.Stdout = nil
					continue
				}
				d.logger.Info("diode-agent stdout", zap.String("log", line))
			case line, open := <-d.proc.Stderr:
				if !open {
					d.proc.Stderr = nil
					continue
				}
				d.logger.Info("diode-agent stderr", zap.String("log", line))
			}
		}
	}()

	// wait for simple startup errors
	time.Sleep(time.Second)

	status := d.proc.Status()

	if status.Error != nil {
		d.logger.Error("diode-agent startup error", zap.Error(status.Error))
		return status.Error
	}

	if status.Complete {
		err = d.proc.Stop()
		if err != nil {
			d.logger.Error("proc.Stop error", zap.Error(err))
		}
		return errors.New("pktvisor startup error, check log")
	}

	d.logger.Info("pktvisor process started", zap.Int("pid", status.PID))

	var readinessError error
	for backoff := 0; backoff < ReadinessBackoff; backoff++ {
		var dConf dconf.Status
		readinessError = d.request("status", &dConf, http.MethodGet, http.NoBody, "application/json", ReadinessTimeout)
		if readinessError == nil {
			d.logger.Info("diode-agent readiness ok, got version ", zap.String("diode_agent_version", dConf.Version))
			break
		}
		backoffDuration := time.Duration(backoff) * time.Second
		d.logger.Info("diode-agent is not ready, trying again with backoff", zap.String("backoff backoffDuration", backoffDuration.String()))
		time.Sleep(backoffDuration)
	}

	if readinessError != nil {
		d.logger.Error("diode-agent error on readiness", zap.Error(readinessError))
		err = d.proc.Stop()
		if err != nil {
			d.logger.Error("proc.Stop error", zap.Error(err))
		}
		return readinessError
	}
	return nil
}

func (d *diodeBackend) Stop(ctx context.Context) error {
	d.logger.Info("routine call to stop diode-agent", zap.Any("routine", ctx.Value("routine")))
	defer d.cancelFunc()
	err := d.proc.Stop()
	finalStatus := <-d.statusChan
	if err != nil {
		d.logger.Error("diode-agent shutdown error", zap.Error(err))
	}
	d.logger.Info("diode-agent process stopped", zap.Int("pid", finalStatus.PID), zap.Int("exit_code", finalStatus.Exit))
	return nil
}

func (d *diodeBackend) FullReset(ctx context.Context) error {
	// force a stop, which stops scrape as well. if proc is dead, it no ops.
	if state, _, _ := d.getProcRunningStatus(); state == backend.Running {
		if err := d.Stop(ctx); err != nil {
			d.logger.Error("failed to stop backend on restart procedure", zap.Error(err))
			return err
		}
	}
	backendCtx, cancelFunc := context.WithCancel(context.WithValue(ctx, "routine", "diode-agent"))
	// start it
	if err := d.Start(backendCtx, cancelFunc); err != nil {
		d.logger.Error("failed to start backend on restart procedure", zap.Error(err))
		return err
	}
	return nil
}

func (d *diodeBackend) ApplyPolicy(data policies.PolicyData, updatePolicy bool) error {
	if updatePolicy {
		// To update a policy it's necessary first remove it and then apply a new version
		err := d.RemovePolicy(data)
		if err != nil {
			d.logger.Warn("policy failed to remove", zap.String("policy_id", data.ID), zap.String("policy_name", data.Name), zap.Error(err))
		}
	}
	d.logger.Debug("diode-agent policy apply", zap.String("policy_id", data.ID), zap.Any("data", data.Data))
	pol := map[string]interface{}{
		data.Name: data.Data,
	}
	policyYaml, err := yaml.Marshal(pol)
	if err != nil {
		d.logger.Warn("yaml policy marshal failure", zap.String("policy_id", data.ID), zap.Any("policy", pol))
		return err
	}
	var resp map[string]interface{}
	err = d.request("policies", &resp, http.MethodPost, bytes.NewBuffer(policyYaml), "application/x-yaml", ApplyPolicyTimeout)
	if err != nil {
		d.logger.Warn("yaml policy application failure", zap.String("policy_id", data.ID), zap.ByteString("policy", policyYaml))
		return err
	}
	return nil
}

func (d *diodeBackend) RemovePolicy(data policies.PolicyData) error {
	d.logger.Debug("diode-agent policy remove", zap.String("policy_id", data.ID))
	var resp interface{}
	var name string
	// Since we use Name for removing policies not IDs, if there is a change, we need to remove the previous name of the policy
	if data.PreviousPolicyData != nil && data.PreviousPolicyData.Name != data.Name {
		name = data.PreviousPolicyData.Name
	} else {
		name = data.Name
	}
	err := d.request(fmt.Sprintf("policies/%s", name), &resp, http.MethodDelete, http.NoBody, "application/json", RemovePolicyTimeout)
	if err != nil {
		return err
	}
	return nil
}
