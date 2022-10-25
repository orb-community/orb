package otel

import (
	"context"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/sinker/otel/bridgeservice"
	kafkaexporter "github.com/ns1labs/orb/sinker/otel/kafkafanoutexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"

	"github.com/ns1labs/orb/sinker/otel/orbreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func StartOtelComponents(ctx context.Context, bridgeService *bridgeservice.SinkerOtelBridgeService, logger *zap.Logger, kafkaUrl string, pubSub mfnats.PubSub) (context.CancelFunc, error) {
	otelContext, otelCancelFunc := context.WithCancel(ctx)

	log := logger.Sugar()
	log.Info("Starting to create Otel Components in routine: ", ctx.Value("routine"))
	exporterFactory := kafkaexporter.NewFactory()
	exporterCtx := context.WithValue(otelContext, "component", "kafkaexporter")
	exporterCreateSettings := component.ExporterCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	expCfg := exporterFactory.CreateDefaultConfig().(*kafkaexporter.Config)
	expCfg.Brokers = []string{kafkaUrl}
	expCfg.Topic = "otlp_metrics"
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
	transformFactory := transformprocessor.NewFactory()
	transformCtx := context.WithValue(otelContext, "component", "transformprocessor")
	log.Info("start to create component", zap.Any("component", transformCtx.Value("component")))
	transformCfg := transformFactory.CreateDefaultConfig().(*transformprocessor.Config)
	transformCfg.OTTLConfig.Metrics.Statements = []string{
		`set(resource.attributes["agent-name"], ctx.Value("agent-name"))`,
		`set(resource.attributes["agent-tags"], ctx.Value("agent-tags"))`,
		`set(resource.attributes["orb-tags"], ctx.Value("orb-tags"))`,
		`set(resource.attributes["agent-groups"], ctx.Value("agent-groups"))`,
		`set(resource.attributes["agent-ownerID"], ctx.Value("agent-ownerID"))`,
		`set(resource.attributes["policy-id"], ctx.Value("policy-id"))`,
		`set(resource.attributes["policy-name"], ctx.Value("policy-name"))`,
		`set(resource.attributes["sink-id"], ctx.Value("sink-id"))`,
		`set(resource.attributes["format"], "otlp")`,
	}
	transformSet := component.ProcessorCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	processor, err := transformFactory.CreateMetricsProcessor(transformCtx, transformSet, transformCfg, exporter)
	if err != nil {
		log.Error("error on creating processor", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	log.Info("created kafka exporter successfully")
	// receiver Factory
	orbReceiverFactory := orbreceiver.NewFactory()
	receiverCtx := context.WithValue(otelContext, "component", "orbreceiver")
	receiverCfg := orbReceiverFactory.CreateDefaultConfig().(*orbreceiver.Config)
	receiverCfg.Logger = logger
	receiverCfg.PubSub = pubSub
	receiverCfg.SinkerService = bridgeService
	receiverSet := component.ReceiverCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
	}
	receiver, err := orbReceiverFactory.CreateMetricsReceiver(receiverCtx, receiverSet, receiverCfg, processor)
	log.Info("created receiver")
	if err != nil {
		log.Error("error on creating receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	err = receiver.Start(receiverCtx, nil)
	log.Info("started receiver")
	if err != nil {
		log.Error("error on starting receiver", err)
		otelCancelFunc()
		ctx.Done()
		return nil, err
	}
	log.Info("created orb receiver successfully")
	return otelCancelFunc, nil
}
