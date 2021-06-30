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

type GroupMembershipRPCPayload struct {
	ChannelIDS []string `json:"channel_ids"`
	FullList   bool     `json:"full_list"`
}

// Edge -> Core

const GroupMembershipReqRPCFunc = "group_membership_req"

type GroupMembershipReqRPCPayload struct {
	// empty
}
