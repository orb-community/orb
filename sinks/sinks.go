/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"context"
	"database/sql/driver"
	"github.com/etaques/orb/pkg/errors"
	"github.com/etaques/orb/pkg/types"
	"github.com/etaques/orb/sinks/backend"
	"time"
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

const (
	Unknown State = iota
	Active
	Error
	Idle
)

type State int

var stateMap = [...]string{
	"unknown",
	"active",
	"error",
	"idle",
}

const MetadataLabelOtel = "opentelemetry"

type Filter struct {
	StateFilter   string
	OpenTelemetry string
}

var stateRevMap = map[string]State{
	"unknown": Unknown,
	"active":  Active,
	"error":   Error,
	"idle":    Idle,
}

func (s State) String() string {
	return stateMap[s]
}

func (s *State) Scan(value interface{}) error {
	asString, ok := value.(string)
	if !ok {
		asBytes, ok := value.([]byte)
		if !ok {
			return errors.New("Scan source is not []byte")
		}
		asString = string(asBytes)
	}
	*s = stateRevMap[string(asString)]
	return nil
}
func (s State) Value() (driver.Value, error) { return s.String(), nil }

type Sink struct {
	ID          string
	Name        types.Identifier
	MFOwnerID   string
	Description string
	Backend     string
	Config      types.Metadata
	Tags        types.Tags
	State       State
	Error       string
	Created     time.Time
}

// Page contains page related metadata as well as list of sinks that
// belong to this page
type Page struct {
	PageMetadata
	Sinks []Sink
}

// SinkService Sink CRUD interface
type SinkService interface {
	// CreateSink creates new data sink
	CreateSink(ctx context.Context, token string, s Sink) (Sink, error)
	// UpdateSink by id
	UpdateSink(ctx context.Context, token string, s Sink) (Sink, error)
	// ListSinks retrieves data about sinks
	ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error)
	// ListSinksInternal retrieves data from sinks filtered by SinksFilter for Services like Maestro, to build DeploymentEntries
	ListSinksInternal(ctx context.Context, filter Filter) ([]Sink, error)
	// ListBackends retrieves a list of available backends
	ListBackends(ctx context.Context, token string) ([]string, error)
	// ViewBackend retrieves a backend by the name
	ViewBackend(ctx context.Context, token string, key string) (backend.Backend, error)
	// ViewSink retrieves a sink by id, for View, does not send password
	ViewSink(ctx context.Context, token string, key string) (Sink, error)
	// ViewSinkInternal retrieves a sink by id, via GRPC, sends password
	ViewSinkInternal(ctx context.Context, ownerID string, key string) (Sink, error)
	// DeleteSink delete a existing sink by id
	DeleteSink(ctx context.Context, token string, key string) error
	// ValidateSink validate a sink configuration without saving
	ValidateSink(ctx context.Context, token string, sink Sink) (Sink, error)
	// ChangeSinkStateInternal change the sink internal state from new/idle/active
	ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state State) error
}

type SinkRepository interface {
	// Save persists the Sink. Successful operation is indicated by non-nil error response.
	Save(ctx context.Context, sink Sink) (string, error)
	// Update performs an update to the existing sink, A non-nil error is
	// returned to indicate operation failure
	Update(ctx context.Context, sink Sink) error
	// RetrieveAllByOwnerID retrieves Sinks by OwnerID
	RetrieveAllByOwnerID(ctx context.Context, owner string, pm PageMetadata) (Page, error)
	// SearchAllSinks search Sinks for internal usage like services
	SearchAllSinks(ctx context.Context, filter Filter) ([]Sink, error)
	// RetrieveById retrieves a Sink by ID
	RetrieveById(ctx context.Context, key string) (Sink, error)
	// RetrieveByOwnerAndId retrieves a By OwnerId And SinkId
	RetrieveByOwnerAndId(ctx context.Context, ownerID string, key string) (Sink, error)
	// Remove an existing Sink by id
	Remove(ctx context.Context, owner string, key string) error
	// UpdateSinkState updates sink state like active, idle, new, unknown
	UpdateSinkState(ctx context.Context, sinkID string, msg string, ownerID string, state State) error
}
