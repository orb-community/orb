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
	SinkPrefix = "sinks."
	SinkCreate = SinkPrefix + "create"
	SinkDelete = SinkPrefix + "delete"
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
	id string
}

func (dse deleteSinkEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"id":        dse.id,
		"operation": SinkDelete,
	}
}

type updateSinkEvent struct {
	sinkID     string
	owner      string
	username   string
	password   string
	remoteHost string
	timestamp  time.Time
}

func (cce updateSinkEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"sink_id":    cce.sinkID,
		"owner":      cce.owner,
		"username":   cce.username,
		"password":   cce.password,
		"remoteHost": cce.timestamp.Unix(),
		"operation":  SinkUpdate,
	}
}
