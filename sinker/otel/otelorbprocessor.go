package otel

import (
	"github.com/ns1labs/orb/pkg/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
)

type OrbProcessorService interface {
	AddSinksInfo(agentId, ownerId, datasetId string) (sinksInfo []types.Metadata)
}

type otelProcessor struct {
	exporterConfig       otlpexporter.Config
	exporterFactory      component.ExporterFactory
	attributeProcConfig  attributesprocessor.Config
	attributeProcFactory component.ProcessorFactory
	receiverConfig       redisreceiver.Config
	receiverFactory      component.ReceiverFactory
}
