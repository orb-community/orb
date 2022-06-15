/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

func (a *orbAgent) handleGroupMembership(rpc fleet.GroupMembershipRPCPayload) {
	// if this is the full list, reset all group subscriptions and subscribed to this list
	if rpc.FullList {
		a.unsubscribeGroupChannels()
		a.subscribeGroupChannels(rpc.Groups)
		err := a.sendAgentPoliciesReq()
		if err != nil {
			a.logger.Error("failed to send agent policies request", zap.Error(err))
		}
	} else {
		// otherwise, just add these subscriptions to the existing list
		a.subscribeGroupChannels(rpc.Groups)
	}
}

func (a *orbAgent) handleAgentPolicies(rpc []fleet.AgentPolicyRPCPayload, fullList bool) {
	if fullList {
		policies, err := a.policyManager.GetRepo().GetAll()
		if err != nil {
			a.logger.Error("failed to retrieve policies on handle subscriptions")
			return
		}
		// Create a map with all the old policies
		policyRemove := map[string]bool{}
		for _, p := range policies {
			policyRemove[p.ID] = true
		}
		for _, payload := range rpc {
			if ok := policyRemove[payload.ID]; ok {
				policyRemove[payload.ID] = false
			}
		}
		// Remove only the policy which should be removed
		for k, v := range policyRemove {
			if v == true {
				policy, err := a.policyManager.GetRepo().Get(k)
				if err != nil {
					a.logger.Warn("failed to retrieve policy", zap.String("policy_id", k), zap.Error(err))
					continue
				}
				err = a.policyManager.RemovePolicy(policy.ID, policy.Name, policy.Backend)
				if err != nil {
					a.logger.Warn("failed to remove a policy, ignoring", zap.String("policy_id", policy.ID), zap.String("policy_name", policy.Name), zap.Error(err))
					continue
				}
			}
		}
	}

	for _, payload := range rpc {
		if payload.Action != "sanitize" {
			a.policyManager.ManagePolicy(payload)
		}
	}

	// heart beat with new policy status after application
	a.sendSingleHeartbeat(time.Now(), fleet.Online)
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
		a.handleAgentPolicies(r.Payload, r.FullList)
	case fleet.GroupRemovedRPCFunc:
		var r fleet.GroupRemovedRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding agent group removal message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleAgentGroupRemoval(r.Payload)
	case fleet.DatasetRemovedRPCFunc:
		var r fleet.DatasetRemovedRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding dataset removal message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleDatasetRemoval(r.Payload)
	default:
		a.logger.Warn("unsupported/unhandled core RPC, ignoring",
			zap.String("func", rpc.Func),
			zap.Any("payload", rpc.Payload))
	}

}

func (a *orbAgent) handleAgentStop(payload fleet.AgentStopRPCPayload) {
	// TODO graceful stop agent https://github.com/ns1labs/orb/issues/466
	panic(fmt.Sprintf("control plane requested we terminate, reason: %s", payload.Reason))
}

func (a *orbAgent) handleAgentGroupRemoval(rpc fleet.GroupRemovedRPCPayload) {
	a.unsubscribeGroupChannel(rpc.ChannelID)
}

func (a *orbAgent) handleDatasetRemoval(rpc fleet.DatasetRemovedRPCPayload) {
	a.removeDatasetFromPolicy(rpc.DatasetID, rpc.PolicyID)
}

func (a *orbAgent) handleAgentReset(payload fleet.AgentResetRPCPayload) {
	if payload.FullReset {
		err := a.RestartAll(payload.Reason)
		if err != nil {
			a.logger.Error("RestartAll failure", zap.Error(err))
		}
	} else {
		// TODO backend specific restart
		// a.RestartBackend()
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
		a.logger.Debug("received group membership, setting success in request")
		a.groupRequestSucceeded <- true
	case fleet.AgentPolicyRPCFunc:
		var r fleet.AgentPolicyRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding agent policy message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleAgentPolicies(r.Payload, r.FullList)
		a.logger.Debug("received agent policies, setting success in request")
		a.policyRequestSucceeded <- true
	case fleet.AgentStopRPCFunc:
		var r fleet.AgentStopRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding agent stop message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleAgentStop(r.Payload)
	case fleet.AgentResetRPCFunc:
		var r fleet.AgentResetRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding agent reset message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleAgentReset(r.Payload)
	default:
		a.logger.Warn("unsupported/unhandled core RPC, ignoring",
			zap.String("func", rpc.Func),
			zap.Any("payload", rpc.Payload))
	}

}
