package otel

import (
	"context"
	"github.com/ns1labs/orb/sinker/otel/orbreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func StartOtelComponents(ctx context.Context, logger *zap.Logger, metricsChannel chan []byte) (context.CancelFunc, error) {
	otelContext, otelCancelFunc := context.WithCancel(ctx)

	log := logger.Sugar()
	log.Info("Starting to create Otel Components", ctx.Value("routine"))
	var bla kafkaexporter.Config
	log.Info("load info on", bla)
	exporterFactory := kafkaexporter.NewFactory()
	exporterCtx := context.WithValue(otelContext, "component", "kafkaexporter")
	exporterCreateSettings := component.ExporterCreateSettings{
		TelemetrySettings: component.TelemetrySettings{},
		BuildInfo:         component.BuildInfo{},
	}
	expCfg := exporterFactory.CreateDefaultConfig().(*kafkaexporter.Config)
	expCfg.Brokers = []string{"kafka1:9092"}
	exporter, err := exporterFactory.CreateMetricsExporter(exporterCtx, exporterCreateSettings, expCfg)
	if err != nil {
		log.Error("error on creating exporter", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = exporter.Start(exporterCtx, nil)
	if err != nil {
		log.Error("error on starting exporter", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	// receiver Factory
	orbReceiverFactory := orbreceiver.NewFactory()
	receiverCtx := context.WithValue(otelContext, "component", "orbreceiver")
	receiverCfg := orbReceiverFactory.CreateDefaultConfig().(*orbreceiver.Config)
	receiverCfg.Logger = logger
	receiverCfg.MetricsChannel = metricsChannel
	receiverSet := component.ReceiverCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	receiver, err := orbReceiverFactory.CreateMetricsReceiver(receiverCtx, receiverSet, receiverCfg, exporter)
	if err != nil {
		log.Error("error on creating receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = receiver.Start(receiverCtx, nil)
	if err != nil {
		log.Error("error on starting receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	return otelCancelFunc, nil
}
