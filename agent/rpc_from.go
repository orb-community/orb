/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"encoding/json"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
)

func (a *orbAgent) handleGroupMembership(rpc fleet.GroupMembershipRPCPayload) {
	// if this is the full list, reset all group subscriptions and subscribed to this list
	if rpc.FullList {
		a.unsubscribeGroupChannels()
		a.groupChannels = a.subscribeGroupChannels(rpc.Groups)
	} else {
		// otherwise, just add these subscriptions to the existing list
		successList := a.subscribeGroupChannels(rpc.Groups)
		a.groupChannels = append(a.groupChannels, successList...)
	}
}

func (a *orbAgent) handleAgentPolicies(rpc []fleet.AgentPolicyRPCPayload) {

	for _, payload := range rpc {
		a.policyManager.ManagePolicy(payload)
	}

}

func (a *orbAgent) handleGroupRPCFromCore(client mqtt.Client, message mqtt.Message) {

	a.logger.Debug("Group RPC message from core", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))

	var rpc fleet.RPC
	if err := json.Unmarshal(message.Payload(), &rpc); err != nil {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaMalformed))
		return
	}
	if rpc.SchemaVersion != fleet.CurrentRPCSchemaVersion {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaVersion))
		return
	}
	if rpc.Func == "" || rpc.Payload == nil {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaMalformed))
		return
	}

	// dispatch
	switch rpc.Func {
	case fleet.AgentPolicyRPCFunc:
		var r fleet.AgentPolicyRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding agent policy message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleAgentPolicies(r.Payload)
	default:
		a.logger.Warn("unsupported/unhandled core RPC, ignoring",
			zap.String("func", rpc.Func),
			zap.Any("payload", rpc.Payload))
	}

}

func (a *orbAgent) handleRPCFromCore(client mqtt.Client, message mqtt.Message) {

	a.logger.Debug("RPC message from core", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))

	var rpc fleet.RPC
	if err := json.Unmarshal(message.Payload(), &rpc); err != nil {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaMalformed))
		return
	}
	if rpc.SchemaVersion != fleet.CurrentRPCSchemaVersion {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaVersion))
		return
	}
	if rpc.Func == "" || rpc.Payload == nil {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaMalformed))
		return
	}

	// dispatch
	switch rpc.Func {
	case fleet.GroupMembershipRPCFunc:
		var r fleet.GroupMembershipRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding group membership message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleGroupMembership(r.Payload)
	case fleet.AgentPolicyRPCFunc:
		var r fleet.AgentPolicyRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding agent policy message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleAgentPolicies(r.Payload)
	default:
		a.logger.Warn("unsupported/unhandled core RPC, ignoring",
			zap.String("func", rpc.Func),
			zap.Any("payload", rpc.Payload))
	}

}
