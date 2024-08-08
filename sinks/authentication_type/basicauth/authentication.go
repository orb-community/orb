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

const AuthType = "basicauth"

type AuthConfig struct {
	Username          string `json:"username" ,yaml:"username"`
	Password          string `json:"password" ,yaml:"password"`
	encryptionService authentication_type.PasswordService
}

func (a *AuthConfig) Metadata() authentication_type.AuthenticationTypeConfig {

	return authentication_type.AuthenticationTypeConfig{
		Type:        AuthType,
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
			if _, ok := value.(string); !ok {
				if key == "password" {
					return errors.Wrap(errors.ErrInvalidPasswordType, errors.New("invalid auth type for field: "+key))
				}
				if key == "type" {
					return errors.Wrap(errors.ErrInvalidAuthType, errors.New("invalid auth type for field: "+key))
				}
				if key == "username" {
					return errors.Wrap(errors.ErrInvalidUsernameType, errors.New("invalid auth type for field: "+key))
				}
			}
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
			retVal, err := yaml.Marshal(input)
			return string(retVal), err
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
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)
		authMeta[PasswordConfigFeature] = ""
		inputMeta[authentication_type.AuthenticationKey] = authMeta
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
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)
		authMeta[PasswordConfigFeature] = ""
		inputMeta[authentication_type.AuthenticationKey] = authMeta
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
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)
		if _, ok := authMeta[PasswordConfigFeature].(string); !ok {
			return nil, errors.Wrap(errors.ErrPasswordNotFound, errors.New("password field was not found"))
		}
		encoded, err := a.encryptionService.EncodePassword(authMeta[PasswordConfigFeature].(string))
		if err != nil {
			return nil, err
		}
		authMeta[PasswordConfigFeature] = encoded
		inputMeta[authentication_type.AuthenticationKey] = authMeta
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
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)
		encoded, err := a.encryptionService.EncodePassword(authMeta[PasswordConfigFeature].(string))
		if err != nil {
			return nil, err
		}
		authMeta[PasswordConfigFeature] = encoded
		inputMeta[authentication_type.AuthenticationKey] = authMeta
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
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)
		decoded, err := a.encryptionService.DecodePassword(authMeta[PasswordConfigFeature].(string))
		if err != nil {
			return nil, err
		}
		authMeta[PasswordConfigFeature] = decoded
		inputMeta[authentication_type.AuthenticationKey] = authMeta
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
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)
		decoded, err := a.encryptionService.DecodePassword(authMeta[PasswordConfigFeature].(string))
		if err != nil {
			return nil, err
		}
		authMeta[PasswordConfigFeature] = decoded
		inputMeta[authentication_type.AuthenticationKey] = authMeta
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
	authentication_type.Register(AuthType, &basicAuth)
}
