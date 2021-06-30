/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"errors"
	"time"
)

var (
	// ErrSchemaVersion a message was received indicating a version we don't support
	ErrSchemaVersion = errors.New("unsupported schema version")
	// ErrSchemaMalformed a message contained a schema we couldn't parse
	ErrSchemaMalformed = errors.New("schema malformed")
	// ErrPayloadTooBig a message contained a payload that was abnormally large
	ErrPayloadTooBig = errors.New("payload too big")
)

// MaxMsgPayloadSize maximum payload size we will process from a client
const MaxMsgPayloadSize = 1024 * 5

// MQTT messaging schemes

type SchemaVersionCheck struct {
	SchemaVersion string `json:"schema_version"`
}

type OrbAgentInfo struct {
	Version string `json:"version"`
}

type BackendInfo struct {
	Version string `json:"version"`
}

const CurrentCapabilitiesSchemaVersion = "1.0"

type Capabilities struct {
	SchemaVersion string                 `json:"schema_version"`
	OrbAgent      OrbAgentInfo           `json:"orb_agent"`
	AgentTags     map[string]string      `json:"agent_tags"`
	Backends      map[string]BackendInfo `json:"backends"`
}

const CurrentHeartbeatSchemaVersion = "1.0"

type Heartbeat struct {
	SchemaVersion string    `json:"schema_version"`
	TimeStamp     time.Time `json:"ts"`
	State         State     `json:"state"`
}

const CurrentRPCSchemaVersion = "1.0"

type RPC struct {
	SchemaVersion string      `json:"schema_version"`
	Func          string      `json:"func"`
	Payload       interface{} `json:"payload"`
}

const GroupMembershipRPCFunc = "group_membership"

type GroupMembershipRPCPayload struct {
	ChannelIDS []string `json:"channel_ids"`
}
