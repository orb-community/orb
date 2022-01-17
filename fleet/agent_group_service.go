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

func (svc fleetService) removeAgentGroupSubscriptions(groupID string, ownerID string) error {
	ag, err := svc.agentGroupRepository.RetrieveByID(context.Background(), groupID, ownerID)
	if err != nil {
		return err
	}
	err = svc.agentComms.NotifyGroupRemoval(ag)
	if err != nil {
		svc.logger.Error("failure during agent group membership comms", zap.Error(err))
	}
	return nil
}

func (svc fleetService) addAgentsToAgentGroupChannel(token string, g AgentGroup) error {
	// first we get all agents, online or not, to connect them to the correct group channel
	list, err := svc.agentRepo.RetrieveAllByAgentGroupID(context.Background(), g.MFOwnerID, g.ID, false)
	if len(list) == 0 {
		return nil
	}
	// TODO need to connect a agent on creatinon to a existing agent group
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
		err := svc.agentComms.NotifyAgentNewGroupMembership(agent, g)
		if err != nil {
			// note we will not make failure to deliver to one agent fatal, just log
			svc.logger.Error("failure during agent group membership comms", zap.Error(err))
		}
	}
	return nil
}

func (svc fleetService) ListAgentGroups(ctx context.Context, token string, pm PageMetadata) (PageAgentGroup, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return PageAgentGroup{}, err
	}

	ag, err := svc.agentGroupRepository.RetrieveAllAgentGroupsByOwner(ctx, ownerID, pm)
	if err != nil {
		return PageAgentGroup{}, err
	}
	return ag, nil
}

func (svc fleetService) ViewAgentGroupByID(ctx context.Context, token string, id string) (AgentGroup, error) {
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

func (svc fleetService) EditAgentGroup(ctx context.Context, token string, group AgentGroup) (AgentGroup, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return AgentGroup{}, err
	}

	if len(group.MatchingAgents) > 0 {
		return AgentGroup{}, errors.ErrUpdateEntity
	}

	ag, err := svc.agentGroupRepository.Update(ctx, ownerID, group)
	if err != nil {
		return AgentGroup{}, err
	}

	list, err := svc.agentRepo.RetrieveAllByAgentGroupID(context.Background(), ownerID, group.ID, true)
	if err != nil {
		return AgentGroup{}, err
	}
	for _, agent := range list {
		err := svc.agentComms.NotifyAgentGroupMemberships(agent)
		if err != nil {
			svc.logger.Error("failure during agent group membership comms", zap.Error(err))
		}
	}

	return ag, nil
}

func (svc fleetService) ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (AgentGroup, error) {
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

	ag, err := svc.agentGroupRepository.RetrieveByID(ctx, id, mfOwnerID)
	if err != nil {
		return AgentGroup{}, errors.Wrap(ErrCreateAgentGroup, err)
	}

	err = svc.addAgentsToAgentGroupChannel(token, ag)
	if err != nil {
		// TODO should we roll back?
		svc.logger.Error("error adding agents to group channel", zap.Error(errors.Wrap(ErrMaintainAgentGroupChannels, err)))
	}

	return ag, err
}

func (svc fleetService) RemoveAgentGroup(ctx context.Context, token, groupId string) error {
	ownerID, err := svc.identify(token)
	if err != nil {
		return err
	}

	err = svc.removeAgentGroupSubscriptions(groupId, ownerID)
	if err != nil {
		svc.logger.Error("removing agents from group channel", zap.Error(errors.Wrap(ErrMaintainAgentGroupChannels, err)))
	}

	err = svc.agentGroupRepository.Delete(ctx, groupId, ownerID)
	if err != nil {
		return err
	}

	return nil
}

func (svc fleetService) ValidateAgentGroup(ctx context.Context, token string, ag AgentGroup) (AgentGroup, error) {
	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return AgentGroup{}, err
	}

	ag.MFOwnerID = mfOwnerID
	res, err := svc.agentRepo.RetrieveMatchingAgents(ctx, mfOwnerID, ag.Tags)
	if err != nil {
		return AgentGroup{}, err
	}
	ag.MatchingAgents = res
	return ag, err
}
