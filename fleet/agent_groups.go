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
	Tags           types.Tags
	Created        time.Time
	MatchingAgents types.Metadata
}

type PageAgentGroup struct {
	PageMetadata
	AgentGroups []AgentGroup
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
	// ViewAgentGroupByID Retrieve an AgentGroup by id
	ViewAgentGroupByID(ctx context.Context, token string, id string) (AgentGroup, error)
	// ViewAgentGroupByIDInternal Retrieve an AgentGroup by id, without a token
	ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (AgentGroup, error)
	// ListAgentGroups Retrieve a list of AgentGroups by owner
	ListAgentGroups(ctx context.Context, token string, pm PageMetadata) (PageAgentGroup, error)
	// EditAgentGroup edit a existing agent group by id and owner
	EditAgentGroup(ctx context.Context, token string, ag AgentGroup) (AgentGroup, error)
	// RemoveAgentGroup Remove a existing agent group by owner an id
	RemoveAgentGroup(ctx context.Context, token string, id string) error
	// ValidateAgentGroup validate AgentGroup
	ValidateAgentGroup(ctx context.Context, token string, s AgentGroup) (AgentGroup, error)
}

type AgentGroupRepository interface {
	// Save persists the AgentGroup. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, group AgentGroup) (string, error)
	// RetrieveAllByAgent get all AgentGroup which an Agent belongs to.
	RetrieveAllByAgent(ctx context.Context, a Agent) ([]AgentGroup, error)
	// RetrieveByID get an AgentGroup by id
	RetrieveByID(ctx context.Context, groupID string, ownerID string) (AgentGroup, error)
	// RetrieveAllAgentGroupsByOwner get all AgentGroup by owner.
	RetrieveAllAgentGroupsByOwner(ctx context.Context, ownerID string, pm PageMetadata) (PageAgentGroup, error)
	// Update a existing agent group by owner and id
	Update(ctx context.Context, ownerID string, group AgentGroup) (AgentGroup, error)
	// Delete a existing agent group by owner and id
	Delete(ctx context.Context, groupID string, ownerID string) error
	// RetrieveMatchingGroups Groups this Agent currently belongs to, according to matching agent and group tags
	RetrieveMatchingGroups(ctx context.Context, ownerID string, thingID string) (MatchingGroups, error)
}
