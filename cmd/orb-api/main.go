package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"

	"github.com/ns1labs/orb/server"
)

func main() {
	rootCmd.Execute()
}

// default config path
var cfgFile = "/etc/orb/api.yml"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "orb",
	Short: "Orb Network observability platform api server",
	Long:  "Orb Network observability platform api sever allows centralized management of orb-agents",
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	var config server.Config
	viper.Unmarshal(&config)
	fmt.Printf("%+v\n", config)
	s := server.New(config)
	cobra.CheckErr(s.Serve())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		cobra.CheckErr(err)
	}
}
