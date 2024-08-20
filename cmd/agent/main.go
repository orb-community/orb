/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"context"
	"fmt"
	"github.com/orb-community/orb/agent/backend/otel"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/orb-community/orb/agent"
	"github.com/orb-community/orb/agent/backend/pktvisor"
	"github.com/orb-community/orb/agent/config"
	"github.com/orb-community/orb/buildinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultConfig = "/opt/orb/agent.yaml"
)

var (
	cfgFiles []string
	Debug    bool
)

func init() {
	pktvisor.Register()
	otel.Register()
}

func Version(_ *cobra.Command, _ []string) {
	fmt.Printf("orb-agent %s\n", buildinfo.GetVersion())
	os.Exit(0)
}

func Run(_ *cobra.Command, _ []string) {

	initConfig()

	// configuration
	var configData config.Config
	err := viper.Unmarshal(&configData)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("agent start up error (configData): %w", err))
		os.Exit(1)
	}

	// logger
	var logger *zap.Logger
	atomicLevel := zap.NewAtomicLevel()
	if Debug {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else {
		atomicLevel.SetLevel(zap.InfoLevel)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		os.Stdout,
		atomicLevel,
	)
	logger = zap.New(core, zap.AddCaller())
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	// include pktvisor backend by default if binary is at default location
	_, err = os.Stat(pktvisor.DefaultBinary)
	logger.Info("backends loaded", zap.Any("backends", configData.OrbAgent.Backends))
	if err == nil && configData.OrbAgent.Backends == nil {
		logger.Info("no backends loaded, adding pktvisor as default")
		configData.OrbAgent.Backends = make(map[string]map[string]string)
		configData.OrbAgent.Backends["pktvisor"] = make(map[string]string)
		configData.OrbAgent.Backends["pktvisor"]["binary"] = pktvisor.DefaultBinary
		configData.OrbAgent.Backends["pktvisor"]["api_host"] = "localhost"
		if _, ok := configData.OrbAgent.Backends["pktvisor"]["api_port"]; !ok {
			configData.OrbAgent.Backends["pktvisor"]["api_port"] = "10853"
		}
		if len(cfgFiles) > 0 {
			configData.OrbAgent.Backends["pktvisor"]["config_file"] = cfgFiles[0]
		}
	}

	// new agent
	a, err := agent.New(logger, configData)
	if err != nil {
		logger.Error("agent start up error", zap.Error(err))
		os.Exit(1)
	}

	// handle signals
	done := make(chan bool, 1)
	rootCtx, cancelFunc := context.WithCancel(context.WithValue(context.Background(), "routine", "mainRoutine"))

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		select {
		case <-sigs:
			logger.Warn("stop signal received stopping agent")
			a.Stop(rootCtx)
			cancelFunc()
		case <-rootCtx.Done():
			logger.Warn("mainRoutine context cancelled")
			done <- true
			return
		}
	}()

	// start agent
	err = a.Start(rootCtx, cancelFunc)
	if err != nil {
		logger.Error("agent startup error", zap.Error(err))
		os.Exit(1)
	}

	<-done
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
	v.SetDefault("orb.cloud.api.address", "https://orb.live")
	v.SetDefault("orb.cloud.api.token", "")
	v.SetDefault("orb.cloud.config.agent_name", "")
	v.SetDefault("orb.cloud.config.auto_provision", true)
	v.SetDefault("orb.cloud.mqtt.address", "tls://agents.orb.live:8883")
	v.SetDefault("orb.cloud.mqtt.id", "")
	v.SetDefault("orb.cloud.mqtt.key", "")
	v.SetDefault("orb.cloud.mqtt.channel_id", "")
	v.SetDefault("orb.db.file", "./orb-agent.db")
	v.SetDefault("orb.tls.verify", true)
	v.SetDefault("orb.otel.host", "localhost")
	v.SetDefault("orb.otel.port", 0)
	v.SetDefault("orb.debug.enable", Debug)

	if len(path) > 0 {
		cobra.CheckErr(v.ReadInConfig())
	}

	var fZero float64

	// check that version of config files are all matched up
	if versionNumber1 := viper.GetFloat64("version"); versionNumber1 != fZero {
		versionNumber2 := v.GetFloat64("version")
		if versionNumber2 == fZero {
			cobra.CheckErr("Failed to parse config version in: " + path)
		}
		if versionNumber2 != versionNumber1 {
			cobra.CheckErr("Config file version mismatch in: " + path)
		}
	}

	// load backend static functions for setting up default values
	backendVarsFunction := make(map[string]func(*viper.Viper))
	backendVarsFunction["pktvisor"] = pktvisor.RegisterBackendSpecificVariables
	backendVarsFunction["otel"] = otel.RegisterBackendSpecificVariables

	// check if backends are configured
	// if not then add pktvisor as default
	if len(path) > 0 && len(v.GetStringMap("orb.backends")) == 0 {
		pktvisor.RegisterBackendSpecificVariables(v)
	} else {
		for backendName := range v.GetStringMap("orb.backends") {
			if backend := v.GetStringMap("orb.backends." + backendName); backend != nil {
				backendVarsFunction[backendName](v)
			}
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
	_ = rootCmd.Execute()
}
