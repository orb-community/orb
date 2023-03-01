/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prometheus

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/backend"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
	"io"
	"net/url"
)

var _ backend.Backend = (*prometheusBackend)(nil)

const (
	RemoteHostURLConfigFeature = "remote_host"
	UsernameConfigFeature      = "username"
	PasswordConfigFeature      = "password"
	ApiTokenConfigFeature      = "api_token"
)

//type PrometheusConfigMetadata = types.Metadata

type AuthType int

const (
	BasicAuth AuthType = iota
	TokenAuth
)

type prometheusBackend struct {
	apiHost     string
	apiPort     uint64
	apiUser     string
	apiPassword string
}

type configParseUtility struct {
	RemoteHost string  `yaml:"remote_host"`
	Username   *string `yaml:"username,omitempty"`
	Password   *string `yaml:"password,omitempty"`
	APIToken   *string `yaml:"api_token,omitempty"`
}

type SinkFeature struct {
	Backend     string                  `json:"backend"`
	Description string                  `json:"description"`
	Config      []backend.ConfigFeature `json:"config"`
}

func (p *prometheusBackend) Metadata() interface{} {
	return SinkFeature{
		Backend:     "prometheus",
		Description: "Prometheus time series database sink",
		Config:      p.CreateFeatureConfig(),
	}
}

func (p *prometheusBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string) error {
	return nil
}

func Register() bool {
	backend.Register("prometheus", &prometheusBackend{})
	return true
}

func (p *prometheusBackend) ConfigToFormat(format string, metadata types.Metadata) (string, error) {
	if format == "yaml" {
		username := metadata[UsernameConfigFeature].(*string)
		password := metadata[PasswordConfigFeature].(string)
		parseUtil := configParseUtility{
			RemoteHost: metadata[RemoteHostURLConfigFeature].(string),
			Username:   username,
			Password:   &password,
		}
		config, err := yaml.Marshal(parseUtil)
		if err != nil {
			return "", err
		}
		return string(config), nil
	} else {
		return "", errors.New("unsupported format")
	}
}

func (p *prometheusBackend) ParseConfig(format string, config string) (configReturn types.Metadata, err error) {
	if format == "yaml" {
		configAsByte := []byte(config)
		// Parse the YAML data into a Config struct
		var configUtil configParseUtility
		err = yaml.Unmarshal(configAsByte, &configUtil)
		if err != nil {
			return nil, errors.New("failed to parse config YAML")
		}
		configReturn = make(types.Metadata)
		// Check for Token Auth
		configReturn[RemoteHostURLConfigFeature] = configUtil.RemoteHost
		configReturn[UsernameConfigFeature] = configUtil.Username
		configReturn[PasswordConfigFeature] = configUtil.Password
		return
	} else {
		return nil, errors.New("unsupported format")
	}
}

func (p *prometheusBackend) ValidateConfiguration(config types.Metadata) error {
	authType := BasicAuth
	for _, key := range maps.Keys(config) {
		if key == ApiTokenConfigFeature {
			authType = TokenAuth
			break
		}
	}
	switch authType {
	case BasicAuth:
		_, userOk := config[UsernameConfigFeature]
		_, passwordOk := config[PasswordConfigFeature]
		if !userOk || !passwordOk {
			return errors.New("basic authentication, must provide username and password fields")
		}
	case TokenAuth:
		return errors.New("not implemented yet")
	}
	remoteUrl, remoteHostOk := config[RemoteHostURLConfigFeature]
	if !remoteHostOk {
		return errors.New("must send valid URL for Remote Write")
	}
	// Validate remote_host
	_, err := url.ParseRequestURI(remoteUrl.(string))
	if err != nil {
		return errors.New("must send valid URL for Remote Write")
	}
	return nil
}

func (p *prometheusBackend) CreateFeatureConfig() []backend.ConfigFeature {
	var configs []backend.ConfigFeature

	remoteHost := backend.ConfigFeature{
		Type:     backend.ConfigFeatureTypeText,
		Input:    "text",
		Title:    "Remote Write URL",
		Name:     RemoteHostURLConfigFeature,
		Required: true,
	}

	userName := backend.ConfigFeature{
		Type:     backend.ConfigFeatureTypeText,
		Input:    "text",
		Title:    "Username",
		Name:     UsernameConfigFeature,
		Required: true,
	}
	password := backend.ConfigFeature{
		Type:     backend.ConfigFeatureTypePassword,
		Input:    "text",
		Title:    "Password",
		Name:     PasswordConfigFeature,
		Required: true,
	}
	configs = append(configs, remoteHost, userName, password)
	return configs
}
