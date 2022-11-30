package main

import (
	"github.com/spf13/cobra"
	"testing"
)

func Test_e2e_orbAgent_ConfigFile(t *testing.T) {
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
	err := rootCmd.Execute()
	if err != nil {
		t.Fail()
	}

}
