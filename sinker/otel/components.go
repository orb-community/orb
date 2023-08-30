package otel

import (
	"context"

	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/orb-community/orb/sinker/otel/bridgeservice"
	kafkaexporter "github.com/orb-community/orb/sinker/otel/kafkafanoutexporter"

	"github.com/orb-community/orb/sinker/otel/orbreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func StartOtelMetricsComponents(ctx context.Context, bridgeService *bridgeservice.SinkerOtelBridgeService, logger *zap.Logger, kafkaUrl string, pubSub mfnats.PubSub) (context.CancelFunc, error) {
	otelContext, otelCancelFunc := context.WithCancel(ctx)

	log := logger.Sugar()
	log.Info("Starting to create Otel Metrics Components in routine: ", ctx.Value("routine"))
	exporterFactory := kafkaexporter.NewFactory()
	exporterCtx := context.WithValue(otelContext, "component", "kafkaexporter")
	exporterCreateSettings := exporter.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	expCfg := exporterFactory.CreateDefaultConfig().(*kafkaexporter.Config)
	expCfg.Brokers = []string{kafkaUrl}
	expCfg.Topic = "otlp_metrics"
	exporter, err := exporterFactory.CreateMetricsExporter(exporterCtx, exporterCreateSettings, expCfg)
	if err != nil {
		log.Error("error on creating metrics exporter", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = exporter.Start(exporterCtx, nil)
	if err != nil {
		log.Error("error on starting metrics exporter", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	transformFactory := transformprocessor.NewFactory()
	transformCtx := context.WithValue(otelContext, "component", "transformprocessor")
	log.Info("start to create metrics component", zap.Any("component", transformCtx.Value("component")))
	transformCfg := transformFactory.CreateDefaultConfig().(*transformprocessor.Config)
	transformSet := processor.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	processor, err := transformFactory.CreateMetricsProcessor(transformCtx, transformSet, transformCfg, exporter)
	if err != nil {
		log.Error("error on creating metrics processor", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	log.Info("created kafka metrics exporter successfully")
	// receiver Factory
	orbReceiverFactory := orbreceiver.NewFactory()
	receiverCtx := context.WithValue(otelContext, "component", "orbreceiver")
	receiverCfg := orbReceiverFactory.CreateDefaultConfig().(*orbreceiver.Config)
	receiverCfg.Logger = logger
	receiverCfg.PubSub = pubSub
	receiverCfg.SinkerService = bridgeService
	receiverSet := receiver.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	receiver, err := orbReceiverFactory.CreateMetricsReceiver(receiverCtx, receiverSet, receiverCfg, processor)
	log.Info("created metrics receiver")
	if err != nil {
		log.Error("error on creating metrics receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = receiver.Start(receiverCtx, nil)
	log.Info("started receiver")
	if err != nil {
		log.Error("error on starting metrics receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	log.Info("created orb metrics receiver successfully")
	return otelCancelFunc, nil
}

func StartOtelLogsComponents(ctx context.Context, bridgeService *bridgeservice.SinkerOtelBridgeService, logger *zap.Logger, kafkaUrl string, pubSub mfnats.PubSub) (context.CancelFunc, error) {
	otelContext, otelCancelFunc := context.WithCancel(ctx)

	log := logger.Sugar()
	log.Info("Starting to create Otel Logs Components in routine: ", ctx.Value("routine"))
	exporterFactory := kafkaexporter.NewFactory()
	exporterCtx := context.WithValue(otelContext, "component", "kafkaexporterlogs")
	exporterCreateSettings := exporter.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	expCfg := exporterFactory.CreateDefaultConfig().(*kafkaexporter.Config)
	expCfg.Brokers = []string{kafkaUrl}
	expCfg.Topic = "otlp_logs"
	exporter, err := exporterFactory.CreateLogsExporter(exporterCtx, exporterCreateSettings, expCfg)
	if err != nil {
		log.Error("error on creating logs exporter", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = exporter.Start(exporterCtx, nil)
	if err != nil {
		log.Error("error on starting logs exporter", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	transformFactory := transformprocessor.NewFactory()
	transformCtx := context.WithValue(otelContext, "component", "transformprocessorlogs")
	log.Info("start to create logs component", zap.Any("component", transformCtx.Value("component")))
	transformCfg := transformFactory.CreateDefaultConfig().(*transformprocessor.Config)
	transformSet := processor.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	processor, err := transformFactory.CreateLogsProcessor(transformCtx, transformSet, transformCfg, exporter)
	if err != nil {
		log.Error("error on creating logs processor", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	log.Info("created kafka logs exporter successfully")
	// receiver Factory
	orbReceiverFactory := orbreceiver.NewFactory()
	receiverCtx := context.WithValue(otelContext, "component", "orbreceiverlogs")
	receiverCfg := orbReceiverFactory.CreateDefaultConfig().(*orbreceiver.Config)
	receiverCfg.Logger = logger
	receiverCfg.PubSub = pubSub
	receiverCfg.SinkerService = bridgeService
	receiverSet := receiver.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	receiver, err := orbReceiverFactory.CreateLogsReceiver(receiverCtx, receiverSet, receiverCfg, processor)
	log.Info("created logs receiver")
	if err != nil {
		log.Error("error on creating receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = receiver.Start(receiverCtx, nil)
	log.Info("started logs receiver")
	if err != nil {
		log.Error("error on starting logs receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	log.Info("created orb logs receiver successfully")
	return otelCancelFunc, nil
}

func StartOtelTracesComponents(ctx context.Context, bridgeService *bridgeservice.SinkerOtelBridgeService, logger *zap.Logger, kafkaUrl string, pubSub mfnats.PubSub) (context.CancelFunc, error) {
	otelContext, otelCancelFunc := context.WithCancel(ctx)

	log := logger.Sugar()
	log.Info("Starting to create Otel Traces Components in routine: ", ctx.Value("routine"))
	exporterFactory := kafkaexporter.NewFactory()
	exporterCtx := context.WithValue(otelContext, "component", "kafkaexporter")
	exporterCreateSettings := exporter.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	expCfg := exporterFactory.CreateDefaultConfig().(*kafkaexporter.Config)
	expCfg.Brokers = []string{kafkaUrl}
	expCfg.Topic = "otlp_traces"
	exporter, err := exporterFactory.CreateTracesExporter(exporterCtx, exporterCreateSettings, expCfg)
	if err != nil {
		log.Error("error on creating traces exporter", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = exporter.Start(exporterCtx, nil)
	if err != nil {
		log.Error("error on starting traces exporter", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	transformFactory := transformprocessor.NewFactory()
	transformCtx := context.WithValue(otelContext, "component", "transformprocessor")
	log.Info("start to create traces component", zap.Any("component", transformCtx.Value("component")))
	transformCfg := transformFactory.CreateDefaultConfig().(*transformprocessor.Config)
	transformSet := processor.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	processor, err := transformFactory.CreateTracesProcessor(transformCtx, transformSet, transformCfg, exporter)
	if err != nil {
		log.Error("error on creating traces processor", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	log.Info("created kafka traces exporter successfully")
	// receiver Factory
	orbReceiverFactory := orbreceiver.NewFactory()
	receiverCtx := context.WithValue(otelContext, "component", "orbreceiver")
	receiverCfg := orbReceiverFactory.CreateDefaultConfig().(*orbreceiver.Config)
	receiverCfg.Logger = logger
	receiverCfg.PubSub = pubSub
	receiverCfg.SinkerService = bridgeService
	receiverSet := receiver.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	receiver, err := orbReceiverFactory.CreateTracesReceiver(receiverCtx, receiverSet, receiverCfg, processor)
	log.Info("created traces receiver")
	if err != nil {
		log.Error("error on creating traces receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = receiver.Start(receiverCtx, nil)
	log.Info("started traces receiver")
	if err != nil {
		log.Error("error on starting traces receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	log.Info("created orb traces receiver successfully")
	return otelCancelFunc, nil
}
