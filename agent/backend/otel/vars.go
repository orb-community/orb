package otel

import "github.com/spf13/viper"

func RegisterBackendSpecificVariables() {
	viper.SetDefault("orb.backends.otel.otlp_port", "4316")
}
