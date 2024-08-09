package noauth

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type"
	"gopkg.in/yaml.v3"
)

var features = []authentication_type.ConfigFeature{}

const AuthType = "noauth"

type AuthConfig struct{}

func (a *AuthConfig) GetFeatureConfig() []authentication_type.ConfigFeature {
	return features
}

func (a *AuthConfig) ValidateConfiguration(_ string, _ interface{}) error {
	return nil
}

func (a *AuthConfig) ConfigToFormat(outputFormat string, input interface{}) (interface{}, error) {
	switch input.(type) {
	case types.Metadata:
		if outputFormat == "yaml" {
			retVal, err := yaml.Marshal(input)
			return string(retVal), err
		}
	case string:
		if outputFormat == "object" {
			retVal := make(types.Metadata)
			val := input.(string)
			err := yaml.Unmarshal([]byte(val), &retVal)
			return retVal, err
		} else {
			return nil, errors.New("unsupported format")
		}
	}
	return nil, errors.New("unsupported format")
}

// OmitInformation just pass-through the config without changes
func (a *AuthConfig) OmitInformation(_ string, c interface{}) (interface{}, error) {
	return c, nil
}

// EncodeInformation just pass-through the config without changes
func (a *AuthConfig) EncodeInformation(_ string, c interface{}) (interface{}, error) {

	return c, nil
}

// DecodeInformation just pass-through the config without changes
func (a *AuthConfig) DecodeInformation(_ string, c interface{}) (interface{}, error) {
	return c, nil
}

func (a *AuthConfig) Metadata() authentication_type.AuthenticationTypeConfig {
	return authentication_type.AuthenticationTypeConfig{
		Type:        AuthType,
		Description: "No authentication",
		Config:      features,
	}
}

func Register(_ authentication_type.PasswordService) {
	authentication_type.Register(AuthType, &AuthConfig{})
}
