// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package redis

import (
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type SinksUpdateEvent struct {
	SinkID    string
	Owner     string
	Config    types.Metadata
	Timestamp time.Time
}

type SinkerUpdateEvent struct {
	SinkID    string
	Owner     string
	State     string
	Msg       string
	Timestamp time.Time
}

func (cse SinkerUpdateEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"sink_id":   cse.SinkID,
		"owner":     cse.Owner,
		"state":     cse.State,
		"msg":       cse.Msg,
		"timestamp": cse.Timestamp.Unix(),
	}
}

type DeploymentEvent struct {
	SinkID         string
	DeploymentYaml string
}
