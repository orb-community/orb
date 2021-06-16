// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"encoding/json"
	"fmt"
	"github.com/mainflux/mainflux/pkg/messaging"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"go.uber.org/zap"
)

type AgentCommsService interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error
}

var _ AgentCommsService = (*fleetCommsService)(nil)

const CapabilitiesChannel = "agent"
const HeartbeatsChannel = "hb"
const RPCToCoreChannel = "tocore"
const RPCFromCoreChannel = "fromcore"
const LogChannel = "log"

type fleetCommsService struct {
	logger *zap.Logger
	// agent comms
	agentPubSub mfnats.PubSub
}

func NewFleetCommsService(logger *zap.Logger, agentPubSub mfnats.PubSub) AgentCommsService {
	return &fleetCommsService{
		logger:      logger,
		agentPubSub: agentPubSub,
	}
}

func (svc fleetCommsService) handleMsgFromAgent(msg messaging.Message) error {
	var payload interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	svc.logger.Info("received message", zap.Any("payload", payload), zap.Any("subtopic", msg.Subtopic), zap.Any("channel", msg.Channel),
		zap.Any("protocol", msg.Protocol), zap.Any("created", msg.Created), zap.Any("publisher", msg.Publisher))
	return nil
}

func (svc fleetCommsService) Start() error {
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", CapabilitiesChannel), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", HeartbeatsChannel), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", RPCToCoreChannel), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", LogChannel), svc.handleMsgFromAgent); err != nil {
		return err
	}
	svc.logger.Info("subscribed to agent channels")
	return nil
}

func (svc fleetCommsService) Stop() error {
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", CapabilitiesChannel)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", HeartbeatsChannel)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", RPCToCoreChannel)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", LogChannel)); err != nil {
		return err
	}
	svc.logger.Info("unsubscribed from agent channels")
	return nil
}
