package otel

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
	"go.uber.org/zap"
)

func StartOtelComponents(ctx context.Context, logger zap.Logger) error {
	log := logger.Sugar()
	log.Info("Starting to create Otel Components", ctx.Value("routine"))
	var bla kafkaexporter.Config
	log.Info("load info on", bla)
	return nil
}
