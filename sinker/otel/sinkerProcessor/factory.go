package otel

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr         = "otlp"
	defaultEndpoint = "localhost:4317"
)

// NewFactory returns a new factory
func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		CreateDefaultConfig,
		component.WithMetricsProcessor(CreateMetricsProcessor))
}

func CreateDefaultConfig() config.Processor {
	return &Config{
		GRPCClientSettings: configgrpc.GRPCClientSettings{
			Endpoint:    defaultEndpoint,
			Compression: "",
			TLSSetting: configtls.TLSClientSetting{
				Insecure: true,
			},
			Keepalive:       nil,
			ReadBufferSize:  0,
			WriteBufferSize: 512 * 1024,
			WaitForReady:    false,
			Headers:         map[string]string{},
			BalancerName:    "",
		},
	}
}

func CreateMetricsProcessor(
	_ context.Context,
	params component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Metrics,
) (component.MetricsProcessor, error) {
	return processorhelper.NewMetricsProcessor(
		cfg,
		nextConsumer,
		newProcessor(cfg).processMetrics,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
	)
}
