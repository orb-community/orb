package noauth

import (
	"github.com/orb-community/orb/maestro/config"
	"github.com/orb-community/orb/pkg/types"
)

type NoAuthBuilder struct {
}

func (b *NoAuthBuilder) GetExtensionsFromMetadata(c types.Metadata) (config.Extensions, string) {
	return config.Extensions{}, ""
}

func (b *NoAuthBuilder) DecodeAuth(c types.Metadata) (types.Metadata, error) {
	return c, nil
}

func (b *NoAuthBuilder) EncodeAuth(c types.Metadata) (types.Metadata, error) {
	return c, nil
}
