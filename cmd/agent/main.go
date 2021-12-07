/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"context"
	"fmt"
	"github.com/ns1labs/orb/agent"
	"github.com/ns1labs/orb/agent/backend/pktvisor"
	config2 "github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/agent/otel"
	"github.com/ns1labs/orb/buildinfo"
	"github.com/ns1labs/orb/pkg/errors"
	promexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	configutil "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	defaultConfig = "/etc/orb/agent.yaml"
	typeStr       = "prometheus_simple"

	defaultEndpoint    = "localhost:10853"
	defaultMetricsPath = "/api/v1/policies/__all/metrics/prometheus"
)

var (
	cfgFiles                  []string
	Debug                     bool
	rng                       = rand.New(rand.NewSource(time.Now().UnixNano()))
	defaultCollectionInterval = 10 * time.Second
)

func init() {

	pktvisor.Register()

}

func Version(cmd *cobra.Command, args []string) {
	fmt.Printf("orb-agent %s\n", buildinfo.GetVersion())
	os.Exit(0)
}

func Run(cmd *cobra.Command, args []string) {

	// logger
	var logger *zap.Logger
	var err error
	if Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	cobra.CheckErr(err)

	initConfig()

	// configuration
	var config config2.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Error("agent start up error (config)", zap.Error(err))
		os.Exit(1)
	}

	config.Debug = Debug
	// include pktvisor backend by default if binary is at default location
	_, err = os.Stat(pktvisor.DefaultBinary)
	if err == nil && config.OrbAgent.Backends == nil {
		config.OrbAgent.Backends = make(map[string]map[string]string)
		config.OrbAgent.Backends["pktvisor"] = make(map[string]string)
		config.OrbAgent.Backends["pktvisor"]["binary"] = pktvisor.DefaultBinary
		if len(cfgFiles) > 0 {
			config.OrbAgent.Backends["pktvisor"]["config_file"] = cfgFiles[0]
		}
	}

	// new agent
	a, err := agent.New(logger, config)
	if err != nil {
		logger.Error("agent start up error", zap.Error(err))
		os.Exit(1)
	}

	ctx := context.Background()

	// new otel receiver
	//factory := otel.NewFactory()
	exporter, err := createNewExporter(ctx, logger)
	receiver, err := createNewReceiver(ctx, exporter, logger)

	// handle signals
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		a.Stop()
		exporter.Shutdown(ctx)
		receiver.Shutdown(ctx)
		done <- true
	}()

	// start agent
	err = a.Start()
	if err != nil {
		logger.Error("agent startup error", zap.Error(err))
		os.Exit(1)
	}

	err = exporter.Start(ctx, nil)
	if err != nil {
		logger.Error("otel exporter startup error", zap.Error(err))
		os.Exit(1)
	}

	err = receiver.Start(ctx, nil)
	if err != nil {
		logger.Error("otel receiver startup error", zap.Error(err))
		os.Exit(1)
	}

	<-done
}

func createNewExporter(ctx context.Context, logger *zap.Logger) (component.MetricsExporter, error) {
	// 2. Create the Prometheus metrics exporter that'll receive and verify the metrics produced.
	exporterCfg := &promexporter.Config{
		ExporterSettings: otelconfig.NewExporterSettings(otelconfig.NewComponentID(typeStr)),
		Namespace:        "test",
		Endpoint:         ":8787",
		SendTimestamps:   true,
		MetricExpiration: 2 * time.Hour,
	}
	exporterFactory := promexporter.NewFactory()
	set := component.ExporterCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.GetMeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	exporter, err := exporterFactory.CreateMetricsExporter(ctx, set, exporterCfg)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func createNewReceiver(ctx context.Context, exporter component.MetricsExporter, logger *zap.Logger) (component.MetricsReceiver, error) {
	receiverFactory := prometheusreceiver.NewFactory()
	receiverCreateSet := component.ReceiverCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.GetMeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	rcvCfg := &otel.Config{
		ReceiverSettings: otelconfig.NewReceiverSettings(otelconfig.NewComponentID(typeStr)),
		TCPAddr: confignet.TCPAddr{
			Endpoint: defaultEndpoint,
		},
		MetricsPath:        defaultMetricsPath,
		CollectionInterval: defaultCollectionInterval,
	}
	// 3.5 Create the Prometheus receiver and pass in the preivously created Prometheus exporter.
	pConfig, err := getPrometheusConfig(rcvCfg)
	if err != nil {
		return nil, errors.Wrap(errors.New("failed to create prometheus receiver config"), err)
	}
	prometheusReceiver, err := receiverFactory.CreateMetricsReceiver(ctx, receiverCreateSet, pConfig, exporter)
	if err != nil {
		return nil, err
	}
	return prometheusReceiver, nil
}

func getPrometheusConfig(cfg *otel.Config) (*prometheusreceiver.Config, error) {
	var bearerToken string
	if cfg.UseServiceAccount {
		restConfig, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		bearerToken = restConfig.BearerToken
		if bearerToken == "" {
			return nil, errors.New("bearer token was empty")
		}
	}

	out := &prometheusreceiver.Config{}
	httpConfig := configutil.HTTPClientConfig{}

	scheme := "http"

	if cfg.TLSEnabled {
		scheme = "https"
		httpConfig.TLSConfig = configutil.TLSConfig{
			CAFile:             cfg.TLSConfig.CAFile,
			CertFile:           cfg.TLSConfig.CertFile,
			KeyFile:            cfg.TLSConfig.KeyFile,
			InsecureSkipVerify: cfg.TLSConfig.InsecureSkipVerify,
		}
	}

	httpConfig.BearerToken = configutil.Secret(bearerToken)

	scrapeConfig := &config.ScrapeConfig{
		ScrapeInterval:  model.Duration(cfg.CollectionInterval),
		ScrapeTimeout:   model.Duration(cfg.CollectionInterval),
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

type getReceiverConfigFn func() otelconfig.Receiver

type createReceiverFn func(
	ctx context.Context,
	params component.ReceiverCreateSettings,
	cfg otelconfig.Receiver,
	nextConsumer consumer.Metrics,
	logger *zap.Logger) (component.Receiver, error)

func createMetricsReceiver() createReceiverFn {
	return func(ctx context.Context, params component.ReceiverCreateSettings, cfg otelconfig.Receiver, nextConsumer consumer.Metrics, logger *zap.Logger) (component.Receiver, error) {
		rCfg := cfg.(*otel.Config)
		//return factory.CreateMetricsReceiver(ctx, params, cfg, nextConsumer)
		return otel.New(params, rCfg, nextConsumer), nil
	}
}

func mergeOrError(path string) {

	v := viper.New()
	if len(path) > 0 {
		v.SetConfigFile(path)
		v.SetConfigType("yaml")
	}

	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	// note: viper seems to require a default (or a BindEnv) to be overridden by environment variables
	v.SetDefault("orb.cloud.api.address", "https://api.orb.live")
	v.SetDefault("orb.cloud.api.token", "")
	v.SetDefault("orb.cloud.config.agent_name", "")
	v.SetDefault("orb.cloud.config.auto_provision", true)
	v.SetDefault("orb.cloud.mqtt.address", "tls://mqtt.orb.live:8883")
	v.SetDefault("orb.cloud.mqtt.id", "")
	v.SetDefault("orb.cloud.mqtt.key", "")
	v.SetDefault("orb.cloud.mqtt.channel_id", "")
	v.SetDefault("orb.db.file", "./orb-agent.db")
	v.SetDefault("orb.tls.verify", true)

	v.SetDefault("orb.backends.pktvisor.binary", "")
	v.SetDefault("orb.backends.pktvisor.config_file", "")
	v.SetDefault("orb.backends.pktvisor.api_host", "localhost")
	v.SetDefault("orb.backends.pktvisor.api_port", "10853")

	if len(path) > 0 {
		cobra.CheckErr(v.ReadInConfig())
	}

	var fZero float64

	// check that version of config files are all matched up
	if versionNumber1 := viper.GetFloat64("version"); versionNumber1 != fZero {
		versionNumber2 := v.GetFloat64("version")
		if versionNumber2 == fZero {
			cobra.CheckErr("Failed to parse config vesrion in: " + path)
		}
		if versionNumber2 != versionNumber1 {
			cobra.CheckErr("Config file version mismatch in: " + path)
		}
	}

	cobra.CheckErr(viper.MergeConfigMap(v.AllSettings()))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// set defaults first
	mergeOrError("")
	if len(cfgFiles) == 0 {
		if _, err := os.Stat(defaultConfig); !os.IsNotExist(err) {
			mergeOrError(defaultConfig)
		}
	} else {
		for _, conf := range cfgFiles {
			mergeOrError(conf)
		}
	}
}

func main() {

	rootCmd := &cobra.Command{
		Use: "orb-agent",
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show agent version",
		Run:   Version,
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run orb-agent and connect to Orb control plane",
		Long:  `Run orb-agent and connect to Orb control plane`,
		Run:   Run,
	}

	runCmd.Flags().StringSliceVarP(&cfgFiles, "config", "c", []string{}, "Path to config files (may be specified multiple times)")
	runCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Enable verbose (debug level) output")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}
