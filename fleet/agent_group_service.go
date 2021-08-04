// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrCreateAgentGroup = errors.New("failed to create agent group")

	ErrMaintainAgentGroupChannels = errors.New("failed to maintain agent group channels")
)

func (svc fleetService) addAgentsToAgentGroupChannel(token string, g AgentGroup) error {

	// first we get all agents, online or not, to connect them to the correct group channel
	list, err := svc.agentRepo.RetrieveAllByAgentGroupID(context.Background(), g.MFOwnerID, g.ID, false)
	if len(list) == 0 {
		return nil
	}
	idList := make([]string, len(list))
	for i, agent := range list {
		idList[i] = agent.MFThingID
	}
	ids := mfsdk.ConnectionIDs{
		ChannelIDs: []string{g.MFChannelID},
		ThingIDs:   idList,
	}
	err = svc.mfsdk.Connect(ids, token)
	if err != nil {
		return err
	}

	// now we get only onlinish agents to notify them in real time
	list, err = svc.agentRepo.RetrieveAllByAgentGroupID(context.Background(), g.MFOwnerID, g.ID, true)
	if err != nil {
		return err
	}
	for _, agent := range list {
		err := svc.agentComms.NotifyNewAgentGroupMembership(agent, g)
		if err != nil {
			// note we will not make failure to deliver to one agent fatal, just log
			svc.logger.Error("failure during agent group membership comms", zap.Error(err))
		}
	}
	return nil
}

func (svc fleetService) RetrieveAgentGroupByID(ctx context.Context, token string, id string) (AgentGroup, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return AgentGroup{}, err
	}
	ag, err := svc.agentGroupRepository.RetrieveByID(ctx, id, ownerID)
	if err != nil {
		return AgentGroup{}, err
	}
	return ag, nil
}

func (svc fleetService) RetrieveAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (AgentGroup, error) {
	return svc.agentGroupRepository.RetrieveByID(ctx, groupID, ownerID)
}

func (svc fleetService) CreateAgentGroup(ctx context.Context, token string, s AgentGroup) (AgentGroup, error) {
	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return AgentGroup{}, err
	}

	s.MFOwnerID = mfOwnerID

	md := map[string]interface{}{"type": "orb_agent_group"}

	// create main Group RPC Channel
	mfChannelID, err := svc.mfsdk.CreateChannel(mfsdk.Channel{
		Name:     s.Name.String(),
		Metadata: md,
	}, token)
	if err != nil {
		return AgentGroup{}, errors.Wrap(ErrCreateAgent, err)
	}

	s.MFChannelID = mfChannelID

	id, err := svc.agentGroupRepository.Save(ctx, s)
	if err != nil {
		return AgentGroup{}, errors.Wrap(ErrCreateAgentGroup, err)
	}

	s.ID = id
	err = svc.addAgentsToAgentGroupChannel(token, s)
	if err != nil {
		// TODO should we roll back?
		svc.logger.Error("error adding agents to group channel", zap.Error(errors.Wrap(ErrMaintainAgentGroupChannels, err)))
	}

	return s, err
}
