package basicauth

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
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

func (a *AuthConfig) Metadata() authentication_type.AuthenticationTypeConfig {

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
		for key, value := range input.(types.Metadata) {
			vs := value.(string)
			if key == UsernameConfigFeature {
				if len(vs) == 0 {
					return errors.New("username cannot be empty")
				}
			}
			if key == PasswordConfigFeature {
				if len(vs) == 0 {
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
	switch input.(type) {
	case types.Metadata:
		if outputFormat == "yaml" {
			return yaml.Marshal(input)
		} else {
			return nil, errors.New("unsupported format")
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

func (a *AuthConfig) OmitInformation(outputFormat string, input interface{}) (interface{}, error) {
	switch input.(type) {
	case types.Metadata:
		inputMeta := input.(types.Metadata)
		inputMeta["password"] = ""
		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		} else if outputFormat == "object" {
			return inputMeta, nil
		} else {
			return nil, errors.New("unsupported format")
		}
	case string:
		iia, err := a.ConfigToFormat("object", input)
		if err != nil {
			return nil, err
		}
		inputMeta := iia.(types.Metadata)
		inputMeta["password"] = ""
		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		} else if outputFormat == "object" {
			return inputMeta, nil
		} else {
			return nil, errors.New("unsupported format")
		}
	}
	return nil, errors.New("unsupported format")
}

func (a *AuthConfig) EncodeInformation(outputFormat string, input interface{}) (interface{}, error) {
	switch input.(type) {
	case types.Metadata:
		inputMeta := input.(types.Metadata)
		encoded, err := a.encryptionService.EncodePassword(inputMeta["password"].(string))
		if err != nil {
			return nil, err
		}
		inputMeta["password"] = encoded
		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		} else if outputFormat == "object" {
			return inputMeta, nil
		} else {
			return nil, errors.New("unsupported format")
		}
	case string:
		iia, err := a.ConfigToFormat("object", input)
		if err != nil {
			return nil, err
		}
		inputMeta := iia.(types.Metadata)
		encoded, err := a.encryptionService.EncodePassword(inputMeta["password"].(string))
		if err != nil {
			return nil, err
		}
		inputMeta["password"] = encoded
		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		} else if outputFormat == "object" {
			return inputMeta, nil
		} else {
			return nil, errors.New("unsupported format")
		}
	}
	return nil, errors.New("unsupported format")
}

func (a *AuthConfig) DecodeInformation(outputFormat string, input interface{}) (interface{}, error) {
	switch input.(type) {
	case types.Metadata:
		inputMeta := input.(types.Metadata)
		decoded, err := a.encryptionService.DecodePassword(inputMeta["password"].(string))
		if err != nil {
			return nil, err
		}
		inputMeta["password"] = decoded
		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		} else if outputFormat == "object" {
			return inputMeta, nil
		} else {
			return nil, errors.New("unsupported format")
		}
	case string:
		iia, err := a.ConfigToFormat("object", input)
		if err != nil {
			return nil, err
		}
		inputMeta := iia.(types.Metadata)
		decoded, err := a.encryptionService.DecodePassword(inputMeta["password"].(string))
		if err != nil {
			return nil, err
		}
		inputMeta["password"] = decoded
		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		} else if outputFormat == "object" {
			return inputMeta, nil
		} else {
			return nil, errors.New("unsupported format")
		}
	}
	return nil, errors.New("unsupported format")
}

func Register(encryptionService authentication_type.PasswordService) {
	basicAuth := AuthConfig{
		encryptionService: encryptionService,
	}
	authentication_type.Register("basicauth", &basicAuth)
}
