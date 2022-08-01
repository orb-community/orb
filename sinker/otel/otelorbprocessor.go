package otel

import (
	"github.com/ns1labs/orb/pkg/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"go.opentelemetry.io/collector/component"
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
