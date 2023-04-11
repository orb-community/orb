package otlpmqttexporter

import (
	"context"
	"fmt"

	"github.com/orb-community/orb/agent/otel"
	"go.uber.org/zap"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
)

const (
	// The value of "type" key in configuration.
	typeStr         = "otlpmqtt"
	defaultMQTTAddr = "localhost:1883"
	defaultMQTTId   = "uuid1"
	defaultMQTTKey  = "uuid2"
	defaultName     = "pktvisor"
	// For testing will disable  TLS
	defaultTLS = false
)

// NewFactory creates a factory for OTLP exporter.
// Reducing the scope to just Metrics since it is our use-case
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		CreateDefaultConfig,
		exporter.WithMetrics(CreateMetricsExporter, component.StabilityLevelStable))
}

func CreateConfig(addr, id, key, channel, pktvisor, metricsTopic string, bridgeService otel.AgentBridgeService) component.Config {
	return &Config{
		TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
		QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
		MetricsTopic:    metricsTopic,
		Address:         addr,
		Id:              id,
		Key:             key,
		ChannelID:       channel,
		PktVisorVersion: pktvisor,
		OrbAgentService: bridgeService,
	}
}

func CreateDefaultSettings(logger *zap.Logger) exporter.CreateSettings {
	return exporter.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
}

func CreateDefaultConfig() component.Config {
	base := fmt.Sprintf("channels/%s/messages", defaultMQTTId)
	metricsTopic := fmt.Sprintf("%s/otlp/%s", base, defaultName)
	return &Config{
		TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
		QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
		Address:         defaultMQTTAddr,
		Id:              defaultMQTTId,
		Key:             defaultMQTTKey,
		ChannelID:       base,
		TLS:             defaultTLS,
		MetricsTopic:    metricsTopic,
	}
}

func CreateConfigClient(client *mqtt.Client, metricsTopic, pktvisor string, bridgeService otel.AgentBridgeService) component.Config {
	return &Config{
		TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
		QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
		Client:          client,
		MetricsTopic:    metricsTopic,
		PktVisorVersion: pktvisor,
		OrbAgentService: bridgeService,
	}
}

func createTracesExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	oce, err := newExporter(cfg, set, ctx)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)

	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		oce.pushTraces,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings))
}

func CreateMetricsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Metrics, error) {
	oce, err := newExporter(cfg, set, ctx)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)
	pFunc := oce.pushMetrics
	if ctx.Value("all").(bool) {
		pFunc = oce.pushAllMetrics
	}
	return exporterhelper.NewMetricsExporter(
		ctx,
		set,
		cfg,
		pFunc,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings))
}

func createLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	oce, err := newExporter(cfg, set, ctx)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)

	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		oce.pushLogs,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings))
}
