package otel

import (
	"context"
	_ "embed"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-cmd/cmd"
	"github.com/orb-community/orb/agent/backend"
	"github.com/orb-community/orb/agent/config"
	"github.com/orb-community/orb/agent/otel"
	"github.com/orb-community/orb/agent/otel/otlpmqttexporter"
	"github.com/orb-community/orb/agent/policies"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

var _ backend.Backend = (*openTelemetryBackend)(nil)

const DefaultPath = "/usr/local/bin/otelcol-contrib"
const DefaultHost = "localhost"
const DefaultPort = 4317

type openTelemetryBackend struct {
	logger    *zap.Logger
	startTime time.Time

	//policies
	policyRepo            policies.PolicyRepo
	policyConfigDirectory string
	agentTags             map[string]string

	// Context for controlling the context cancellation
	mainContext        context.Context
	runningCollectors  map[string]runningPolicy
	mainCancelFunction context.CancelFunc

	// MQTT Config for OTEL MQTT Exporter
	mqttConfig config.MQTTConfig
	mqttClient *mqtt.Client

	otlpMetricsTopic string
	otlpTracesTopic  string
	otlpLogsTopic    string
	otelReceiverTaps []string
	otelCurrVersion  string

	otelReceiverHost   string
	otelReceiverPort   int
	otelExecutablePath string

	metricsReceiver receiver.Metrics
	metricsExporter exporter.Metrics
	tracesReceiver  receiver.Traces
	tracesExporter  exporter.Traces
	logsReceiver    receiver.Logs
	logsExporter    exporter.Logs
}

// Configure initializes the backend with the given configuration
func (o *openTelemetryBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo,
	config map[string]string, otelConfig map[string]interface{}) error {
	o.logger = logger
	o.logger.Info("configuring OpenTelemetry backend")
	o.policyRepo = repo
	var err error
	o.otelReceiverTaps = []string{"otelcol-contrib", "receivers", "processors", "extensions"}
	o.policyConfigDirectory, err = os.MkdirTemp("", "otel-policies")
	if path, ok := config["binary"]; ok {
		o.otelExecutablePath = path
	} else {
		o.otelExecutablePath = DefaultPath
	}
	if err != nil {
		o.logger.Error("failed to create temporary directory for policy configs", zap.Error(err))
		return err
	}
	if agentTags, ok := otelConfig["agent_tags"]; ok {
		o.agentTags = agentTags.(map[string]string)
	}
	if otelPort, ok := config["otlp_port"]; ok {
		o.otelReceiverPort, err = strconv.Atoi(otelPort)
		if err != nil {
			o.logger.Error("failed to parse otlp port using default", zap.Error(err))
			o.otelReceiverPort = DefaultPort
		}
	} else {
		o.otelReceiverPort = DefaultPort
	}
	if otelHost, ok := config["otlp_host"]; ok {
		o.otelReceiverHost = otelHost
	} else {
		o.otelReceiverHost = DefaultHost
	}

	return nil
}

func (o *openTelemetryBackend) GetInitialState() backend.RunningStatus {
	return backend.Waiting
}

func (o *openTelemetryBackend) Version() (string, error) {
	if o.otelCurrVersion != "" {
		return o.otelCurrVersion, nil
	}
	ctx, cancel := context.WithTimeout(o.mainContext, 60*time.Second)
	var versionOutput string
	command := cmd.NewCmd(o.otelExecutablePath, "--version")
	status := command.Start()
	select {
	case finalStatus := <-status:
		if finalStatus.Error != nil {
			o.logger.Error("error during call of otelcol-contrib version", zap.Error(finalStatus.Error))
			return "", finalStatus.Error
		} else {
			output := finalStatus.Stdout
			o.otelCurrVersion = output[0]
			versionOutput = output[0]
		}
	case <-ctx.Done():
		o.logger.Error("timeout during getting version", zap.Error(ctx.Err()))
	}

	cancel()
	o.logger.Info("running opentelemetry-contrib version", zap.String("version", versionOutput))

	return versionOutput, nil

}

func (o *openTelemetryBackend) Start(ctx context.Context, cancelFunc context.CancelFunc) (err error) {
	o.runningCollectors = make(map[string]runningPolicy)
	o.mainCancelFunction = cancelFunc
	o.mainContext = ctx
	o.startTime = time.Now()
	currentWd, err := os.Getwd()
	if err != nil {
		o.otelExecutablePath = currentWd + "/otelcol-contrib"
	}
	currentVersion, err := o.Version()
	if err != nil {
		cancelFunc()
		o.logger.Error("error during getting current version", zap.Error(err))
		return err
	}
	o.receiveOtlp()
	o.logger.Info("starting open-telemetry backend using version", zap.String("version", currentVersion))
	policiesData, err := o.policyRepo.GetAll()
	if err != nil {
		cancelFunc()
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

func (o *openTelemetryBackend) Stop(_ context.Context) error {
	o.logger.Info("stopping all running policies")
	o.mainCancelFunction()
	for policyID, policyEntry := range o.runningCollectors {
		o.logger.Debug("stopping policy context", zap.String("policy_id", policyID))
		policyEntry.ctx.Done()
	}
	return nil
}

func (o *openTelemetryBackend) FullReset(ctx context.Context) error {
	o.logger.Info("restarting otel backend", zap.Int("running policies", len(o.runningCollectors)))
	backendCtx, cancelFunc := context.WithCancel(context.WithValue(ctx, "routine", "otel"))
	if err := o.Start(backendCtx, cancelFunc); err != nil {
		return err
	}
	return nil
}

func Register() bool {
	backend.Register("otel", &openTelemetryBackend{})
	return true
}

func (o *openTelemetryBackend) GetStartTime() time.Time {
	return o.startTime
}

// GetCapabilities this will only print a default backend config
func (o *openTelemetryBackend) GetCapabilities() (capabilities map[string]interface{}, err error) {
	capabilities = make(map[string]interface{})
	capabilities["taps"] = o.otelReceiverTaps
	return
}

// GetRunningStatus returns cross-reference the Processes using the os, with the policies and contexts
func (o *openTelemetryBackend) GetRunningStatus() (backend.RunningStatus, string, error) {
	amountCollectors := len(o.runningCollectors)
	if amountCollectors > 0 {
		return backend.Running, fmt.Sprintf("opentelemetry backend running with %d policies", amountCollectors), nil
	}
	return backend.Waiting, "opentelemetry backend is waiting for policy to come to start running", nil
}

func (o *openTelemetryBackend) createOtlpMetricMqttExporter(ctx context.Context, cancelFunc context.CancelCauseFunc) (exporter.Metrics, error) {
	bridgeService := otel.NewBridgeService(ctx, cancelFunc, &o.policyRepo, o.agentTags)
	var cfg component.Config
	if o.mqttClient != nil {
		cfg = otlpmqttexporter.CreateConfigClient(o.mqttClient, o.otlpMetricsTopic, "", bridgeService)
	} else {
		cfg = otlpmqttexporter.CreateConfig(o.mqttConfig.Address, o.mqttConfig.Id, o.mqttConfig.Key,
			o.mqttConfig.ChannelID, "", o.otlpMetricsTopic, bridgeService)
	}

	set := otlpmqttexporter.CreateDefaultSettings(o.logger)
	// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
	return otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)

}

func (o *openTelemetryBackend) createOtlpTraceMqttExporter(ctx context.Context, cancelFunc context.CancelCauseFunc) (exporter.Traces, error) {
	bridgeService := otel.NewBridgeService(ctx, cancelFunc, &o.policyRepo, o.agentTags)
	if o.mqttClient != nil {
		cfg := otlpmqttexporter.CreateConfigClient(o.mqttClient, o.otlpTracesTopic, "", bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(o.logger)
		// Create the OTLP metrics metricsExporter that'll receive and verify the metrics produced.
		tracerExporter, err := otlpmqttexporter.CreateTracesExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return tracerExporter, nil
	} else {
		cfg := otlpmqttexporter.CreateConfig(o.mqttConfig.Address, o.mqttConfig.Id, o.mqttConfig.Key,
			o.mqttConfig.ChannelID, "", o.otlpTracesTopic, bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(o.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		tracerExporter, err := otlpmqttexporter.CreateTracesExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return tracerExporter, nil
	}

}

func (o *openTelemetryBackend) createOtlpLogsMqttExporter(ctx context.Context, cancelFunc context.CancelCauseFunc) (exporter.Logs, error) {
	bridgeService := otel.NewBridgeService(ctx, cancelFunc, &o.policyRepo, o.agentTags)
	if o.mqttClient != nil {
		cfg := otlpmqttexporter.CreateConfigClient(o.mqttClient, o.otlpLogsTopic, "", bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(o.logger)
		// Create the OTLP metrics metricsExporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateLogsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	} else {
		cfg := otlpmqttexporter.CreateConfig(o.mqttConfig.Address, o.mqttConfig.Id, o.mqttConfig.Key,
			o.mqttConfig.ChannelID, "", o.otlpLogsTopic, bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(o.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateLogsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	}

}
