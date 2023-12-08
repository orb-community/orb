package config

import (
	maestrobasicauth "github.com/orb-community/orb/maestro/config/basicauth"
	maestronoauth "github.com/orb-community/orb/maestro/config/noauth"
	"github.com/orb-community/orb/maestro/password"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type/basicauth"
	"github.com/orb-community/orb/sinks/authentication_type/noauth"
)

type AuthBuilderService interface {
	GetExtensionsFromMetadata(config types.Metadata) (Extensions, string)
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
