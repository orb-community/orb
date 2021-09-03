/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

const CurrentRPCSchemaVersion = "1.0"

type RPC struct {
	SchemaVersion string      `json:"schema_version"`
	Func          string      `json:"func"`
	Payload       interface{} `json:"payload"`
}

const DatasetReqRPCFunc = "dataset_policy"

type DatasetRPC struct {
	SchemaVersion string              `json:"schema_version"`
	Func          string              `json:"func"`
	Payload       []DatasetRPCPayload `json:"payload"`
}

type DatasetRPCPayload struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	AgentGroupID  string      `json:"agent_group_id"`
	AgentPolicyID string      `json:"agent_policy_id"`
	SinkID        string      `json:"sink_id"`
	Data          interface{} `json:"data"`
}
