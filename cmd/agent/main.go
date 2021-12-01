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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	otelconfig "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
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

	defaultEndpoint    = "http://localhost:10853"
	defaultMetricsPath = "/api/v1/policies/__all/metrics/bucket/1"
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

	// new otel receiver
	factory := otel.NewFactory()
	var getConfigFn getReceiverConfigFn
	if getConfigFn == nil {
		getConfigFn = func() otelconfig.Receiver {
			config := &otel.Config{
				ReceiverSettings: otelconfig.NewReceiverSettings(otelconfig.NewComponentID(typeStr)),
				TCPAddr: confignet.TCPAddr{
					Endpoint: defaultEndpoint,
				},
				MetricsPath:        defaultMetricsPath,
				CollectionInterval: defaultCollectionInterval,
			}
			return config
		}
	}
	ctx := context.Background()
	wrap := createMetricsReceiver(factory)
	receiverCreateSet := componenttest.NewNopReceiverCreateSettings()
	receiver, err := wrap(ctx, receiverCreateSet, getConfigFn(), nil, logger)

	// handle signals
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		a.Stop()
		receiver.Shutdown(ctx)
		done <- true
	}()

	// start agent
	err = a.Start()
	if err != nil {
		logger.Error("agent startup error", zap.Error(err))
		os.Exit(1)
	}

	err = receiver.Start(context.Background(), nil)
	if err != nil {
		logger.Error("otel startup error", zap.Error(err))
		os.Exit(1)
	}

	<-done
}

type getReceiverConfigFn func() otelconfig.Receiver

type createReceiverFn func(
	ctx context.Context,
	params component.ReceiverCreateSettings,
	cfg otelconfig.Receiver,
	nextConsumer consumer.Metrics,
	logger *zap.Logger) (component.Receiver, error)

func createMetricsReceiver(factory component.ReceiverFactory) createReceiverFn {
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
