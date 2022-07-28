package config

import (
	"context"
	"go.uber.org/zap"
)

// ContextWithLogger encapsulates logger within context
func ContextWithLogger(parent context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(parent, "logger", logger)
}

// LoggerFromContext extracts logger from context
func LoggerFromContext(ctx context.Context) *zap.Logger {
	return ctx.Value("logger").(*zap.Logger)
}
