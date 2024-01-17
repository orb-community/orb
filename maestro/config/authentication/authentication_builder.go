package authentication

import (
	maestrobasicauth "github.com/orb-community/orb/maestro/config/authentication/basicauth"
	maestronoauth "github.com/orb-community/orb/maestro/config/authentication/noauth"
	"github.com/orb-community/orb/maestro/config/output"
	"github.com/orb-community/orb/maestro/password"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type/basicauth"
	"github.com/orb-community/orb/sinks/authentication_type/noauth"
)

type AuthBuilderService interface {
	GetExtensionsFromMetadata(config types.Metadata) (output.Extensions, string)
	DecodeAuth(config types.Metadata) (types.Metadata, error)
	EncodeAuth(config types.Metadata) (types.Metadata, error)
}

func GetAuthService(authType string, service password.EncryptionService) AuthBuilderService {
	switch authType {
	case basicauth.AuthType:
		return &maestrobasicauth.BasicAuthBuilder{
			EncryptionService: service,
		}
	case noauth.AuthType:
		return &maestronoauth.NoAuthBuilder{}
	}
	return nil
}
