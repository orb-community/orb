/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-cmd/cmd"
	"github.com/go-co-op/gocron"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/agent/policies"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

const (
	DefaultBinary       = "/usr/local/sbin/pktvisord"
	ReadinessBackoff    = 10
	ReadinessTimeout    = 10
	ApplyPolicyTimeout  = 10
	RemovePolicyTimeout = 20
	VersionTimeout      = 2
	ScrapeTimeout       = 5
	TapsTimeout         = 5
)

// AppInfo represents server application information
type AppInfo struct {
	App struct {
		Version   string  `json:"version"`
		UpTimeMin float64 `json:"up_time_min"`
	} `json:"app"`
}

type pktvisorBackend struct {
	logger          *zap.Logger
	binary          string
	configFile      string
	pktvisorVersion string
	proc            *cmd.Cmd
	statusChan      <-chan cmd.Status
	startTime       time.Time
	cancelFunc      context.CancelFunc
	ctx             context.Context

	// MQTT Config for OTEL MQTT Exporter
	mqttConfig config.MQTTConfig

	mqttClient       mqtt.Client
	metricsTopic     string
	otlpMetricsTopic string
	scraper          *gocron.Scheduler
	policyRepo       policies.PolicyRepo

	receiver map[string]component.MetricsReceiver
	exporter map[string]component.MetricsExporter

	adminAPIHost     string
	adminAPIPort     string
	adminAPIProtocol string

	scrapeOtel bool

	// Go routine manager
	RoutineMap     map[string]context.CancelFunc
	RoutineChannel map[string]chan struct{}
}

func (p *pktvisorBackend) AddGoroutine(cancel context.CancelFunc, key string) {
	p.RoutineMap[key] = cancel
}

func (p *pktvisorBackend) KillGoroutine(key string) {
	cancel := p.RoutineMap[key]
	if cancel != nil {
		cancel()
	}
}

func (p *pktvisorBackend) GetOtelEnabled() bool {
	return p.scrapeOtel
}

func (p *pktvisorBackend) RestartScrapeOpenTelemetry(policyID string, policyName string) {
	p.KillGoroutine(policyID)
	exeCtx, execCancelF := context.WithCancel(p.ctx)
	p.AddGoroutine(execCancelF, policyID)
	attributeCtx := context.WithValue(exeCtx, "policy_id", policyID)
	attributeCtx = context.WithValue(attributeCtx, "policy_name", policyName)
	attributeCtx = context.WithValue(attributeCtx, "cancelFunc", execCancelF)
	p.scrapeOpenTelemetry(attributeCtx)
}

func (p *pktvisorBackend) GetStartTime() time.Time {
	return p.startTime
}

func (p *pktvisorBackend) SetCommsClient(agentID string, client mqtt.Client, baseTopic string) {
	p.mqttClient = client
	metricsTopic := strings.Replace(baseTopic, "?", "be", 1)
	otelMetricsTopic := strings.Replace(baseTopic, "?", "otlp", 1)
	p.metricsTopic = fmt.Sprintf("%s/m/%c", metricsTopic, agentID[0])
	p.otlpMetricsTopic = fmt.Sprintf("%s/m/%c", otelMetricsTopic, agentID[0])
}

func (p *pktvisorBackend) GetRunningStatus() (backend.RunningStatus, string, error) {
	// first check process status
	runningStatus, errMsg, err := p.getProcRunningStatus()
	// if it's not running, we're done
	if runningStatus != backend.Running {
		return runningStatus, errMsg, err
	}
	// if it's running, check REST API availability too
	_, aiErr := p.getAppInfo()
	if aiErr != nil {
		// process is running, but REST API is not accessible
		return backend.BackendError, "process running, REST API unavailable", aiErr
	}
	return runningStatus, "", nil
}

func (p *pktvisorBackend) Version() (string, error) {
	appInfo, err := p.getAppInfo()
	if err != nil {
		return "", err
	}
	p.pktvisorVersion = appInfo.App.Version
	return appInfo.App.Version, nil
}

func (p *pktvisorBackend) Start(ctx context.Context, cancelFunc context.CancelFunc) error {

	// this should record the start time whether it's successful or not
	// because it is used by the automatic restart system for last attempt
	p.startTime = time.Now()
	p.cancelFunc = cancelFunc
	p.ctx = ctx

	if p.RoutineMap == nil {
		p.RoutineMap = make(map[string]context.CancelFunc)
	}

	if p.receiver == nil {
		p.receiver = make(map[string]component.MetricsReceiver)
	}

	if p.exporter == nil {
		p.exporter = make(map[string]component.MetricsExporter)
	}

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

	// the macros should be properly configured to enable crashpad
	// pvOptions = append(pvOptions, "--cp-token", PKTVISOR_CP_TOKEN)
	// pvOptions = append(pvOptions, "--cp-url", PKTVISOR_CP_URL)
	// pvOptions = append(pvOptions, "--cp-path", PKTVISOR_CP_PATH)
	// pvOptions = append(pvOptions, "--default-geo-city", "/geo-db/city.mmdb")
	// pvOptions = append(pvOptions, "--default-geo-asn", "/geo-db/asn.mmdb")

	p.logger.Info("pktvisor startup", zap.Strings("arguments", pvOptions))

	p.proc = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, p.binary, pvOptions...)
	p.statusChan = p.proc.Start()

	// log STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer func() {
			if doneChan != nil {
				close(doneChan)
			}
		}()
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
		var appMetrics AppInfo
		readinessError = p.request("metrics/app", &appMetrics, http.MethodGet, http.NoBody, "application/json", ReadinessTimeout)
		if readinessError == nil {
			p.logger.Info("pktvisor readiness ok, got version ", zap.String("pktvisor_version", appMetrics.App.Version))
			break
		}
		backoffDuration := time.Duration(backoff) * time.Second
		p.logger.Info("pktvisor is not ready, trying again with backoff", zap.String("backoff backoffDuration", backoffDuration.String()))
		time.Sleep(backoffDuration)
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

	if !p.scrapeOtel {
		if err := p.scrapeDefault(); err != nil {
			return err
		}
	}

	return nil
}

func (p *pktvisorBackend) Stop(ctx context.Context) error {
	p.logger.Info("routine call to stop pktvisor", zap.Any("routine", ctx.Value("routine")))
	defer p.cancelFunc()
	err := p.proc.Stop()
	finalStatus := <-p.statusChan
	if err != nil {
		p.logger.Error("pktvisor shutdown error", zap.Error(err))
	}
	p.scraper.Stop()
	if p.scrapeOtel {
		for key, cancelScrap := range p.RoutineMap {
			if cancelScrap != nil {
				cancelScrap()
			}
			p.logger.Info("Requested to stop scrap function policy: " + key)
		}
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
			p.logger.Info("OpenTelemetry enabled")
			p.scrapeOtel = v.(bool)
		}
	}

	return nil
}

func (p *pktvisorBackend) GetCapabilities() (map[string]interface{}, error) {
	var taps interface{}
	err := p.request("taps", &taps, http.MethodGet, http.NoBody, "application/json", TapsTimeout)
	if err != nil {
		return nil, err
	}
	jsonBody := make(map[string]interface{})
	jsonBody["taps"] = taps
	return jsonBody, nil
}

func (p *pktvisorBackend) FullReset(ctx context.Context) error {

	// force a stop, which stops scrape as well. if proc is dead, it no ops.
	if state, _, _ := p.getProcRunningStatus(); state == backend.Running {
		if err := p.Stop(ctx); err != nil {
			p.logger.Error("failed to stop backend on restart procedure", zap.Error(err))
			return err
		}
	}

	backendCtx, cancelFunc := context.WithCancel(context.WithValue(ctx, "routine", "pktvisor"))
	// start it
	if err := p.Start(backendCtx, cancelFunc); err != nil {
		p.logger.Error("failed to start backend on restart procedure", zap.Error(err))
		return err
	}

	return nil
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{
		adminAPIProtocol: "http",
	})
	return true
}
