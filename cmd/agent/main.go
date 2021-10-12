/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/ns1labs/orb/agent"
	"github.com/ns1labs/orb/agent/backend/pktvisor"
	config2 "github.com/ns1labs/orb/agent/config"
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

	// handle signals
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
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

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run orb-agent and connect to Orb control plane",
		Long:  `Run orb-agent and connect to Orb control plane`,
		Run:   Run,
	}

	runCmd.Flags().StringSliceVarP(&cfgFiles, "config", "c", []string{}, "Path to config files (may be specified multiple times)")
	runCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Enable verbose (debug level) output")

	rootCmd.AddCommand(runCmd)
	rootCmd.Execute()
}
