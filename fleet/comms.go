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

	// agent comms
	agentPubSub mfnats.PubSub
}

func (svc fleetCommsService) NotifyNewAgentGroupMembership(a Agent, ag AgentGroup) error {

	payload := GroupMembershipRPCPayload{
		ChannelIDS: []string{ag.MFChannelID},
		FullList:   false,
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

func (svc fleetCommsService) NotifyAgentGroupMembership(a Agent) error {

	list, err := svc.agentGroupRepo.RetrieveAllByAgent(context.Background(), a)
	if err != nil {
		return err
	}

	fullList := make([]string, len(list))
	for i, agentGroup := range list {
		fullList[i] = agentGroup.MFChannelID
	}

	payload := GroupMembershipRPCPayload{
		ChannelIDS: fullList,
		FullList:   true,
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

func NewFleetCommsService(logger *zap.Logger, agentRepo AgentRepository, agentGroupRepo AgentGroupRepository, agentPubSub mfnats.PubSub) AgentCommsService {
	return &fleetCommsService{
		logger:         logger,
		agentRepo:      agentRepo,
		agentGroupRepo: agentGroupRepo,
		agentPubSub:    agentPubSub,
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
	agent.LastHBData["ts"] = hb.TimeStamp.UnixNano()
	// accept "offline" state request to indicate agent is going offline
	if hb.State == Offline {
		agent.State = Offline
	} else {
		// otherwise, state is always "online"
		agent.State = Online
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
			svc.logger.Error(" failure", zap.Error(err))
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
