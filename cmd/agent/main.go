/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
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

	rootCmd = &cobra.Command{
		Use:   "orb-agent",
		Short: "orb-agent connects to orb control plane",
		Long:  "orb-agent connects to orb control plane",
		Run:   Run,
	}
)

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

	// configuration
	var config config2.Config
	viper.Unmarshal(&config)

	config.Debug = Debug

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

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	pktvisor.Register()

	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringSliceVarP(&cfgFiles, "config", "c", []string{}, "Path to config files (may be specified multiple times)")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "verbose debug output")
}

func mergeOrError(path string) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	// note: viper seems to require a default (or a BindEnv) to be overridden by environment variables
	v.SetDefault("orb.cloud.api.address", "api.orb.live")
	v.SetDefault("orb.cloud.api.token", "")
	v.SetDefault("orb.cloud.config.agent_name", "")
	v.SetDefault("orb.cloud.config.auto_provision", true)
	v.SetDefault("orb.cloud.mqtt.address", "mqtt.orb.live")
	v.SetDefault("orb.cloud.mqtt.id", "")
	v.SetDefault("orb.cloud.mqtt.key", "")
	v.SetDefault("orb.cloud.mqtt.channel_id", "")
	v.SetDefault("orb.db.file", "./orb-agent.db")
	v.SetDefault("orb.tls.verify", true)

	cobra.CheckErr(v.ReadInConfig())
	var config config2.Config
	v.Unmarshal(&config)

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
	fmt.Fprintln(os.Stderr, "Using config file:", v.ConfigFileUsed())

	cobra.CheckErr(viper.MergeConfigMap(v.AllSettings()))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if len(cfgFiles) == 0 {
		mergeOrError(defaultConfig)
	} else {
		for _, conf := range cfgFiles {
			mergeOrError(conf)
		}
	}
}

func main() {
	rootCmd.Execute()
}
