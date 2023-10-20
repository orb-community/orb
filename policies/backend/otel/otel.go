package otel

import (
	"errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/policies/backend"
	"gopkg.in/yaml.v3"
	"strings"
)

type otelBackend struct {
	version string
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
	err = yaml.Unmarshal([]byte(policy), &metadata)
	return
}

// Validate Will not validate anything until we have a better way to do this
func (o otelBackend) Validate(_ types.Metadata) error {
	// block everything related to the exporter tag, this is not supported
	return nil
}

func Register() bool {
	backend.Register("otel", &otelBackend{})
	return true
}
