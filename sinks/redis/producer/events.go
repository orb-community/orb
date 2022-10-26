// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"encoding/json"
	"time"

	"github.com/ns1labs/orb/pkg/types"
)

const (
	SinkPrefix = "sinks."
	SinkCreate = SinkPrefix + "create"
	SinkDelete = SinkPrefix + "remove"
	SinkUpdate = SinkPrefix + "update"
)

type event interface {
	Encode() (map[string]interface{}, error)
}

var (
	_ event = (*createSinkEvent)(nil)
)

type createSinkEvent struct {
	sinkID    string
	owner     string
	config    types.Metadata
	timestamp time.Time
}

func (cce createSinkEvent) Encode() (map[string]interface{}, error) {
	config, err := json.Marshal(cce.config)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sink_id":   cce.sinkID,
		"owner":     cce.owner,
		"config":    config,
		"timestamp": cce.timestamp.Unix(),
		"operation": SinkCreate,
	}, nil
}

type deleteSinkEvent struct {
	sinkID    string
	owner     string
	config    types.Metadata
	timestamp time.Time
}

func (dse deleteSinkEvent) Encode() (map[string]interface{}, error) {
	config, err := json.Marshal(dse.config)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sink_id":   dse.sinkID,
		"owner":     dse.owner,
		"config":    config,
		"timestamp": dse.timestamp.Unix(),
		"operation": SinkDelete,
	}, nil
}

type updateSinkEvent struct {
	sinkID    string
	owner     string
	config    types.Metadata
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
		"timestamp": cce.timestamp.Unix(),
		"operation": SinkUpdate,
	}, nil

}
