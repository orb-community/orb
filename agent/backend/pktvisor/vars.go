package pktvisor

import (
	"github.com/spf13/viper"
)

func RegisterBackendSpecificVariables(v *viper.Viper) {
	v.SetDefault("orb.backends.pktvisor.binary", "/usr/local/sbin/pktvisord")
	v.SetDefault("orb.backends.pktvisor.config_file", "/opt/orb/agent.yaml")
	v.SetDefault("orb.backends.pktvisor.api_host", "localhost")
	v.SetDefault("orb.backends.pktvisor.api_port", "10853")
}
