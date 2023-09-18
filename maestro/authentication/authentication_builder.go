package authentication

import (
	"github.com/orb-community/orb/maestro/config"
	"github.com/orb-community/orb/maestro/deployment"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type/basicauth"
)

const AuthenticationKey = "authentication"

type AuthBuilderService interface {
	GetExtensionsFromMetadata(config types.Metadata) (config.Extensions, string)
	DecodeAuth(config types.Metadata) (types.Metadata, error)
	EncodeAuth(config types.Metadata) (types.Metadata, error)
}

func GetAuthService(authType string, service deployment.EncryptionService) AuthBuilderService {
	switch authType {
	case basicauth.AuthType:
		return &BasicAuthBuilder{
			encryptionService: service,
		}
	}
	return nil
}

type BasicAuthBuilder struct {
	encryptionService deployment.EncryptionService
}

func (b *BasicAuthBuilder) GetExtensionsFromMetadata(c types.Metadata) (config.Extensions, string) {
	authcfg := c.GetSubMetadata(AuthenticationKey)
	username := authcfg["username"].(string)
	password := authcfg["password"].(string)
	return config.Extensions{
		BasicAuth: &config.BasicAuthenticationExtension{
			ClientAuth: &config.ClientAuth{
				Username: username,
				Password: password,
			},
		},
	}, "basicauth/exporter"
}

func (b *BasicAuthBuilder) DecodeAuth(config types.Metadata) (types.Metadata, error) {
	authCfg := config.GetSubMetadata(AuthenticationKey)
	password := authCfg["password"].(string)
	decodedPassword, err := b.encryptionService.DecodePassword(password)
	if err != nil {
		return nil, err
	}
	authCfg["password"] = decodedPassword
	return config, nil
}

func (b *BasicAuthBuilder) EncodeAuth(config types.Metadata) (types.Metadata, error) {
	authcfg := config.GetSubMetadata(AuthenticationKey)
	password := authcfg["password"].(string)
	encodedPassword, err := b.encryptionService.EncodePassword(password)
	if err != nil {
		return nil, err
	}
	authcfg["password"] = encodedPassword
	return config, nil
}
