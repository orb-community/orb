package pktvisorreceiver

import (
	"context"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
)

// This file implements factory for prometheus_simple receiver
const (
	// The value of "type" key in configuration.
	typeStr = "pktvisor_prometheus"

	defaultEndpoint    = "localhost:10853"
	defaultMetricsPath = "/api/v1/policies/__all/metrics/prometheus"
)

var defaultCollectionInterval = 60 * time.Second

// NewFactory creates a factory for "Simple" Prometheus receiver.
func NewFactory() component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		CreateDefaultConfig,
		component.WithMetricsReceiver(CreateMetricsReceiver))
}

func CreateDefaultSettings(logger *zap.Logger) component.ReceiverCreateSettings {
	return component.ReceiverCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
}

func CreateDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
		TCPAddr: confignet.TCPAddr{
			Endpoint: defaultEndpoint,
		},
		MetricsPath:        defaultMetricsPath,
		CollectionInterval: defaultCollectionInterval,
	}
}

func CreateMetricsReceiver(
	_ context.Context,
	params component.ReceiverCreateSettings,
	cfg config.Receiver,
	nextConsumer consumer.Metrics,
) (component.MetricsReceiver, error) {
	rCfg := cfg.(*Config)
	return New(params, rCfg, nextConsumer), nil
}
