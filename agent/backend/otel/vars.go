package otel

import "github.com/spf13/viper"

func RegisterBackendSpecificVariables(v *viper.Viper) {
	v.SetDefault("orb.backends.otel.otlp_port", "4316")
}
