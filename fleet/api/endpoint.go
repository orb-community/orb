// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/ns1labs/orb/fleet"
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
			Name: nID,
			Tags: req.Tags,
		}
		saved, err := svc.CreateAgentGroup(c, req.token, group)
		if err != nil {
			return nil, err
		}

		res := agentGroupRes{
			ID:      saved.ID,
			Name:    saved.Name.String(),
			Tags:    saved.Tags,
			created: true,
		}

		return res, nil
	}
}

func viewAgentGroupEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return response, fleet.ErrNotFound
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
			Name:    nID,
			OrbTags: req.OrbTags,
		}
		saved, err := svc.CreateAgent(c, req.token, agent)
		if err != nil {
			return nil, err
		}

		res := agentRes{
			Name:      saved.Name.String(),
			ID:        saved.MFThingID,
			State:     saved.State.String(),
			Key:       saved.MFKeyID,
			ChannelID: saved.MFChannelID,
			created:   true,
		}

		return res, nil
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
			Agents: []viewAgentRes{},
		}
		for _, agent := range page.Agents {
			view := viewAgentRes{
				ID:           agent.MFThingID,
				ChannelID:    agent.MFChannelID,
				Owner:        agent.MFOwnerID,
				Name:         agent.Name.String(),
				State:        agent.State.String(),
				Capabilities: agent.AgentMetadata,
			}
			res.Agents = append(res.Agents, view)
		}

		return res, nil
	}
}
