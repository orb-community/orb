package otlpexporter

import (
	"github.com/orb-community/orb/sinks/backend"
)

var _ backend.Backend = (*Backend)(nil)

type Backend struct {
	apiHost     string
	apiPort     uint64
	apiUser     string
	apiPassword string
}

type SinkFeature struct {
	Backend     string                  `json:"backend"`
	Description string                  `json:"description"`
	Config      []backend.ConfigFeature `json:"config"`
}

func (b Backend) Metadata() interface{} {
	return SinkFeature{
		Backend:     "otlpexporter",
		Description: "OTLP gRPC Exporter, configuration is documented in https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md",
		Config:      b.CreateFeatureConfig(),
	}
}
