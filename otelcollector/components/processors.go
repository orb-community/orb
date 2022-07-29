package components

import (
	"context"
	"github.com/ns1labs/orb/pkg/config"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/attraction"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
)

func GetAttributeProcessorWithOwnerAndSinkData(ctx context.Context, factories component.Factories) error {
	logger := config.LoggerFromContext(ctx)
	slog := logger.Sugar()
	name := "attributeprocessor"
	_ = context.WithValue(ctx, "name", name)
	slog.Debug("create processor", name, "")
	factory := factories.Processors[otelconfig.Type(name)]
	cfg := factory.CreateDefaultConfig().(*attributesprocessor.Config)
	cfg.Actions = []attraction.ActionKeyValue{
		{Key: "ownerID", Value: "???", FromAttribute: "???", Action: "insert"}, // Mainflux Owner ID
		{Key: "agentID", Value: "???", FromAttribute: "???", Action: "insert"}, // Agent ID
		{Key: "sinks", Value: "???", FromAttribute: "???", Action: "insert"},   // Sink slice with id and config metadata
	}

	return nil
}
