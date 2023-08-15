package otel

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"github.com/amenzhinsky/go-memexec"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/orb-community/orb/agent/backend"
	"github.com/orb-community/orb/agent/config"
	"github.com/orb-community/orb/agent/policies"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"os"
	"time"
)

var _ backend.Backend = (*openTelemetryBackend)(nil)

//go:embed otelcol-contrib
var openTelemetryContribBinary []byte

type openTelemetryBackend struct {
	logger    *zap.Logger
	startTime time.Time

	//policies
	policyRepo            policies.PolicyRepo
	policyConfigDirectory string
	agentTags             map[string]string

	// Context for controlling the context cancellation
	mainContext        context.Context
	runningCollectors  map[string]context.Context
	mainCancelFunction context.CancelFunc

	// MQTT Config for OTEL MQTT Exporter
	mqttConfig config.MQTTConfig
	mqttClient *mqtt.Client

	otlpMetricsTopic string
	otlpTracesTopic  string
	otlpLogsTopic    string
	otelReceiverTaps []string

	otelReceiverHost string
	otelReceiverPort int

	metricsReceiver receiver.Metrics
	metricsExporter exporter.Metrics
	//tracesReceiver  receiver.Traces
	//tracesExporter  exporter.Traces
	//logsReceiver    receiver.Logs
	//logsExporter    exporter.Logs
}

// Configure initializes the backend with the given configuration
func (o openTelemetryBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo,
	_ map[string]string, otelConfig map[string]interface{}) error {
	o.logger = logger
	o.logger.Info("configuring OpenTelemetry backend")
	o.policyRepo = repo
	var err error
	o.otelReceiverTaps = []string{}
	o.policyConfigDirectory, err = os.MkdirTemp("", "otel-policies")
	if err != nil {
		o.logger.Error("failed to create temporary directory for policy configs", zap.Error(err))
		return err
	}
	if agentTags, ok := otelConfig["agent_tags"]; ok {
		o.agentTags = agentTags.(map[string]string)
	}
	for k, v := range otelConfig {
		switch k {
		case "Host":
			o.otelReceiverHost = v.(string)
		case "Port":
			o.otelReceiverPort = v.(int)
		}
	}

	return nil
}

func (o openTelemetryBackend) Version() (string, error) {
	executable, err := memexec.New(openTelemetryContribBinary)
	if err != nil {
		return "", err
	}
	defer func(executable *memexec.Exec) {
		err := executable.Close()
		if err != nil {
			o.logger.Error("error closing executable", zap.Error(err))
		}
	}(executable)
	ctx, cancel := context.WithTimeout(o.mainContext, 5*time.Second)
	defer cancel()
	cmd := executable.CommandContext(ctx, "--version")
	if cmd.Err != nil {
		return "", cmd.Err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	var versionOutput string
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			o.logger.Info("DEBUG", zap.String("line", scanner.Text()))
			versionOutput = scanner.Text()
		}
	}()
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	o.logger.Info("running opentelemetry-contrib version", zap.String("version", versionOutput))
	return versionOutput, nil
}

func (o openTelemetryBackend) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	o.runningCollectors = make(map[string]context.Context)
	o.mainCancelFunction = cancelFunc
	o.mainContext = ctx
	o.startTime = time.Now()
	// apply sample policy - remove after POC
	err := o.ApplyPolicy(samplePolicy, false)
	if err != nil {
		o.logger.Error("error updating policies", zap.Error(err))
		return err
	}
	currentVersion, err := o.Version()
	if err != nil {
		o.logger.Error("error during ")
	}
	o.logger.Info("starting open-telemetry backend using version", zap.String("version", currentVersion))

	policiesData, err := o.policyRepo.GetAll()
	if err != nil {
		defer cancelFunc()
		o.logger.Error("failed to start otel backend, policies are absent")
		return err
	}
	for _, policyData := range policiesData {
		if err := o.ApplyPolicy(policyData, true); err != nil {
			o.logger.Error("failed to start otel backend, failed to apply policy", zap.Error(err))
			cancelFunc()
			return err
		}
		o.logger.Info("policy applied successfully", zap.String("policy_id", policyData.ID))
	}

	return nil
}

func (o openTelemetryBackend) Stop(_ context.Context) error {
	o.logger.Info("stopping all running policies")
	o.mainCancelFunction()
	for policyID, policyCtx := range o.runningCollectors {
		o.logger.Debug("stopping policy context", zap.String("policy_id", policyID))
		policyCtx.Done()
	}
	return nil
}

func (o openTelemetryBackend) FullReset(_ context.Context) error {
	o.logger.Info("resetting all policies and restarting")
	for policyID, policyCtx := range o.runningCollectors {
		o.logger.Debug("stopping policy context", zap.String("policy_id", policyID))
		policyCtx.Done()
		policy, err := o.policyRepo.Get(policyID)
		if err != nil {
			o.logger.Error("failed to get policy", zap.Error(err))
			return err
		}
		err = o.ApplyPolicy(policy, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func Register() bool {
	backend.Register("otel", &openTelemetryBackend{})
	return true
}

func (o openTelemetryBackend) GetStartTime() time.Time {
	return o.startTime
}

// this will only print a default backend config
func (o openTelemetryBackend) GetCapabilities() (capabilities map[string]interface{}, err error) {
	capabilities["taps"] = o.otelReceiverTaps
	capabilities["version"], err = o.Version()
	if err != nil {
		return
	}
	return
}

// cross reference the Processes using the os, with the policies and contexts
func (o openTelemetryBackend) GetRunningStatus() (backend.RunningStatus, string, error) {
	amountCollectors := len(o.runningCollectors)
	if amountCollectors > 0 {
		return backend.Running, fmt.Sprintf("opentelemetry backend running with %d policies", amountCollectors), nil
	}
	return backend.Offline, "opentelemetry backend offline, waiting for policy to come to start running", nil
}
