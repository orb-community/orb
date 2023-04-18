package pktvisorreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		CreateDefaultConfig,
		receiver.WithMetrics(CreateMetricsReceiver, component.StabilityLevelAlpha))
}

func CreateDefaultSettings(logger *zap.Logger) receiver.CreateSettings {
	return receiver.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
}

func CreateReceiverConfig(endpoint, metricsPath string) component.Config {
	return &Config{
		TCPAddr: confignet.TCPAddr{
			Endpoint: endpoint,
		},
		MetricsPath:        metricsPath,
		CollectionInterval: defaultCollectionInterval,
	}
}

func CreateDefaultConfig() component.Config {
	return &Config{
		TCPAddr: confignet.TCPAddr{
			Endpoint: defaultEndpoint,
		},
		MetricsPath:        defaultMetricsPath,
		CollectionInterval: defaultCollectionInterval,
	}
}

func CreateMetricsReceiver(
	_ context.Context,
	params receiver.CreateSettings,
	config component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	rCfg := config.(*Config)
	return New(params, rCfg, consumer), nil
}
