package otel

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

func StartOtelComponents(ctx context.Context, logger zap.Logger) (context.CancelFunc, error) {
	otelContext, otelCancelFunc := context.WithCancel(ctx)

	log := logger.Sugar()
	log.Info("Starting to create Otel Components", ctx.Value("routine"))
	var bla kafkaexporter.Config
	log.Info("load info on", bla)
	exporterFactory := kafkaexporter.NewFactory()
	exporterCtx := context.WithValue(otelContext, "component", "kafkaexporter")
	set := component.ExporterCreateSettings{
		TelemetrySettings: component.TelemetrySettings{},
		BuildInfo:         component.BuildInfo{},
	}
	cfg := exporterFactory.CreateDefaultConfig().(*kafkaexporter.Config)
	cfg.Brokers = []string{"kafka1:9092"}
	exporter, err := exporterFactory.CreateMetricsExporter(exporterCtx, set, cfg)
	if err != nil {
		log.Error("error on creating exporter", err)
		otelCancelFunc()
		return nil, err
	}
	err = exporter.Start(exporterCtx, nil)
	if err != nil {
		log.Error("error on starting exporter", err)
		otelCancelFunc()
		return nil, err
	}
	// receiver Factory

	return otelCancelFunc, nil
}
