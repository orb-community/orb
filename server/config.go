package server

import (
	"github.com/ns1labs/orb/datasource"
)

type ErrorLogLevel int

const (
	Error ErrorLogLevel = iota
	Warn
	Info
)

type DataStore struct {
	StoreType string `mapstructure:"type"`
}

type Services struct {
	DataStore datasource.ConsulConfig `mapstructure:"datasource"`
}

type AutoTLS struct {
	Enabled      bool   `mapstructure:"enabled"`
	Domain       string `mapstructure:"domain"`
	CertCacheDir string `mapstructure:"cache_dir"`
}

type TLSConfig struct {
	AutoTLS AutoTLS `mapstructure:"auto"`
	Cert    string  `mapstructure:"cert"`
	Key     string  `mapstructure:"key"`
}

type Listener struct {
	BindAddr  string    `mapstructure:"bind_addr"`
	TLSConfig TLSConfig `mapstructure:"tls"`
}

type ErrorLog struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

type AccessLog struct {
	Format string `mapstructure:"format"`
	Path   string `mapstructure:"path"`
}

type Loggers struct {
	ErrorLog  ErrorLog  `mapstructure:"error_log"`
	AccessLog AccessLog `mapstructure:"access_log"`
}

type Config struct {
	Listener Listener `mapstructure:"listener"`
	Services Services `mapstructure:"services"`
	Loggers  Loggers  `mapstructure:"loggers"`
}
