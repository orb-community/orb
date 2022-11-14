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

const TagsRPCFunc = ""

const GroupMembershipRPCFunc = "group_membership"

type GroupMembershipRPC struct {
	SchemaVersion string                    `json:"schema_version"`
	Func          string                    `json:"func"`
	Payload       GroupMembershipRPCPayload `json:"payload"`
}

type GroupMembershipData struct {
	GroupID   string `json:"group_id"`
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
	FullList      bool                    `json:"full_list"`
}

type AgentPolicyRPCPayload struct {
	Action       string      `json:"action"`
	ID           string      `json:"id"`
	DatasetID    string      `json:"dataset_id"`
	AgentGroupID string      `json:"agent_group_id"`
	Name         string      `json:"name"`
	Backend      string      `json:"backend"`
	Version      int32       `json:"version"`
	Data         interface{} `json:"data"`
}

const GroupRemovedRPCFunc = "group_removed"

type GroupRemovedRPC struct {
	SchemaVersion string                 `json:"schema_version"`
	Func          string                 `json:"func"`
	Payload       GroupRemovedRPCPayload `json:"payload"`
}

type GroupRemovedRPCPayload struct {
	AgentGroupID string   `json:"agent_group_id"`
	ChannelID    string   `json:"channel_id"`
	Datasets     []string `json:"datasets"`
}

const DatasetRemovedRPCFunc = "dataset_removed"

type DatasetRemovedRPC struct {
	SchemaVersion string                   `json:"schema_version"`
	Func          string                   `json:"func"`
	Payload       DatasetRemovedRPCPayload `json:"payload"`
}

type DatasetRemovedRPCPayload struct {
	DatasetID string `json:"dataset_id"`
	PolicyID  string `json:"policy_id"`
}

const AgentStopRPCFunc = "agent_stop"

type AgentStopRPCPayload struct {
	Reason string `json:"reason"`
}

type AgentStopRPC struct {
	SchemaVersion string              `json:"schema_version"`
	Func          string              `json:"func"`
	Payload       AgentStopRPCPayload `json:"payload"`
}

const AgentResetRPCFunc = "agent_reset"

type AgentResetRPCPayload struct {
	FullReset bool   `json:"full_reset"`
	Reason    string `json:"reason"`
}

type AgentResetRPC struct {
	SchemaVersion string               `json:"schema_version"`
	Func          string               `json:"func"`
	Payload       AgentResetRPCPayload `json:"payload"`
}

// Edge -> Core

const GroupMembershipReqRPCFunc = "group_membership_req"

type GroupMembershipReqRPCPayload struct {
	// empty
}

const AgentPoliciesReqRPCFunc = "agent_policies_req"

type AgentPoliciesReqRPCPayload struct {
	// empty
}

const AgentOrbConfigReqRPCFunc = "agent_tags_req"

type AgentOrbConfigReqRPCPayload struct {
	// empty
}

type AgentTagsRPCPayload struct {
	AgentName string            `json:"agent_name"`
	OrbTags   map[string]string `json:"orb_tags"`
}

const AgentMetricsRPCFunc = "agent_metrics"

type AgentMetricsRPC struct {
	SchemaVersion string                   `json:"schema_version"`
	Func          string                   `json:"func"`
	Payload       []AgentMetricsRPCPayload `json:"payload"`
}

type AgentMetricsRPCPayload struct {
	PolicyID   string   `json:"policy_id"`
	PolicyName string   `json:"policy_name"`
	Datasets   []string `json:"datasets"`
	Format     string   `json:"format"`
	BEVersion  string   `json:"be_version"`
	Data       []byte   `json:"data"`
}
