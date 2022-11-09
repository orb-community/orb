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
		`set(attributes["agent-name"], resource.attributes["agent_name"])`,
		`set(attributes["agent-tags"], resource.attributes["agent_tags"])`,
		`set(attributes["orb-tags"], resource.attributes["orb_tags"])`,
		`set(attributes["agent-groups"], resource.attributes["agent_groups"])`,
		`set(attributes["agent-ownerID"], resource.attributes["agent_ownerID"])`,
		`set(attributes["policy-id"], resource.attributes["policy_id"])`,
		`set(attributes["policy-name"], resource.attributes["policy_name"])`,
		`set(attributes["sink-id"], resource.attributes["sink_id"])`,
		`set(attributes["format"], "otlp")`,
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
