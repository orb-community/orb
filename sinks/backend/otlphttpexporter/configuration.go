package otlphttpexporter

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

const EndpointFieldName = "endpoint"

type OTLPHTTPBackend struct {
	Endpoint string `yaml:"endpoint"`
	//TODO will keep TLS until we confirm there is no need for those
	//Tls      *tlsConfig `yaml:"tls,omitempty,flow"`
}

func (b *OTLPHTTPBackend) Metadata() interface{} {
	return backend.SinkFeature{
		Backend:     "otlphttp",
		Description: "OTLP Exporter over HTTP",
		Config:      b.CreateFeatureConfig(),
	}
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

func Register() bool {
	backend.Register("otlphttp", &OTLPHTTPBackend{})
	return true
}

// CreateFeatureConfig Not available since this is only supported in YAML configuration
func (b *OTLPHTTPBackend) CreateFeatureConfig() []backend.ConfigFeature {
	var configs []backend.ConfigFeature

	remoteHost := backend.ConfigFeature{
		Type:     backend.ConfigFeatureTypeText,
		Input:    "text",
		Title:    "Remote Write URL",
		Name:     EndpointFieldName,
		Required: true,
	}

	configs = append(configs, remoteHost)
	return configs
}

func (b *OTLPHTTPBackend) ValidateConfiguration(config types.Metadata) error {
	if config[EndpointFieldName] == "" {
		return errors.New("malformed entity specification. endpoint must not be empty")
	}
	if _, ok := config[EndpointFieldName]; !ok {
		return errors.New("malformed entity specification. endpoint field is expected on exporter field")
	}
	return nil
}

func (b *OTLPHTTPBackend) ParseConfig(format string, config string) (retConfig types.Metadata, err error) {
	if format == "yaml" {
		var parsedConfig OTLPHTTPBackend
		err = yaml.Unmarshal([]byte(config), &parsedConfig)
		if err != nil {
			return nil, errors.Wrap(errors.New("failed to unmarshal config"), err)
		}
		retConfig = make(types.Metadata)
		retConfig[EndpointFieldName] = parsedConfig.Endpoint

	} else {
		return nil, errors.New("format not supported")
	}
	return
}

func (b *OTLPHTTPBackend) ConfigToFormat(format string, metadata types.Metadata) (string, error) {
	if format == "yaml" {
		value, err := yaml.Marshal(metadata)
		return string(value), err
	} else {
		return "", errors.New("format not supported")
	}
}
