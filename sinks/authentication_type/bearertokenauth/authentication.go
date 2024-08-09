package bearertokenauth

import (
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/backend"
)

const (
	AuthType            = "bearertokenauth"
	SchemeConfigFeature = "scheme"
	TokenConfigFeature  = "token"
)

var features = []authentication_type.ConfigFeature{
	{
		Type:     backend.ConfigFeatureTypeText,
		Input:    "text",
		Title:    "Scheme",
		Name:     SchemeConfigFeature,
		Required: true,
	},
	{
		Type:     backend.ConfigFeatureTypeText,
		Input:    "text",
		Title:    "Token",
		Name:     TokenConfigFeature,
		Required: true,
	},
}

type (
	AuthConfig struct {
		encryptionService authentication_type.PasswordService

		Scheme *string `json:"scheme" ,yaml:"scheme"`
		Token  *string `json:"token" ,yaml:"token"`
	}
)

func (a *AuthConfig) GetFeatureConfig() []authentication_type.ConfigFeature {
	return features
}

func (a *AuthConfig) ValidateConfiguration(inputFormat string, input any) error {
	switch inputFormat {
	case "object":
		if _, ok := input.(types.Metadata)[SchemeConfigFeature]; !ok {
			return errors.Wrap(errors.ErrAuthSchemeNotFound, errors.New("scheme field was not found"))
		}

		if _, ok := input.(types.Metadata)[TokenConfigFeature]; !ok {
			return errors.Wrap(errors.ErrAuthTokenNotFound, errors.New("token field was not found"))
		}

		for key, value := range input.(types.Metadata) {
			if key == SchemeConfigFeature {
				if _, ok := value.(string); !ok {
					return errors.Wrap(errors.ErrAuthInvalidSchemeType, errors.New("invalid auth type for field: "+key))
				}

				if len(strings.Fields(value.(string))) != 1 {
					return errors.Wrap(errors.ErrAuthInvalidSchemeType, errors.New("invalid authentication scheme"))
				}
			}

			if key == TokenConfigFeature {
				if _, ok := value.(string); !ok {
					return errors.Wrap(errors.ErrAuthInvalidTokenType, errors.New("invalid auth type for field: "+key))
				}
				if len(strings.Fields(value.(string))) != 1 {
					return errors.Wrap(errors.ErrAuthInvalidTokenType, errors.New("invalid authentication token"))
				}
			}
		}
	case "yaml":
		err := yaml.Unmarshal([]byte(input.(string)), &a)
		if err != nil {
			return err
		}

		if a.Scheme == nil {
			return errors.Wrap(errors.ErrAuthSchemeNotFound, errors.New("scheme field was not found"))
		}

		if len(strings.Fields(*a.Scheme)) != 1 {
			return errors.Wrap(errors.ErrAuthInvalidSchemeType, errors.New("invalid authentication scheme"))
		}

		if a.Token == nil {
			return errors.Wrap(errors.ErrAuthTokenNotFound, errors.New("token field was not found"))
		}

		if len(strings.Fields(*a.Token)) != 1 {
			return errors.Wrap(errors.ErrAuthInvalidTokenType, errors.New("invalid authentication token"))
		}
	}

	return nil
}

func (a *AuthConfig) ConfigToFormat(outputFormat string, input any) (any, error) {
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
		}
	}

	return nil, errors.New("unsupported format")
}

func (a *AuthConfig) OmitInformation(outputFormat string, input any) (any, error) {
	switch input.(type) {
	case types.Metadata:
		inputMeta := input.(types.Metadata)
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)
		authMeta[TokenConfigFeature] = ""
		inputMeta[authentication_type.AuthenticationKey] = authMeta

		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		}

		if outputFormat == "object" {
			return inputMeta, nil
		}

		return nil, errors.New("unsupported format")
	case string:
		iia, err := a.ConfigToFormat("object", input)
		if err != nil {
			return nil, err
		}
		inputMeta := iia.(types.Metadata)
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)
		authMeta[TokenConfigFeature] = ""
		inputMeta[authentication_type.AuthenticationKey] = authMeta

		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		}

		if outputFormat == "object" {
			return inputMeta, nil
		}

		return nil, errors.New("unsupported format")
	}

	return nil, errors.New("unsupported format")
}

func (a *AuthConfig) EncodeInformation(outputFormat string, input interface{}) (interface{}, error) {
	switch input.(type) {
	case types.Metadata:
		inputMeta := input.(types.Metadata)
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)

		if _, ok := authMeta[TokenConfigFeature].(string); !ok {
			return nil, errors.Wrap(errors.ErrAuthTokenNotFound, errors.New("token field was not found"))
		}

		encoded, err := a.encryptionService.EncodePassword(authMeta[TokenConfigFeature].(string))
		if err != nil {
			return nil, err
		}

		authMeta[TokenConfigFeature] = encoded
		inputMeta[authentication_type.AuthenticationKey] = authMeta

		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		}

		if outputFormat == "object" {
			return inputMeta, nil
		}

		return nil, errors.New("unsupported format")
	case string:
		iia, err := a.ConfigToFormat("object", input)
		if err != nil {
			return nil, err
		}
		inputMeta := iia.(types.Metadata)
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)

		encoded, err := a.encryptionService.EncodePassword(authMeta[TokenConfigFeature].(string))
		if err != nil {
			return nil, err
		}

		authMeta[TokenConfigFeature] = encoded
		inputMeta[authentication_type.AuthenticationKey] = authMeta

		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		}

		if outputFormat == "object" {
			return inputMeta, nil
		}

		return nil, errors.New("unsupported format")
	}

	return nil, errors.New("unsupported format")
}

func (a *AuthConfig) DecodeInformation(outputFormat string, input any) (any, error) {
	switch input.(type) {
	case types.Metadata:
		inputMeta := input.(types.Metadata)
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)

		decoded, err := a.encryptionService.DecodePassword(authMeta[TokenConfigFeature].(string))
		if err != nil {
			return nil, err
		}

		authMeta[TokenConfigFeature] = decoded
		inputMeta[authentication_type.AuthenticationKey] = authMeta

		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		}

		if outputFormat == "object" {
			return inputMeta, nil
		}

		return nil, errors.New("unsupported format")
	case string:
		iia, err := a.ConfigToFormat("object", input)
		if err != nil {
			return nil, err
		}

		inputMeta := iia.(types.Metadata)
		authMeta := inputMeta.GetSubMetadata(authentication_type.AuthenticationKey)

		decoded, err := a.encryptionService.DecodePassword(authMeta[TokenConfigFeature].(string))
		if err != nil {
			return nil, err
		}

		authMeta[TokenConfigFeature] = decoded
		inputMeta[authentication_type.AuthenticationKey] = authMeta

		if outputFormat == "yaml" {
			return a.ConfigToFormat("yaml", inputMeta)
		}

		if outputFormat == "object" {
			return inputMeta, nil
		}

		return nil, errors.New("unsupported format")
	}

	return nil, errors.New("unsupported format")
}

func (a *AuthConfig) Metadata() authentication_type.AuthenticationTypeConfig {
	return authentication_type.AuthenticationTypeConfig{
		Type:        AuthType,
		Description: "Token authentication",
		Config:      features,
	}
}

func Register(encryptionService authentication_type.PasswordService) {
	bearerTokenAuth := AuthConfig{
		encryptionService: encryptionService,
	}
	authentication_type.Register(AuthType, &bearerTokenAuth)
}
