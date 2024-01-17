package noauth

import (
	"github.com/orb-community/orb/maestro/config/output"
	"github.com/orb-community/orb/pkg/types"
)

type NoAuthBuilder struct {
}

func (b *NoAuthBuilder) GetExtensionsFromMetadata(c types.Metadata) (output.Extensions, string) {
	return output.Extensions{}, ""
}

func (b *NoAuthBuilder) DecodeAuth(c types.Metadata) (types.Metadata, error) {
	return c, nil
}

func (b *NoAuthBuilder) EncodeAuth(c types.Metadata) (types.Metadata, error) {
	return c, nil
}
