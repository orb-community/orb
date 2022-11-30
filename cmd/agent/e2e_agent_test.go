package main

import (
	"context"
	"github.com/spf13/cobra"
	"testing"
	"time"
)

func Test_e2e_orbAgent_ConfigFile(t *testing.T) {
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
