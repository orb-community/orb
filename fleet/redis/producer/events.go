// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"time"
)

const (
	AgentPrefix      = "agent."
	AgentCreate      = AgentPrefix + "create"
	AgentGroupPrefix = "agent_group."
	AgentGroupRemove = AgentGroupPrefix + "remove"
)

type event interface {
	encode() map[string]interface{}
}

var (
	_ event = (*createAgentEvent)(nil)
	_ event = (*removeAgentGroupEvent)(nil)
)

type createAgentEvent struct {
	mfThing   string
	owner     string
	name      string
	content   string
	timestamp time.Time
}

type removeAgentGroupEvent struct {
	groupID   string
	token     string
	timestamp time.Time
}

func (rde removeAgentGroupEvent) encode() map[string]interface{} {
	return map[string]interface{}{
		"group_id":  rde.groupID,
		"token":     rde.token,
		"timestamp": rde.timestamp.Unix(),
		"operation": AgentGroupRemove,
	}
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
