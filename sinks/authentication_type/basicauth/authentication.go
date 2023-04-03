package basicauth

import (
	"github.com/orb-community/orb/sinks/authentication_type"
)

type AuthConfig struct {
	username string `json:"username", yaml:"username"`
	password string `json:"password", yaml:"password"`
}

func (a *AuthConfig) Metadata() interface{} {
	//TODO implement me
	panic("implement me")
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
