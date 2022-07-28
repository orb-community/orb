package components

import (
	"context"
	"github.com/ns1labs/orb/pkg/config"
	"go.opentelemetry.io/collector/component"
)

func GetAttributeProcessorWithOwnerAndSinkData(ctx context.Context, factories component.Factories) error {
	logger := config.LoggerFromContext(ctx)
	slog := logger.Sugar()
	logger.Debug("create attribute processor")
	name := "attributeprocessor"
	factories[""]

	return nil
}
