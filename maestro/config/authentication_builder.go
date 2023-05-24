package config

import (
	"github.com/orb-community/orb/pkg/types"
)

const AuthenticationKey = "authentication"

type AuthBuilderService interface {
	GetExtensionsFromMetadata(config types.Metadata) (Extensions, string)
}

func GetAuthService(authType string) AuthBuilderService {
	switch authType {
	case "basicauth":
		return &BasicAuthBuilder{}
	}
	return nil
}

type BasicAuthBuilder struct {
}

func (b *BasicAuthBuilder) GetExtensionsFromMetadata(config types.Metadata) (Extensions, string) {

	authcfg := config.GetSubMetadata(AuthenticationKey)
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
