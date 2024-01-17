package basicauth

import (
	"github.com/orb-community/orb/maestro/config/output"
	"github.com/orb-community/orb/maestro/password"
	"github.com/orb-community/orb/pkg/types"
)

const AuthenticationKey = "authentication"

type BasicAuthBuilder struct {
	EncryptionService password.EncryptionService
}

func (b *BasicAuthBuilder) GetExtensionsFromMetadata(c types.Metadata) (output.Extensions, string) {
	authcfg := c.GetSubMetadata(AuthenticationKey)
	username := authcfg["username"].(string)
	password := authcfg["password"].(string)
	return output.Extensions{
		BasicAuth: &output.BasicAuthenticationExtension{
			ClientAuth: &output.ClientAuth{
				Username: username,
				Password: password,
			},
		},
	}, "basicauth/exporter"
}

func (b *BasicAuthBuilder) DecodeAuth(c types.Metadata) (types.Metadata, error) {
	authCfg := c.GetSubMetadata(AuthenticationKey)
	password := authCfg["password"].(string)
	decodedPassword, err := b.EncryptionService.DecodePassword(password)
	if err != nil {
		return nil, err
	}
	authCfg["password"] = decodedPassword
	c[AuthenticationKey] = authCfg
	return c, nil
}

func (b *BasicAuthBuilder) EncodeAuth(c types.Metadata) (types.Metadata, error) {
	authcfg := c.GetSubMetadata(AuthenticationKey)
	password := authcfg["password"].(string)
	encodedPassword, err := b.EncryptionService.EncodePassword(password)
	if err != nil {
		return nil, err
	}
	authcfg["password"] = encodedPassword
	c[AuthenticationKey] = authcfg
	return c, nil
}
