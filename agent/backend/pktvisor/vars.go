package pktvisor

import (
	"github.com/spf13/viper"
)

func RegisterBackendSpecificVariables() {
	viper.SetDefault("orb.backends.pktvisor.binary", "/usr/local/sbin/pktvisord")
	viper.SetDefault("orb.backends.pktvisor.config_file", "/opt/orb/agent.yaml")
	viper.SetDefault("orb.backends.pktvisor.api_host", "localhost")
	viper.SetDefault("orb.backends.pktvisor.api_port", "10853")
}
