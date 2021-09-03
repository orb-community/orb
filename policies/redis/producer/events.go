// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

const (
	DatasetPrefix = "dataset."
	DatasetCreate = DatasetPrefix + "create"
	DatasetRemove = DatasetPrefix + "remove"
	PolicyPrefix  = "policy."
	PolicyCreate  = PolicyPrefix + "create"
	PolicyUpdate  = PolicyPrefix + "update"
)

type event interface {
	Encode() map[string]interface{}
}

var (
	_ event = (*createDatasetEvent)(nil)
	_ event = (*createPolicyEvent)(nil)
	_ event = (*updatePolicyEvent)(nil)
)

type createDatasetEvent struct {
	id           string
	ownerID      string
	name         string
	agentGroupID string
	policyID     string
	sinkID       string
	timestamp    time.Time
}

type createPolicyEvent struct {
	id        string
	ownerID   string
	name      string
	timestamp time.Time
}

type updatePolicyEvent struct {
	id       string
	ownerID  string
	groupIDs []string
	policy   types.Metadata
}

func (cce createDatasetEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"id":        cce.id,
		"owner_id":  cce.ownerID,
		"name":      cce.name,
		"group_id":  cce.agentGroupID,
		"policy_id": cce.policyID,
		"sink_id":   cce.sinkID,
		"timestamp": cce.timestamp.Unix(),
		"operation": DatasetCreate,
	}
}

func (cce createPolicyEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"id":        cce.id,
		"owner_id":  cce.ownerID,
		"name":      cce.name,
		"timestamp": cce.timestamp.Unix(),
		"operation": PolicyCreate,
	}
}

func (cce updatePolicyEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"id":        cce.id,
		"owner_id":  cce.ownerID,
		"groups_id": cce.groupIDs,
		"policy":    cce.policy,
		"operation": PolicyUpdate,
	}
}
