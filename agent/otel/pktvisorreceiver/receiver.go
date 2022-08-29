package pktvisorreceiver

import (
	"context"
	"errors"
	"fmt"
	configutil "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
)

type prometheusReceiverWrapper struct {
	logger             *zap.Logger
	params             component.ReceiverCreateSettings
	config             *Config
	consumer           consumer.Metrics
	prometheusReceiver component.MetricsReceiver
}

// New returns a prometheusReceiverWrapper
func New(params component.ReceiverCreateSettings, cfg *Config, consumer consumer.Metrics) *prometheusReceiverWrapper {
	var logger *zap.Logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		logger.Error("failed to create logger", zap.Error(err))
	}
	return &prometheusReceiverWrapper{params: params, config: cfg, consumer: consumer, logger: logger}
}

// Start creates and starts the prometheus receiver.
func (prw *prometheusReceiverWrapper) Start(ctx context.Context, host component.Host) error {
	pFactory := NewFactory()

	pConfig, err := GetPrometheusConfig(prw.config)
	if err != nil {
		return fmt.Errorf("failed to create prometheus receiver config: %v", err)
	}

	pr, err := pFactory.CreateMetricsReceiver(ctx, prw.params, pConfig, prw.consumer)
	if err != nil {
		return fmt.Errorf("failed to create prometheus receiver: %v", err)
	}

	prw.prometheusReceiver = pr
	return prw.prometheusReceiver.Start(ctx, host)
}

func GetPrometheusConfig(cfg *Config) (*Config, error) {
	var bearerToken string
	// TODO check what was UseServiceAccount field
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	bearerToken = restConfig.BearerToken
	if bearerToken == "" {
		return nil, errors.New("bearer token was empty")
	}

	out := &Config{}
	httpConfig := configutil.HTTPClientConfig{}

	scheme := "http"

	httpConfig.BearerToken = configutil.Secret(bearerToken)

	scrapeConfig := &config.ScrapeConfig{
		ScrapeInterval:  model.Duration(cfg.BufferPeriod),
		ScrapeTimeout:   model.Duration(cfg.BufferPeriod),
		JobName:         fmt.Sprintf("%s/%s", typeStr, cfg.Endpoint),
		HonorTimestamps: true,
		Scheme:          scheme,
		MetricsPath:     cfg.MetricsPath,
		Params:          cfg.Params,
		ServiceDiscoveryConfigs: discovery.Configs{
			&discovery.StaticConfig{
				{
					Targets: []model.LabelSet{
						{model.AddressLabel: model.LabelValue(cfg.Endpoint)},
					},
				},
			},
		},
	}

	scrapeConfig.HTTPClientConfig = httpConfig
	out.PrometheusConfig = &config.Config{ScrapeConfigs: []*config.ScrapeConfig{
		scrapeConfig,
	}}

	return out, nil
}

// Shutdown stops the underlying Prometheus receiver.
func (prw *prometheusReceiverWrapper) Shutdown(ctx context.Context) error {
	return prw.prometheusReceiver.Shutdown(ctx)
}
