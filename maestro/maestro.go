package maestro

import (
	"context"
	"time"

	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
)

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrScanMetadata indicates problem with metadata in db
	ErrScanMetadata = errors.New("failed to scan metadata in db")

	// ErrSelectEntity indicates error while reading entity from database
	ErrSelectEntity = errors.New("select entity from db error")

	// ErrEntityConnected indicates error while checking connection in database
	ErrEntityConnected = errors.New("check connection in database error")

	// ErrUpdateEntity indicates error while updating a entity
	ErrUpdateEntity = errors.New("failed to update entity")

	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	ErrRemoveEntity = errors.New("failed to remove entity")

	ErrInvalidBackend = errors.New("No available backends")
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

type MaestroService interface {
	// CreateOtelCollector - create a existing collector by id
	CreateOtelCollector(ctx context.Context, sinkID string, msg string, ownerID string) error

	// DeleteOtelCollector - delete a existing collector by id
	DeleteOtelCollector(ctx context.Context, sinkID string, msg string, ownerID string) error

	// UpdateOtelCollector - update a existing collector by id
	UpdateOtelCollector(ctx context.Context, sinkID string, msg string, ownerID string) error
}
