package prometheus

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/backend"
	"gopkg.in/yaml.v3"
	"net/url"
)

var invalidCustomHeaders = []string{
	"Content-Encoding", "Content-Type", "X-Prometheus-Remote-Write-Version", "User-Agent",
}

func (p *Backend) ConfigToFormat(format string, metadata types.Metadata) (string, error) {
	if format == "yaml" {
		remoteHost := metadata[RemoteHostURLConfigFeature].(string)
		parseUtil := Backend{
			RemoteHost: remoteHost,
		}
		p.RemoteHost = remoteHost
		config, err := yaml.Marshal(parseUtil)
		if err != nil {
			return "", err
		}
		return string(config), nil
	} else {
		return "", errors.New("unsupported format")
	}
}

func (p *Backend) ParseConfig(format string, config string) (configReturn types.Metadata, err error) {
	if format == "yaml" {
		configAsByte := []byte(config)
		// Parse the YAML data into a Config struct
		configReturn = make(types.Metadata)
		err = yaml.Unmarshal(configAsByte, &configReturn)
		if err != nil {
			return nil, errors.Wrap(errors.New("failed to parse config YAML"), err)
		}
		return
	} else {
		return nil, errors.New("unsupported format")
	}
}

func (p *Backend) ValidateConfiguration(config types.Metadata) error {
	remoteUrl, remoteHostOk := config[RemoteHostURLConfigFeature]
	if !remoteHostOk {
		return errors.New("must send valid URL for Remote Write")
	}
	// Validate remote_host
	_, err := url.ParseRequestURI(remoteUrl.(string))
	if err != nil {
		return errors.New("must send valid URL for Remote Write")
	}
	// check for custom http headers
	customHeaders, customHeadersOk := config[CustomHeadersConfigFeature]
	if customHeadersOk {
		headersAsMap := customHeaders.(map[string]interface{})
		for _, header := range invalidCustomHeaders {
			if _, ok := headersAsMap[header]; !ok {
				return errors.New("invalid custom headers")
			}
		}
	}
	return nil
}

func (p *Backend) CreateFeatureConfig() []backend.ConfigFeature {
	var configs []backend.ConfigFeature

	remoteHost := backend.ConfigFeature{
		Type:     backend.ConfigFeatureTypeText,
		Input:    "text",
		Title:    "Remote Write URL",
		Name:     RemoteHostURLConfigFeature,
		Required: true,
	}

	configs = append(configs, remoteHost)
	return configs
}
