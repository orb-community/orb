package otel

import (
	"github.com/ns1labs/orb/pkg/types"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
)

type OrbProcessorService interface {
	AddSinksInfo(agentId, ownerId, datasetId string) (sinksInfo []types.Metadata)
}

type otelProcessor struct {
	exporterConfig       otlpexporter.Config
	exporterFactory      component.ExporterFactory
	attributeProcConfig  attributeprocessor.Config
	attributeProcFactory component.ProcessorFactory
}
