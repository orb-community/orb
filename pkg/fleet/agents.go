/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"time"
)

const (
	new     = "new"
	online  = "online"
	offline = "offline"
	stale   = "stale"
	removed = "removed"
)

type Agent struct {
	MFThingID     string
	MFOwnerID     string
	Created       time.Time
	OrbTags       Tags
	AgentTags     Tags
	AgentMetadata Metadata
	State         string
	LastHBData    Metadata
	LastHB        time.Time
}

type AgentService interface {
	// CreateAgent creates new agent
	CreateAgent(ctx context.Context, token string, a Agent) (Agent, error)
}

type AgentRepository interface {
	// Save persists the Agent. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, agent Agent) error
}
