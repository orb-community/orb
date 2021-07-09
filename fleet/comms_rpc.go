/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

const CurrentRPCSchemaVersion = "1.0"

type RPC struct {
	SchemaVersion string      `json:"schema_version"`
	Func          string      `json:"func"`
	Payload       interface{} `json:"payload"`
}

// Core -> Edge

const GroupMembershipRPCFunc = "group_membership"

type GroupMembershipRPC struct {
	SchemaVersion string                    `json:"schema_version"`
	Func          string                    `json:"func"`
	Payload       GroupMembershipRPCPayload `json:"payload"`
}

type GroupMembershipRPCPayload struct {
	ChannelIDS []string `json:"channel_ids"`
	FullList   bool     `json:"full_list"`
}

const AgentPolicyRPCFunc = "agent_policy"

type AgentPolicyRPC struct {
	SchemaVersion string                  `json:"schema_version"`
	Func          string                  `json:"func"`
	Payload       []AgentPolicyRPCPayload `json:"payload"`
}

type AgentPolicyRPCPayload struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Backend string      `json:"backend"`
	Version int32       `json:"version"`
	Data    interface{} `json:"data"`
}

// Edge -> Core

const GroupMembershipReqRPCFunc = "group_membership_req"

type GroupMembershipReqRPCPayload struct {
	// empty
}
