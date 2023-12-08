package noauth

import "github.com/orb-community/orb/sinks/authentication_type"

var features = []authentication_type.ConfigFeature{}

const AuthType = "noauth"

type AuthConfig struct{}

func (a *AuthConfig) GetFeatureConfig() []authentication_type.ConfigFeature {
	return features
}

func (a *AuthConfig) ValidateConfiguration(_ string, _ interface{}) error {
	return nil
}

func (a *AuthConfig) ConfigToFormat(_ string, _ interface{}) (interface{}, error) {
	return nil, nil
}

func (a *AuthConfig) OmitInformation(_ string, _ interface{}) (interface{}, error) {
	return nil, nil
}

func (a *AuthConfig) EncodeInformation(_ string, _ interface{}) (interface{}, error) {
	return nil, nil
}

func (a *AuthConfig) DecodeInformation(_ string, _ interface{}) (interface{}, error) {
	return nil, nil
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
