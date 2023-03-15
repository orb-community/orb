package otlpexporter

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/backend"
	"gopkg.in/yaml.v3"
)

// OTLP Exporter Examples
// exporters:
//  otlp:
//    endpoint: myserver.local:55690
//    tls:
//      insecure: false
//      ca_file: server.crt
//      cert_file: client.crt
//      key_file: client.key
//      min_version: "1.1"
//      max_version: "1.2"
//  otlp/insecure:
//    endpoint: myserver.local:55690
//    tls:
//      insecure: true
//  otlp/secure_no_verify:
//    endpoint: myserver.local:55690
//    tls:
//      insecure: false
//      insecure_skip_verify: true

type parserHelper struct {
	Endpoint string     `yaml:"endpoint"`
	Tls      *tlsConfig `yaml:"tls,omitempty,flow"`
	Auth     authConfig `yaml:"auth,flow"`
}

type tlsConfig struct {
	Insecure           *bool   `yaml:"insecure,omitempty"`
	CaFile             *string `yaml:"ca_file,omitempty"`
	CertFile           *string `yaml:"cert_file,omitempty"`
	KeyFile            *string `yaml:"key_file,omitempty"`
	MinVersion         *string `yaml:"min_version,omitempty"`
	MaxVersion         *string `yaml:"max_version,omitempty"`
	InsecureSkipVerify *bool   `yaml:"insecure_skip_verify,omitempty"`
}
type authConfig struct {
	Username *string `yaml:"username,omitempty"`
	Password *string `yaml:"password,omitempty"`
}

// CreateFeatureConfig Not available since this is only supported in YAML configuration
func (b Backend) CreateFeatureConfig() []backend.ConfigFeature {
	return nil
}

func (b Backend) ValidateConfiguration(config types.Metadata) error {
	if _, ok := config["endpoint"]; !ok {
		return errors.New("endpoint is required")
	}
	if _, ok := config["auth"]; !ok {
		return errors.New("auth config is required")
	}
	return nil
}

func (b Backend) ParseConfig(format string, config string) (retConfig types.Metadata, err error) {
	if format == "yaml" {
		var parsedConfig parserHelper
		err = yaml.Unmarshal([]byte(config), &parsedConfig)
		if err != nil {
			return nil, errors.Wrap(errors.New("failed to unmarshal config"), err)
		}
		retConfig = make(types.Metadata)
		retConfig["endpoint"] = parsedConfig.Endpoint
		authMap := make(types.Metadata)
		tlsMap := make(types.Metadata)
		if parsedConfig.Auth.Username != nil {
			authMap["username"] = *parsedConfig.Auth.Username
		}
		if parsedConfig.Auth.Password != nil {
			authMap["password"] = *parsedConfig.Auth.Password
		}
		retConfig["auth"] = authMap
		if parsedConfig.Tls != nil {
			if parsedConfig.Tls.Insecure != nil {
				tlsMap["insecure"] = *parsedConfig.Tls.Insecure
			}
			if parsedConfig.Tls.InsecureSkipVerify != nil {
				tlsMap["insecure_skip_verify"] = *parsedConfig.Tls.InsecureSkipVerify
			}
			if parsedConfig.Tls.CaFile != nil {
				tlsMap["ca_file"] = *parsedConfig.Tls.CaFile
			}
			if parsedConfig.Tls.CertFile != nil {
				tlsMap["cert_file"] = *parsedConfig.Tls.CertFile
			}
			if parsedConfig.Tls.KeyFile != nil {
				tlsMap["key_file"] = *parsedConfig.Tls.KeyFile
			}
			if parsedConfig.Tls.MinVersion != nil {
				tlsMap["min_version"] = *parsedConfig.Tls.MinVersion
			}
			if parsedConfig.Tls.MaxVersion != nil {
				tlsMap["max_version"] = *parsedConfig.Tls.MaxVersion
			}
			retConfig["tls"] = tlsMap
		}
	} else {
		return nil, errors.New("format not supported")
	}
	return
}

func (b Backend) ConfigToFormat(format string, metadata types.Metadata) (string, error) {
	if format == "yaml" {
		value, err := yaml.Marshal(metadata)
		return string(value), err
	} else {
		return "", errors.New("format not supported")
	}
}
