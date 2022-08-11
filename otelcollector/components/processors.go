package components

import (
	"context"
	"github.com/ns1labs/orb/pkg/config"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/attraction"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
)

func GetAttributeProcessorWithOwnerAndSinkData(ctx context.Context, factories component.Factories, nextConsumer consumer.Metrics) error {
	logger := config.LoggerFromContext(ctx)
	slog := logger.Sugar()
	name := "attributeprocessor"
	subCtx := context.WithValue(ctx, "name", name)
	slog.Debug("create processor:", name)
	factory := factories.Processors[otelconfig.Type(name)]
	cfg := factory.CreateDefaultConfig().(*attributesprocessor.Config)

	cfg.Actions = []attraction.ActionKeyValue{
		{Key: "orb_private.ownerID", FromContext: "ownerID", Action: "insert"}, // Mainflux Owner ID
		{Key: "orb_private.agentID", FromContext: "agentID", Action: "insert"}, // Agent ID
		{Key: "orb_private.sinks", FromContext: "sinks", Action: "insert"},     // Sink slice with id and config metadata
	}
	set := component.ProcessorCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.BuildInfo{},
	}
	factory.CreateMetricsProcessor(subCtx, set, cfg, nextConsumer)

	return nil
}
