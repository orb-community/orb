/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"github.com/ns1labs/orb/agent"
	"github.com/ns1labs/orb/agent/backend/pktvisor"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/buildinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	defaultConfig = "/etc/orb/agent.yaml"
)

var (
	cfgFiles []string
	Debug    bool
)

func init() {

	pktvisor.Register()

}

func Version(_ *cobra.Command, _ []string) {
	fmt.Printf("orb-agent %s\n", buildinfo.GetVersion())
	os.Exit(0)
}

func Run(_ *cobra.Command, _ []string) {

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
	var configData config.Config
	err = viper.Unmarshal(&configData)
	if err != nil {
		logger.Error("agent start up error (config)", zap.Error(err))
		os.Exit(1)
	}

	configData.Debug = Debug
	// include pktvisor backend by default if binary is at default location
	_, err = os.Stat(pktvisor.DefaultBinary)
	if err == nil && configData.OrbAgent.Backends == nil {
		configData.OrbAgent.Backends = make(map[string]map[string]string)
		configData.OrbAgent.Backends["pktvisor"] = make(map[string]string)
		configData.OrbAgent.Backends["pktvisor"]["binary"] = pktvisor.DefaultBinary
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

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sigs
		a.Stop()
		done <- true
	}()

	// start agent
	err = a.Start()
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
	v.SetDefault("orb.otel.enable", false)
	v.SetDefault("orb.debug.enable", false)

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
			cobra.CheckErr("Failed to parse config version in: " + path)
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
	_ = rootCmd.Execute()
}
