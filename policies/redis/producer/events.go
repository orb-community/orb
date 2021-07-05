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
	DatasetPrefix = "dataset."
	DatasetCreate = DatasetPrefix + "create"
	PolicyPrefix  = "policy."
	PolicyCreate  = PolicyPrefix + "create"
)

type event interface {
	Encode() map[string]interface{}
}

var (
	_ event = (*createDatasetEvent)(nil)
	_ event = (*createPolicyEvent)(nil)
)

type createDatasetEvent struct {
	id           string
	owner        string
	name         string
	agentGroupID string
	policyID     string
	sinkID       string
	timestamp    time.Time
}

type createPolicyEvent struct {
	id        string
	owner     string
	name      string
	timestamp time.Time
}

func (cce createDatasetEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"id":        cce.id,
		"owner":     cce.owner,
		"name":      cce.name,
		"timestamp": cce.timestamp.Unix(),
		"operation": DatasetCreate,
	}
}

func (cce createPolicyEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"id":        cce.id,
		"owner":     cce.owner,
		"name":      cce.name,
		"timestamp": cce.timestamp.Unix(),
		"operation": PolicyCreate,
	}
}
