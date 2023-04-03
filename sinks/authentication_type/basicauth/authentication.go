package basicauth

import (
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/backend"
)

const (
	UsernameConfigFeature = "username"
	PasswordConfigFeature = "password"
)

type AuthConfig struct {
	username string `json:"username", yaml:"username"`
	password string `json:"password", yaml:"password"`
}

func (a *AuthConfig) Metadata() interface{} {
	return authentication_type.AuthenticationTypeConfig{
		Type:        "basicauth",
		Description: "Basic username and password authentication",
		Config: []authentication_type.ConfigFeature{
			{
				Type:     backend.ConfigFeatureTypeText,
				Input:    "text",
				Title:    "Username",
				Name:     UsernameConfigFeature,
				Required: true,
			},
			{
				Type:     backend.ConfigFeatureTypePassword,
				Input:    "text",
				Title:    "Password",
				Name:     PasswordConfigFeature,
				Required: true,
			},
		},
	}
}

func (a *AuthConfig) GetFeatureConfig() []authentication_type.ConfigFeature {
	//TODO implement me
	panic("implement me")
}

func (a *AuthConfig) ValidateConfiguration(inputFormat string, input interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthConfig) ConfigToFormat(outputFormat string, input interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthConfig) OmitInformation(outputFormat string, input interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthConfig) EncodeInformation(outputFormat string, input interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthConfig) DecodeInformation(outputFormat string, input interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
