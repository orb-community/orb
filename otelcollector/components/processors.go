package components

import (
	"context"
	"errors"
	"github.com/ns1labs/orb/pkg/config"

	"github.com/ns1labs/orb/otelcollector/components/orbattributesprocessor"
	"go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
)

func GetAttributeProcessorWithOwnerAndSinkData(ctx context.Context, factories component.Factories, nextConsumer consumer.Metrics) (component.MetricsProcessor, error) {
	// ensure owner and Sink attribute to be in context
	if ctx.Value("sinkData") == nil {
		return nil, errors.New("data must contain sinkData")
	}
	if ctx.Value("ownerID") == nil {
		return nil, errors.New("data must contain ownerID")
	}
	logger := config.LoggerFromContext(ctx)
	slog := logger.Sugar()
	name := "attributesprocessor"
	subCtx := context.WithValue(ctx, "name", name)
	slog.Debug("create processor:", name)
	factory := factories.Processors[otelconfig.Type(name)]
	cfg := factory.CreateDefaultConfig().(*orbattributesprocessor.Config)
	cfg.AddUpsertActionFromContext("orb.sinkData", "sinkData")
	cfg.AddUpsertActionFromContext("orb.ownerID", "ownerId")
	set := component.ProcessorCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.BuildInfo{},
	}
	return factory.CreateMetricsProcessor(subCtx, set, cfg, nextConsumer)

}
