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

type GroupMembershipData struct {
	Name      string `json:"name"`
	ChannelID string `json:"channel_id"`
}

type GroupMembershipRPCPayload struct {
	Groups   []GroupMembershipData `json:"groups"`
	FullList bool                  `json:"full_list"`
}

const AgentPolicyRPCFunc = "agent_policy"

type AgentPolicyRPC struct {
	SchemaVersion string                  `json:"schema_version"`
	Func          string                  `json:"func"`
	Payload       []AgentPolicyRPCPayload `json:"payload"`
}

type AgentPolicyRPCPayload struct {
	Action    string      `json:"action"`
	ID        string      `json:"id"`
	DatasetID string      `json:"dataset_id"`
	Name      string      `json:"name"`
	Backend   string      `json:"backend"`
	Version   int32       `json:"version"`
	Data      interface{} `json:"data"`
}

const GroupRemovedRPCFunc = "group_removed"

// Edge -> Core

const GroupMembershipReqRPCFunc = "group_membership_req"

type GroupMembershipReqRPCPayload struct {
	// empty
}

const AgentPoliciesReqRPCFunc = "agent_policies_req"

type AgentPoliciesReqRPCPayload struct {
	// empty
}

const AgentMetricsRPCFunc = "agent_metrics"

type AgentMetricsRPC struct {
	SchemaVersion string                   `json:"schema_version"`
	Func          string                   `json:"func"`
	Payload       []AgentMetricsRPCPayload `json:"payload"`
}

type AgentMetricsRPCPayload struct {
	PolicyID  string   `json:"policy_id"`
	Datasets  []string `json:"datasets"`
	Format    string   `json:"format"`
	BEVersion string   `json:"be_version"`
	Data      []byte   `json:"data"`
}
