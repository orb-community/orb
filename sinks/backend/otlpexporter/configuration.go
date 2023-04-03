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
	Endpoint string `yaml:"endpoint"`
	//TODO will keep TLS until we confirm there is no need for those
	//Tls      *tlsConfig `yaml:"tls,omitempty,flow"`
}

// TODO will keep TLS until we confirm there is no need for those
type tlsConfig struct {
	Insecure           *bool   `yaml:"insecure,omitempty"`
	CaFile             *string `yaml:"ca_file,omitempty"`
	CertFile           *string `yaml:"cert_file,omitempty"`
	KeyFile            *string `yaml:"key_file,omitempty"`
	MinVersion         *string `yaml:"min_version,omitempty"`
	MaxVersion         *string `yaml:"max_version,omitempty"`
	InsecureSkipVerify *bool   `yaml:"insecure_skip_verify,omitempty"`
}

// CreateFeatureConfig Not available since this is only supported in YAML configuration
func (b parserHelper) CreateFeatureConfig() []backend.ConfigFeature {
	return nil
}

func (b parserHelper) ValidateConfiguration(config types.Metadata) error {
	if _, ok := config["endpoint"]; !ok {
		return errors.New("endpoint is required")
	}
	return nil
}

func (b parserHelper) ParseConfig(format string, config string) (retConfig types.Metadata, err error) {
	if format == "yaml" {
		var parsedConfig parserHelper
		err = yaml.Unmarshal([]byte(config), &parsedConfig)
		if err != nil {
			return nil, errors.Wrap(errors.New("failed to unmarshal config"), err)
		}
		retConfig = make(types.Metadata)
		retConfig["endpoint"] = parsedConfig.Endpoint

	} else {
		return nil, errors.New("format not supported")
	}
	return
}

func (b parserHelper) ConfigToFormat(format string, metadata types.Metadata) (string, error) {
	if format == "yaml" {
		value, err := yaml.Marshal(metadata)
		return string(value), err
	} else {
		return "", errors.New("format not supported")
	}
}
