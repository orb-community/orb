// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"context"
	"encoding/json"
	"github.com/mainflux/mainflux/pkg/messaging"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/fleet/pb"
	"go.uber.org/zap"
	"time"
)

const publisher = "orb_policies"

type PolicyCommsService interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error
	//NotifyDatasetPolicyUpdate RPC Core -> Datasets
	NotifyDatasetPolicyUpdate(ctx context.Context, policy Policy, dataset Dataset) error
}

var _ PolicyCommsService = (*policiesCommsService)(nil)

const RPCFromCoreTopic = "fromcore"

type policiesCommsService struct {
	logger      *zap.Logger
	fleetClient pb.FleetServiceClient
	// policy comms
	policyPubSub mfnats.PubSub
}

func (p policiesCommsService) Start() error {
	return nil
}

func (p policiesCommsService) Stop() error {
	return nil
}

func (p policiesCommsService) NotifyDatasetPolicyUpdate(ctx context.Context, policy Policy, dataset Dataset) error {
	ag, err := p.fleetClient.RetrieveAgentGroup(ctx, &pb.AgentGroupByIDReq{
		AgentGroupID: dataset.AgentGroupID,
		OwnerID:      dataset.MFOwnerID,
	})

	payload := []DatasetRPCPayload{{
		ID:            dataset.ID,
		Name:          dataset.Name.String(),
		AgentGroupID:  dataset.AgentGroupID,
		AgentPolicyID: dataset.PolicyID,
		SinkID:        dataset.SinkID,
		Data:          policy,
	}}

	data := RPC{
		SchemaVersion: CurrentRPCSchemaVersion,
		Func:          DatasetReqRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := messaging.Message{
		Channel:   ag.GetChannel(),
		Subtopic:  RPCFromCoreTopic,
		Publisher: publisher,
		Payload:   body,
		Created:   time.Now().UnixNano(),
	}
	if err := p.policyPubSub.Publish(msg.Channel, msg); err != nil {
		return err
	}

	return nil
}

func NewPoliciesCommsService(logger *zap.Logger, fleetClient pb.FleetServiceClient, policyPubSub mfnats.PubSub) PolicyCommsService {
	return &policiesCommsService{
		logger:       logger,
		fleetClient:  fleetClient,
		policyPubSub: policyPubSub,
	}
}
