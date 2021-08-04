/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type AgentGroup struct {
	ID             string
	MFOwnerID      string
	Name           types.Identifier
	Description    string
	MFChannelID    string
	MatchingAgents types.Metadata
	Tags           types.Tags
	Created        time.Time
}

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")
	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")
	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")
	// ErrUnauthorizedAccess indicates while checking the credentials
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")
	// ErrScanMetadata indicates problem with metadata in db
	ErrScanMetadata = errors.New("failed to scan metadata in db")
	// ErrSelectEntity indicates error while reading entity from database
	ErrSelectEntity = errors.New("select entity from db error")
	// ErrEntityConnected indicates error while checking connection in database
	ErrEntityConnected = errors.New("check connection in database error")
	// ErrUpdateEntity indicates error while updating a entity
	ErrUpdateEntity = errors.New("failed to update entity")
	// ErrRemoveEntity indicates a error while deleting a agent group
	ErrRemoveEntity = errors.New("failed to remove entity")
)

type AgentGroupService interface {
	// CreateAgentGroup creates new AgentGroup, associated channel, applies to Agents as appropriate
	CreateAgentGroup(ctx context.Context, token string, s AgentGroup) (AgentGroup, error)
	// RetrieveAgentGroupByIDInternal Retrieve an AgentGroup by id, without a token
	RetrieveAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (AgentGroup, error)
}

type AgentGroupRepository interface {
	// Save persists the AgentGroup. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, group AgentGroup) (string, error)
	// RetrieveAllByAgent get all AgentGroup which an Agent belongs to.
	RetrieveAllByAgent(ctx context.Context, a Agent) ([]AgentGroup, error)
	// RetrieveByID get an AgentGroup by id
	RetrieveByID(ctx context.Context, groupID string, ownerID string) (AgentGroup, error)
}
