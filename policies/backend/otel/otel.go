package otel

import (
	"errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/policies/backend"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"strings"
)

type otelBackend struct {
	version string
	logger  *zap.Logger
}

func (o otelBackend) SupportsFormat(format string) bool {
	if strings.EqualFold(format, "yaml") {
		return true
	}
	return false
}

func (o otelBackend) ConvertFromFormat(format string, policy string) (metadata types.Metadata, err error) {
	if !o.SupportsFormat(format) {
		return nil, errors.New("unsupported format")
	}
	o.logger.Info("converting policy from yaml", zap.String("policy", policy), zap.String("format", format))
	err = yaml.Unmarshal([]byte(policy), &metadata)
	return
}

// Validate Will not validate anything until we have a better way to do this
func (o otelBackend) Validate(_ types.Metadata) error {
	// block everything related to the exporter tag, this is not supported
	return nil
}

func Register(logger *zap.Logger) bool {
	l := logger.Named("otel-backend")
	backend.Register("otel", &otelBackend{logger: l})
	return true
}
