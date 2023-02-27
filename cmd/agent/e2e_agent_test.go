package main

import (
	"context"
	"fmt"
	"github.com/orb-community/orb/agent"
	"github.com/orb-community/orb/agent/backend/pktvisor"
	"github.com/orb-community/orb/agent/config"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_e2e_orbAgent_ConfigFile(t *testing.T) {
	t.Skip("local run only, skip in CICD")
	defer profile.Start().Stop()
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
	rootCmd.SetArgs([]string{"run", "-d", "-c", "/home/lpegoraro/workspace/orb/localconfig/config.yaml"})
	ctx, cancelF := context.WithTimeout(context.Background(), 2*time.Minute)
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		t.Fail()
	}

	select {
	case <-ctx.Done():
		cancelF()
		return
	}
}

func Test_main(t *testing.T) {
	t.Skip("local run only, skip in CICD")

	mergeOrError("/home/lpegoraro/workspace/orb/localconfig/config.yaml")

	// configuration
	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("agent start up error (config): %w", err))
		os.Exit(1)
	}

	cfg.OrbAgent.Debug.Enable = true

	// include pktvisor backend by default if binary is at default location
	_, err = os.Stat(pktvisor.DefaultBinary)
	if err == nil && cfg.OrbAgent.Backends == nil {
		cfg.OrbAgent.Backends = make(map[string]map[string]string)
		cfg.OrbAgent.Backends["pktvisor"] = make(map[string]string)
		cfg.OrbAgent.Backends["pktvisor"]["binary"] = pktvisor.DefaultBinary
		if len(cfgFiles) > 0 {
			cfg.OrbAgent.Backends["pktvisor"]["config_file"] = "/home/lpegoraro/workspace/orb/localconfig/config.yaml"
		}
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

	// new agent
	a, err := agent.New(logger, cfg)
	if err != nil {
		logger.Error("agent start up error", zap.Error(err))
		os.Exit(1)
	}

	// handle signals
	done := make(chan bool, 1)
	rootCtx, cancelFunc := context.WithTimeout(context.WithValue(context.Background(), "routine", "mainRoutine"), 15*time.Minute)

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
