package config

import (
	"github.com/orb-community/orb/maestro/password"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type/basicauth"
	"github.com/orb-community/orb/sinks/authentication_type/bearertokenauth"
)

const AuthenticationKey = "authentication"

type AuthBuilderService interface {
	GetExtensionsFromMetadata(config types.Metadata) (Extensions, string)
	DecodeAuth(config types.Metadata) (types.Metadata, error)
	EncodeAuth(config types.Metadata) (types.Metadata, error)
}

func GetAuthService(authType string, service password.EncryptionService) AuthBuilderService {
	switch authType {
	case basicauth.AuthType:
		return &BasicAuthBuilder{
			encryptionService: service,
		}
	case bearertokenauth.AuthType:
		return &BearerTokenAuthBuilder{
			encryptionService: service,
		}
	}

	return nil
}

type BasicAuthBuilder struct {
	encryptionService password.EncryptionService
}

func (b *BasicAuthBuilder) GetExtensionsFromMetadata(c types.Metadata) (Extensions, string) {
	authcfg := c.GetSubMetadata(AuthenticationKey)
	username := authcfg["username"].(string)
	password := authcfg["password"].(string)
	return Extensions{
		BasicAuth: &BasicAuthenticationExtension{
			ClientAuth: &ClientAuth{
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
	config[AuthenticationKey] = authCfg
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
	config[AuthenticationKey] = authcfg
	return config, nil
}

type BearerTokenAuthBuilder struct {
	encryptionService password.EncryptionService
}

func (b *BearerTokenAuthBuilder) GetExtensionsFromMetadata(c types.Metadata) (Extensions, string) {
	authcfg := c.GetSubMetadata(AuthenticationKey)
	scheme := authcfg["scheme"].(string)
	token := authcfg["token"].(string)

	return Extensions{
		BearerAuth: &BearerTokenAuthExtension{
			Scheme: scheme,
			Token:  token,
		},
	}, "bearertokenauth/withscheme"
}

func (b *BearerTokenAuthBuilder) DecodeAuth(config types.Metadata) (types.Metadata, error) {
	authCfg := config.GetSubMetadata(AuthenticationKey)
	token := authCfg["token"].(string)

	decodedToken, err := b.encryptionService.DecodePassword(token)
	if err != nil {
		return nil, err
	}

	authCfg["token"] = decodedToken
	config[AuthenticationKey] = authCfg

	return config, nil
}

func (b *BearerTokenAuthBuilder) EncodeAuth(config types.Metadata) (types.Metadata, error) {
	authcfg := config.GetSubMetadata(AuthenticationKey)
	token := authcfg["token"].(string)

	encodedToken, err := b.encryptionService.EncodePassword(token)
	if err != nil {
		return nil, err
	}

	authcfg["token"] = encodedToken
	config[AuthenticationKey] = authcfg

	return config, nil
}
