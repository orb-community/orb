package prometheus

import (
	"net/url"

	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/backend"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

func (p *Backend) ConfigToFormat(format string, metadata types.Metadata, reqType string) (string, error) {

	if reqType == "get" {
		username := metadata[UsernameConfigFeature].(string)
		password := metadata[PasswordConfigFeature].(string)
		parseUtil := configParseUtility{
			RemoteHost: metadata[RemoteHostURLConfigFeature].(string),
			Username:   &username,
			Password:   &password,
		}
		config, err := yaml.Marshal(parseUtil)
		if err != nil {
			return "", err
		}
		return string(config), nil
	} else {
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
	}

}

func (p *Backend) ParseConfig(format string, config string) (configReturn types.Metadata, err error) {

	configAsByte := []byte(config)
	// Parse the YAML data into a Config struct
	var configUtil configParseUtility
	err = yaml.Unmarshal(configAsByte, &configUtil)
	if err != nil {
		return nil, errors.Wrap(errors.New("failed to parse config YAML"), err)
	}
	configReturn = make(types.Metadata)
	// Check for Token Auth
	configReturn[RemoteHostURLConfigFeature] = configUtil.RemoteHost
	configReturn[UsernameConfigFeature] = configUtil.Username
	configReturn[PasswordConfigFeature] = configUtil.Password
	return

}

func (p *Backend) ValidateConfiguration(config types.Metadata) error {
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

func (p *Backend) CreateFeatureConfig() []backend.ConfigFeature {
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
