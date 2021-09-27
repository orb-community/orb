// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mainflux/mainflux/pkg/messaging"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/policies/pb"
	"go.uber.org/zap"
	"time"
)

const publisher = "orb-fleet"

type AgentCommsService interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error

	// NotifyNewAgentGroupMembership RPC Core -> Agent: Notify Agent of new AgentGroup membership
	NotifyNewAgentGroupMembership(a Agent, ag AgentGroup) error
	// NotifyAgentGroupMembership RPC Core -> Agent: Notify Agent of all AgentGroup memberships
	NotifyAgentGroupMembership(a Agent) error
	// NotifyAgentAllDatasets RPC Core -> Agent: Notify Agent of all Policy it should currently run based on group membership and current Datasets
	NotifyAgentAllDatasets(a Agent) error
	// NotifyGroupNewDataset RPC Core -> Agent: Notify AgentGroup of a newly created Dataset, exposing a new Policy to run
	NotifyGroupNewDataset(ctx context.Context, ag AgentGroup, datasetID string, policyID string, ownerID string) error
	// NotifyGroupRemoval unsubscribe the agent membership when delete a agent group
	NotifyGroupRemoval(ag AgentGroup) error
	// NotifyPolicyRemoval stop agent policy utilization after Policy removal
	NotifyPolicyRemoval(policyID string, ag AgentGroup) error
	// NotifyDatasetRemoval usubscribe the agent membership when delete a dataset
	NofityDatasetRemoval(ag AgentGroup, dsID string) error
}

var _ AgentCommsService = (*fleetCommsService)(nil)

const CapabilitiesTopic = "agent"
const HeartbeatsTopic = "hb"
const RPCToCoreTopic = "tocore"
const RPCFromCoreTopic = "fromcore"
const LogTopic = "log"

type fleetCommsService struct {
	logger         *zap.Logger
	agentRepo      AgentRepository
	agentGroupRepo AgentGroupRepository
	policyClient   pb.PolicyServiceClient

	// agent comms
	agentPubSub mfnats.PubSub
}

func (svc fleetCommsService) NotifyGroupNewDataset(ctx context.Context, ag AgentGroup, datasetID string, policyID string, ownerID string) error {
	p, err := svc.policyClient.RetrievePolicy(ctx, &pb.PolicyByIDReq{PolicyID: policyID, OwnerID: ownerID})
	if err != nil {
		return err
	}

	var pdata interface{}
	if err := json.Unmarshal(p.Data, &pdata); err != nil {
		return err
	}

	payload := []AgentPolicyRPCPayload{{
		Action:    "manage",
		ID:        policyID,
		Name:      p.Name,
		Backend:   p.Backend,
		Version:   p.Version,
		Data:      pdata,
		DatasetID: datasetID,
	}}

	data := RPC{
		SchemaVersion: CurrentRPCSchemaVersion,
		Func:          AgentPolicyRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := messaging.Message{
		Channel:   ag.MFChannelID,
		Subtopic:  RPCFromCoreTopic,
		Publisher: publisher,
		Payload:   body,
		Created:   time.Now().UnixNano(),
	}
	if err := svc.agentPubSub.Publish(msg.Channel, msg); err != nil {
		return err
	}

	return nil
}

func (svc fleetCommsService) NotifyNewAgentGroupMembership(a Agent, ag AgentGroup) error {
	payload := GroupMembershipRPCPayload{
		Groups:   []GroupMembershipData{{Name: ag.Name.String(), ChannelID: ag.MFChannelID}},
		FullList: false,
	}

	data := RPC{
		SchemaVersion: CurrentRPCSchemaVersion,
		Func:          GroupMembershipRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := messaging.Message{
		Channel:   a.MFChannelID,
		Subtopic:  RPCFromCoreTopic,
		Publisher: publisher,
		Payload:   body,
		Created:   time.Now().UnixNano(),
	}
	if err := svc.agentPubSub.Publish(msg.Channel, msg); err != nil {
		return err
	}

	return nil

}

func (svc fleetCommsService) NotifyAgentAllDatasets(a Agent) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	groups, err := svc.agentGroupRepo.RetrieveAllByAgent(ctx, a)
	if err != nil {
		return err
	}

	if len(groups) == 0 {
		// no groups, nothing to do
		return nil
	}

	groupIDs := make([]string, len(groups))
	for i, group := range groups {
		groupIDs[i] = group.ID
	}

	// MQTT he doesn't have OwnerID, we need to look it up
	a, err = svc.agentRepo.RetrieveByIDWithChannel(ctx, a.MFThingID, a.MFChannelID)
	if err != nil {
		return err
	}

	p, err := svc.policyClient.RetrievePoliciesByGroups(ctx, &pb.PoliciesByGroupsReq{GroupIDs: groupIDs, OwnerID: a.MFOwnerID})
	if err != nil {
		return err
	}

	payload := make([]AgentPolicyRPCPayload, len(p.Policies))
	for i, policy := range p.Policies {

		var pdata interface{}
		if err := json.Unmarshal(policy.Data, &pdata); err != nil {
			return err
		}

		payload[i] = AgentPolicyRPCPayload{
			Action:    "manage",
			ID:        policy.Id,
			Name:      policy.Name,
			Backend:   policy.Backend,
			Version:   policy.Version,
			Data:      pdata,
			DatasetID: policy.DatasetId,
		}

	}

	data := RPC{
		SchemaVersion: CurrentRPCSchemaVersion,
		Func:          AgentPolicyRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := messaging.Message{
		Channel:   a.MFChannelID,
		Subtopic:  RPCFromCoreTopic,
		Publisher: publisher,
		Payload:   body,
		Created:   time.Now().UnixNano(),
	}
	if err := svc.agentPubSub.Publish(msg.Channel, msg); err != nil {
		return err
	}

	return nil
}

func (svc fleetCommsService) NotifyAgentGroupMembership(a Agent) error {

	list, err := svc.agentGroupRepo.RetrieveAllByAgent(context.Background(), a)
	if err != nil {
		return err
	}

	fullList := make([]GroupMembershipData, len(list))
	for i, agentGroup := range list {
		fullList[i].Name = agentGroup.Name.String()
		fullList[i].ChannelID = agentGroup.MFChannelID
	}

	payload := GroupMembershipRPCPayload{
		Groups:   fullList,
		FullList: true,
	}

	data := RPC{
		SchemaVersion: CurrentRPCSchemaVersion,
		Func:          GroupMembershipRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := messaging.Message{
		Channel:   a.MFChannelID,
		Subtopic:  RPCFromCoreTopic,
		Publisher: publisher,
		Payload:   body,
		Created:   time.Now().UnixNano(),
	}
	if err := svc.agentPubSub.Publish(msg.Channel, msg); err != nil {
		return err
	}

	return nil

}

func (svc fleetCommsService) NotifyGroupRemoval(ag AgentGroup) error {

	data := RPC{
		SchemaVersion: CurrentRPCSchemaVersion,
		Func:          GroupRemovedRPCFunc,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := messaging.Message{
		Channel:   ag.MFChannelID,
		Subtopic:  RPCFromCoreTopic,
		Publisher: publisher,
		Payload:   body,
		Created:   time.Now().UnixNano(),
	}
	if err := svc.agentPubSub.Publish(msg.Channel, msg); err != nil {
		return err
	}
	return nil
}

func (svc fleetCommsService) NotifyPolicyRemoval(policyID string, ag AgentGroup) error {

	payload := AgentPolicyRPCPayload{
		Action: "remove",
		ID:     policyID,
	}

	data := RPC{
		SchemaVersion: CurrentRPCSchemaVersion,
		Func:          AgentPolicyRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := messaging.Message{
		Channel:   ag.MFChannelID,
		Subtopic:  RPCFromCoreTopic,
		Publisher: publisher,
		Payload:   body,
		Created:   time.Now().UnixNano(),
	}
	if err := svc.agentPubSub.Publish(msg.Channel, msg); err != nil {
		return err
	}
	return nil
}

func (svc fleetCommsService) NofityDatasetRemoval(ag AgentGroup, dsID string) error {

	payload := DatasetRemovedRPCPayload{DatasetID: dsID}

	data := RPC{
		SchemaVersion: CurrentRPCSchemaVersion,
		Func:          DatasetRemovedRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := messaging.Message{
		Channel:   ag.MFChannelID,
		Subtopic:  RPCFromCoreTopic,
		Publisher: publisher,
		Payload:   body,
		Created:   time.Now().UnixNano(),
	}
	if err := svc.agentPubSub.Publish(msg.Channel, msg); err != nil {
		return err
	}
	return nil
}

func NewFleetCommsService(logger *zap.Logger, policyClient pb.PolicyServiceClient, agentRepo AgentRepository, agentGroupRepo AgentGroupRepository, agentPubSub mfnats.PubSub) AgentCommsService {
	return &fleetCommsService{
		logger:         logger,
		agentRepo:      agentRepo,
		agentGroupRepo: agentGroupRepo,
		agentPubSub:    agentPubSub,
		policyClient:   policyClient,
	}
}

func (svc fleetCommsService) handleCapabilities(thingID string, channelID string, payload []byte) error {
	var versionCheck SchemaVersionCheck
	if err := json.Unmarshal(payload, &versionCheck); err != nil {
		return ErrSchemaMalformed
	}
	if versionCheck.SchemaVersion != CurrentCapabilitiesSchemaVersion {
		return ErrSchemaVersion
	}
	var capabilities Capabilities
	if err := json.Unmarshal(payload, &capabilities); err != nil {
		return ErrSchemaMalformed
	}
	agent := Agent{MFThingID: thingID, MFChannelID: channelID}
	agent.AgentMetadata = make(map[string]interface{})
	agent.AgentMetadata["backends"] = capabilities.Backends
	agent.AgentMetadata["orb_agent"] = capabilities.OrbAgent
	agent.AgentTags = capabilities.AgentTags
	err := svc.agentRepo.UpdateDataByIDWithChannel(context.Background(), agent)
	if err != nil {
		return err
	}
	return nil
}

func (svc fleetCommsService) handleHeartbeat(thingID string, channelID string, payload []byte) error {
	var versionCheck SchemaVersionCheck
	if err := json.Unmarshal(payload, &versionCheck); err != nil {
		return ErrSchemaMalformed
	}
	if versionCheck.SchemaVersion != CurrentHeartbeatSchemaVersion {
		return ErrSchemaVersion
	}
	var hb Heartbeat
	if err := json.Unmarshal(payload, &hb); err != nil {
		return ErrSchemaMalformed
	}
	agent := Agent{MFThingID: thingID, MFChannelID: channelID}
	agent.LastHBData = make(map[string]interface{})
	// accept "offline" state request to indicate agent is going offline
	if hb.State == Offline {
		agent.State = Offline
		agent.LastHBData["backend_state"] = BackendStateInfo{}
		agent.LastHBData["policy_state"] = PolicyStateInfo{}
	} else {
		// otherwise, state is always "online"
		agent.State = Online
		agent.LastHBData["backend_state"] = hb.BackendState
		agent.LastHBData["policy_state"] = hb.PolicyState
	}
	err := svc.agentRepo.UpdateHeartbeatByIDWithChannel(context.Background(), agent)
	if err != nil {
		return err
	}
	return nil
}

func (svc fleetCommsService) handleRPCToCore(thingID string, channelID string, payload []byte) error {
	var versionCheck SchemaVersionCheck
	if err := json.Unmarshal(payload, &versionCheck); err != nil {
		return ErrSchemaMalformed
	}
	if versionCheck.SchemaVersion != CurrentRPCSchemaVersion {
		return ErrSchemaVersion
	}
	var rpc RPC
	if err := json.Unmarshal(payload, &rpc); err != nil {
		return ErrSchemaMalformed
	}

	// dispatch
	switch rpc.Func {
	case GroupMembershipReqRPCFunc:
		if err := svc.NotifyAgentGroupMembership(Agent{MFThingID: thingID, MFChannelID: channelID}); err != nil {
			svc.logger.Error("notify group membership failure", zap.Error(err))
			return nil
		}
	case AgentPoliciesReqRPCFunc:
		if err := svc.NotifyAgentAllDatasets(Agent{MFThingID: thingID, MFChannelID: channelID}); err != nil {
			svc.logger.Error("notify agent policies failure", zap.Error(err))
			return nil
		}
	default:
		svc.logger.Warn("unsupported/unhandled agent RPC, ignoring",
			zap.String("func", rpc.Func),
			zap.Any("payload", rpc.Payload))
	}

	return nil
}

func (svc fleetCommsService) handleMsgFromAgent(msg messaging.Message) error {

	// NOTE: we need to consider ALL input from the agent as untrusted, the same as untrusted HTTP API would be
	// Given security context is that to get this far we know mainflux MQTT proxy has authenticated a
	// username/password/channelID combination (thingID/thingKey/thingChannel which are all UUIDv4)
	// channelID is globally unique across all owners and things, and can therefore substitute for an ownerID (which we do not have here)
	// mainflux will not allow a thing to communicate on a channelID it does not belong to - thus it is not possible
	// to brute force a channelID from another tenant without brute forcing all three UUIDs which is a lot of entropy

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}

	svc.logger.Debug("received agent message",
		zap.Any("payload", payload),
		zap.String("subtopic", msg.Subtopic),
		zap.String("channel", msg.Channel),
		zap.String("protocol", msg.Protocol),
		zap.Int64("created", msg.Created),
		zap.String("publisher", msg.Publisher))

	if len(msg.Payload) > MaxMsgPayloadSize {
		return ErrPayloadTooBig
	}

	// dispatch
	switch msg.Subtopic {
	case CapabilitiesTopic:
		if err := svc.handleCapabilities(msg.Publisher, msg.Channel, msg.Payload); err != nil {
			svc.logger.Error("capabilities failure", zap.Error(err))
			return err
		}
	case HeartbeatsTopic:
		if err := svc.handleHeartbeat(msg.Publisher, msg.Channel, msg.Payload); err != nil {
			svc.logger.Error("heartbeat failure", zap.Error(err))
			return err
		}
	case RPCToCoreTopic:
		if err := svc.handleRPCToCore(msg.Publisher, msg.Channel, msg.Payload); err != nil {
			svc.logger.Error("RPC to core failure", zap.Error(err))
			return err
		}
	case LogTopic:
		svc.logger.Error("implement me: LogChannel")
	default:
		svc.logger.Warn("unsupported/unhandled agent subtopic, ignoring",
			zap.String("subtopic", msg.Subtopic),
			zap.String("thing_id", msg.Publisher),
			zap.String("channel_id", msg.Channel))
	}

	return nil
}

func (svc fleetCommsService) Start() error {
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", CapabilitiesTopic), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", HeartbeatsTopic), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", RPCToCoreTopic), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", LogTopic), svc.handleMsgFromAgent); err != nil {
		return err
	}
	svc.logger.Info("subscribed to agent channels")
	return nil
}

func (svc fleetCommsService) Stop() error {
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", CapabilitiesTopic)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", HeartbeatsTopic)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", RPCToCoreTopic)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", LogTopic)); err != nil {
		return err
	}
	svc.logger.Info("unsubscribed from agent channels")
	return nil
}
