package components

import (
	"context"
	"github.com/ns1labs/orb/pkg/config"
	"go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
)

func CreateOTLPSinkerReceiver(ctx context.Context,
	factories component.Factories,
	sinkerGrpcCfg config.GRPCConfig,
	mc consumer.Metrics) (r component.Receiver, err error) {
	logger := config.LoggerFromContext(ctx)
	slog := logger.Sugar()
	slog.Debug("create otlpreceiver from orb-sinker", sinkerGrpcCfg)
	name := "otlpreceiver"
	// Receiver Context
	receiverCtx := context.WithValue(ctx, "id", name)
	factory := factories.Receivers[otelconfig.Type(name)]

	cfg := factory.CreateDefaultConfig().(*otlpreceiver.Config)
	cfg.SetIDName(name)
	cfg.GRPC.NetAddr.Endpoint = sinkerGrpcCfg.GetEndpoint()
	cfg.HTTP = nil

	set := component.ReceiverCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.BuildInfo{},
	}
	r, err = factory.CreateMetricsReceiver(receiverCtx, set, cfg, mc)
	if err != nil {
		slog.Error("error during CreateMetricsReceiver", err)
	}
	return
}
