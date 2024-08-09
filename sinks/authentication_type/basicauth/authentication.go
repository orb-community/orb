package basicauth

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/backend"
)

const (
	AuthType              = "basicauth"
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
	Username          *string `json:"username" yaml:"username"`
	Password          *string `json:"password" yaml:"password"`
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
		if _, ok := input.(types.Metadata)[UsernameConfigFeature]; !ok {
			return errors.Wrap(errors.ErrAuthUsernameNotFound, errors.New("username field was not found"))
		}

		if _, ok := input.(types.Metadata)[PasswordConfigFeature]; !ok {
			return errors.Wrap(errors.ErrAuthPasswordNotFound, errors.New("password field was not found"))
		}

		for key, value := range input.(types.Metadata) {
			if key == UsernameConfigFeature {
				if _, ok := value.(string); !ok {
					return errors.Wrap(errors.ErrAuthInvalidUsernameType, errors.New("invalid auth type for field: "+key))
				}

				if len(strings.Fields(value.(string))) == 0 {
					return errors.Wrap(errors.ErrAuthInvalidUsernameType, errors.New("invalid authentication username"))
				}
			}

			if key == PasswordConfigFeature {
				if _, ok := value.(string); !ok {
					return errors.Wrap(errors.ErrAuthInvalidPasswordType, errors.New("invalid auth type for field: "+key))
				}

				if len(strings.Fields(value.(string))) == 0 {
					return errors.Wrap(errors.ErrAuthInvalidPasswordType, errors.New("invalid authentication password"))
				}
			}
		}
	case "yaml":
		err := yaml.Unmarshal([]byte(input.(string)), &a)
		if err != nil {
			return err
		}

		if a.Username == nil {
			return errors.Wrap(errors.ErrAuthUsernameNotFound, errors.New("username field was not found"))
		}

		if len(strings.Fields(*a.Username)) == 0 {
			return errors.Wrap(errors.ErrAuthInvalidUsernameType, errors.New("invalid authentication username"))
		}

		if a.Password == nil {
			return errors.Wrap(errors.ErrAuthPasswordNotFound, errors.New("password field was not found"))
		}

		if len(strings.Fields(*a.Password)) == 0 {
			return errors.Wrap(errors.ErrAuthInvalidPasswordType, errors.New("invalid authentication password"))
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
			return nil, errors.Wrap(errors.ErrAuthPasswordNotFound, errors.New("password field was not found"))
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
