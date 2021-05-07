// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package producer

import (
	"time"
)

const (
	AgentPrefix = "agent."
	AgentCreate = AgentPrefix + "create"
)

type event interface {
	encode() map[string]interface{}
}

var (
	_ event = (*createAgentEvent)(nil)
)

type createAgentEvent struct {
	mfThing   string
	owner     string
	name      string
	content   string
	timestamp time.Time
}

func (cce createAgentEvent) encode() map[string]interface{} {
	return map[string]interface{}{
		"thing_id":  cce.mfThing,
		"owner":     cce.owner,
		"name":      cce.name,
		"content":   cce.content,
		"timestamp": cce.timestamp.Unix(),
		"operation": AgentCreate,
	}
}
