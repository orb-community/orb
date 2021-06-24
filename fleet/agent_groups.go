/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type AgentGroup struct {
	Name      types.Identifier
	MFOwnerID string
	Metadata  types.Metadata
	Created   time.Time
}

type AgentGroupService interface {
	// CreateAgentGroup creates new AgentGroup
	CreateAgentGroup(ctx context.Context, token string, s AgentGroup) (AgentGroup, error)
}

type AgentGroupRepository interface {
	// Save persists the AgentGroup. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, group AgentGroup) error
}
