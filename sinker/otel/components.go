package otel

import (
	"context"
	"go.uber.org/zap"
)

func StartOtelComponents(ctx context.Context, logger zap.Logger) error {
	log := logger.Sugar()
	log.Info("Starting to create Otel Components", ctx.Value("routine"))
	return nil
}
