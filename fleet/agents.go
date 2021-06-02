/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"database/sql/driver"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

const (
	New State = iota
	Online
	Offline
	Stale
	Removed
)

type State int

var stateMap = [...]string{
	"new",
	"online",
	"offline",
	"stale",
	"removed",
}

var stateRevMap = map[string]State{
	"new":     New,
	"online":  Online,
	"offline": Offline,
	"stale":   Stale,
	"removed": Removed,
}

func (s State) String() string {
	return stateMap[s]
}

func (s *State) Scan(value interface{}) error { *s = stateRevMap[value.(string)]; return nil }
func (s State) Value() (driver.Value, error)  { return s.String(), nil }

type Agent struct {
	Name          types.Identifier
	MFOwnerID     string
	MFThingID     string
	MFKeyID       string
	MFChannelID   string
	Created       time.Time
	OrbTags       Tags
	AgentTags     Tags
	AgentMetadata Metadata
	State         State
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

type AgentComms interface {
	StartComms() error
}
