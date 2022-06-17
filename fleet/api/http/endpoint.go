// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
)

func addAgentGroupEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(c context.Context, request interface{}) (interface{}, error) {
		req := request.(addAgentGroupReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		group := fleet.AgentGroup{
			Name:        nID,
			Description: req.Description,
			Tags:        req.Tags,
		}
		saved, err := svc.CreateAgentGroup(c, req.token, group)
		if err != nil {
			return nil, err
		}

		res := agentGroupRes{
			ID:             saved.ID,
			Name:           saved.Name.String(),
			Description:    saved.Description,
			Tags:           saved.Tags,
			MatchingAgents: saved.MatchingAgents,
			created:        true,
		}

		return res, nil
	}
}

func viewAgentGroupEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		agentGroup, err := svc.ViewAgentGroupByID(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}
		res := agentGroupRes{
			ID:             agentGroup.ID,
			Name:           agentGroup.Name.String(),
			Description:    agentGroup.Description,
			Tags:           agentGroup.Tags,
			TsCreated:      agentGroup.Created,
			MatchingAgents: agentGroup.MatchingAgents,
		}
		return res, nil
	}
}

func listAgentGroupsEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listResourcesReq)

		if err := req.validate(); err != nil {
			return nil, err
		}
		page, err := svc.ListAgentGroups(ctx, req.token, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		res := agentGroupsPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
				Order:  page.Order,
				Dir:    page.Dir,
			},
			AgentGroups: []agentGroupRes{},
		}
		for _, ag := range page.AgentGroups {
			view := agentGroupRes{
				ID:             ag.ID,
				Name:           ag.Name.String(),
				Description:    ag.Description,
				Tags:           ag.Tags,
				TsCreated:      ag.Created,
				MatchingAgents: ag.MatchingAgents,
			}
			res.AgentGroups = append(res.AgentGroups, view)
		}
		return res, nil
	}
}

func editAgentGroupEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateAgentGroupReq)
		if err := req.validate(); err != nil {
			return agentGroupRes{}, err
		}

		validName, err := types.NewIdentifier(req.Name)
		if err != nil {
			return agentGroupRes{}, err
		}
		ag := fleet.AgentGroup{
			ID:          req.id,
			Name:        validName,
			Description: req.Description,
			Tags:        req.Tags,
		}

		data, err := svc.EditAgentGroup(ctx, req.token, ag)
		if err != nil {
			return agentGroupRes{}, err
		}

		res := agentGroupRes{
			ID:             data.ID,
			Name:           data.Name.String(),
			Description:    data.Description,
			Tags:           data.Tags,
			TsCreated:      data.Created,
			MatchingAgents: data.MatchingAgents,
		}

		return res, nil
	}
}

func removeAgentGroupEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemoveAgentGroup(ctx, req.token, req.id); err != nil {
			return nil, err
		}
		return removeRes{}, nil
	}
}

func addAgentEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(c context.Context, request interface{}) (interface{}, error) {
		req := request.(addAgentReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		agent := fleet.Agent{
			Name:      nID,
			OrbTags:   req.OrbTags,
			AgentTags: req.AgentTags,
		}
		saved, err := svc.CreateAgent(c, req.token, agent)
		if err != nil {
			return nil, err
		}

		res := agentRes{
			Name:          saved.Name.String(),
			ID:            saved.MFThingID,
			State:         saved.State.String(),
			Key:           saved.MFKeyID,
			OrbTags:       saved.OrbTags,
			AgentTags:     saved.AgentTags,
			AgentMetadata: saved.AgentMetadata,
			LastHBData:    saved.LastHBData,
			TsCreated:     saved.Created,
			created:       true,
			ChannelID:     saved.MFChannelID,
		}

		return res, nil
	}
}

func viewAgentEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		ag, err := svc.ViewAgentByID(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := agentRes{
			ID:            ag.MFThingID,
			Name:          ag.Name.String(),
			ChannelID:     ag.MFChannelID,
			AgentTags:     ag.AgentTags,
			OrbTags:       ag.OrbTags,
			TsCreated:     ag.Created,
			AgentMetadata: ag.AgentMetadata,
			State:         ag.State.String(),
			LastHBData:    ag.LastHBData,
			TsLastHB:      ag.LastHB,
		}
		return res, nil
	}
}

func viewAgentMatchingGroups(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		matchingGroups, err := svc.ViewAgentMatchingGroupsByID(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		var res []matchingGroupsRes
		for _, group := range matchingGroups.Groups {
			res = append(res, matchingGroupsRes{
				GroupID:   group.GroupID,
				GroupName: group.GroupName.String(),
			})
		}
		return res, nil
	}
}

func resetAgentEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if err := svc.ResetAgent(ctx, req.token, req.id); err != nil {
			return nil, err
		}
		return response, nil
	}
}

func listAgentsEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listResourcesReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListAgents(ctx, req.token, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		res := agentsPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
				Order:  page.Order,
				Dir:    page.Dir,
			},
			Agents: []agentRes{},
		}

		for _, ag := range page.Agents {

			policyState, err := svc.GetPoliciesState(ag, nil)
			if err != nil {
				return nil, err
			}

			view := agentRes{
				ID:          ag.MFThingID,
				Name:        ag.Name.String(),
				ChannelID:   ag.MFChannelID,
				AgentTags:   ag.AgentTags,
				OrbTags:     ag.OrbTags,
				TsCreated:   ag.Created,
				State:       ag.State.String(),
				TsLastHB:    ag.LastHB,
				PolicyState: policyState,
			}
			res.Agents = append(res.Agents, view)
		}

		return res, nil
	}
}

func validateAgentGroupEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(c context.Context, request interface{}) (interface{}, error) {
		req := request.(addAgentGroupReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		group := fleet.AgentGroup{
			Name: nID,
			Tags: req.Tags,
		}
		validated, err := svc.ValidateAgentGroup(c, req.token, group)
		if err != nil {
			return nil, err
		}

		res := validateAgentGroupRes{
			Name:           validated.Name.String(),
			Tags:           validated.Tags,
			MatchingAgents: validated.MatchingAgents,
		}

		return res, nil
	}
}

func editAgentEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateAgentReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		validName, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}
		agent := fleet.Agent{
			Name:      validName,
			MFThingID: req.id,
			OrbTags:   req.Tags,
		}

		ag, err := svc.EditAgent(ctx, req.token, agent)
		if err != nil {
			return nil, err
		}

		res := agentRes{
			ID:            ag.MFThingID,
			Name:          ag.Name.String(),
			ChannelID:     ag.MFChannelID,
			AgentTags:     ag.AgentTags,
			OrbTags:       ag.OrbTags,
			TsCreated:     ag.Created,
			AgentMetadata: ag.AgentMetadata,
			State:         ag.State.String(),
			LastHBData:    ag.LastHBData,
			TsLastHB:      ag.LastHB,
		}

		return res, nil

	}
}

func validateAgentEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(c context.Context, request interface{}) (interface{}, error) {
		req := request.(addAgentReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		agent := fleet.Agent{
			Name:    nID,
			OrbTags: req.OrbTags,
		}
		validated, err := svc.ValidateAgent(c, req.token, agent)
		if err != nil {
			return nil, err
		}

		res := validateAgentRes{
			Name:    validated.Name.String(),
			OrbTags: validated.OrbTags,
		}
		return res, nil
	}
}

func removeAgentEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemoveAgent(ctx, req.token, req.id); err != nil {
			return nil, err
		}
		return removeRes{}, nil
	}
}

func listAgentBackendsEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listAgentBackendsReq)
		if err := req.validate(); err != nil {
			return agentBackendsRes{}, err
		}

		bks, err := svc.ListAgentBackends(ctx, req.token)
		if err != nil {
			return agentBackendsRes{}, err
		}

		var res []interface{}
		for _, be := range bks {
			mt, err := svc.ViewAgentBackend(ctx, req.token, be)
			if err != nil {
				return agentBackendsRes{}, err
			}
			if mt == nil {
				return agentBackendsRes{}, errors.ErrNotFound
			}
			res = append(res, mt)
		}

		return agentBackendsRes{
			Backends: res,
		}, nil
	}
}
