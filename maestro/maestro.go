package maestro

import (
	"context"
	"time"

	"github.com/ns1labs/orb/pkg/types"
)

type Maestro struct {
	ID          string
	Name        types.Identifier
	MFOwnerID   string
	Description string
	Backend     string
	Config      types.Metadata
	Error       string
	Created     time.Time
}

type SinkConfig struct {
	Id       string
	Url      string
	Username string
	Password string
}

type Service interface {
	// Start starts the service - load the configuration
	Start(ctx context.Context, cancelFunction context.CancelFunc) error
}
