package basicauth

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/backend"
	"gopkg.in/yaml.v3"
)

const (
	UsernameConfigFeature = "username"
	PasswordConfigFeature = "password"
)

var (
	features = []authentication_type.ConfigFeature{
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
	}
)

type AuthConfig struct {
	Username          string `json:"username" ,yaml:"username"`
	Password          string `json:"password" ,yaml:"password"`
	encryptionService authentication_type.PasswordService
}

func (a *AuthConfig) Metadata() interface{} {

	return authentication_type.AuthenticationTypeConfig{
		Type:        "basicauth",
		Description: "Basic username and password authentication",
		Config:      features,
	}
}

func (a *AuthConfig) GetFeatureConfig() []authentication_type.ConfigFeature {
	return features
}

func (a *AuthConfig) ValidateConfiguration(inputFormat string, input interface{}) error {
	switch inputFormat {
	case "object":
		for key, value := range input.(map[string]string) {
			if key == UsernameConfigFeature {
				if len(value) == 0 {
					return errors.New("username cannot be empty")
				}
			}
			if key == PasswordConfigFeature {
				if len(value) == 0 {
					return errors.New("password cannot be empty")
				}
			}
		}
	case "yaml":
		err := yaml.Unmarshal([]byte(input.(string)), &a)
		if err != nil {
			return err
		}
		if len(a.Username) == 0 {
			return errors.New("username cannot be empty")
		} else if len(a.Password) == 0 {
			return errors.New("password cannot be empty")
		}
	}
	return nil
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
