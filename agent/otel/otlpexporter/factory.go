package otlpexporter

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	typeStr        = "otlp"
	defaultEnpoint = "localhost:1234"
)

func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		CreateDefaultConfig,
		exporterhelper.WithMetrics(CreateMetricsExporter))
}

func CreateDefaultSettings(logger *zap.Logger) component.ExporterCreateSettings {
	return component.ExporterCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.GetMeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
}

func CreateDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		TimeoutSettings:  exporterhelper.DefaultTimeoutSettings(),
		QueueSettings:    exporterhelper.DefaultQueueSettings(),
		RetrySettings:    exporterhelper.DefaultRetrySettings(),
		GRPCClientSettings: configgrpc.GRPCClientSettings{
			Endpoint:        defaultEnpoint,
			Headers:         map[string]string{},
			WriteBufferSize: 512 * 1024,
			TLSSetting: configtls.TLSClientSetting{
				Insecure: true,
			},
		},
	}
}

func CreateMetricsExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.MetricsExporter, error) {
	oce, err := newExporter(cfg, set.TelemetrySettings, set.BuildInfo)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)
	return exporterhelper.NewMetricsExporter(
		cfg,
		set,
		oce.pushMetrics,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(oCfg.TimeoutSettings),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings),
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithShutdown(oce.shutdown),
	)
}
