package datasource

type TlsConfig struct {
	Cert   string `mapstructure:"cert"`
	Key    string `mapstructure:"key"`
	RootCA string `mapstructure:"rootCA"`
	Verify bool   `mapstructure:"verify"`
}

type ConsulConfig struct {
	Address      string    `mapstructure:"address"`
	Access_token string    `mapstructure:"access_token"`
	Tls          TlsConfig `mapstructure:"tls,omitempty"`
}
