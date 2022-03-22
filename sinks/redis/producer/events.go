// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"encoding/json"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

const (
	SinkPrefix = "sinks."
	SinkCreate = SinkPrefix + "create"
	SinkDelete = SinkPrefix + "remove"
	SinkUpdate = SinkPrefix + "update"
)

type event interface {
	Encode() map[string]interface{}
}

var (
	_ event = (*createSinkEvent)(nil)
)

type createSinkEvent struct {
	mfThing   string
	owner     string
	name      string
	content   string
	timestamp time.Time
}

func (cce createSinkEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"thing_id":  cce.mfThing,
		"owner":     cce.owner,
		"name":      cce.name,
		"content":   cce.content,
		"timestamp": cce.timestamp.Unix(),
		"operation": SinkCreate,
	}
}

type deleteSinkEvent struct {
	sinkID  string
	ownerID string
}

func (dse deleteSinkEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"sink_id":   dse.sinkID,
		"owner_id":  dse.ownerID,
		"operation": SinkDelete,
	}
}

type updateSinkEvent struct {
	sinkID    string
	owner     string
	config    types.Metadata
	state     string
	timestamp time.Time
}

func (cce updateSinkEvent) Encode() (map[string]interface{}, error) {
	config, err := json.Marshal(cce.config)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sink_id":   cce.sinkID,
		"owner":     cce.owner,
		"config":    config,
		"state":     cce.state,
		"timestamp": cce.timestamp.Unix(),
		"operation": SinkUpdate,
	}, nil

}
