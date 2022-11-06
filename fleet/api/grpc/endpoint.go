package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/etaques/orb/fleet"
)

func retrieveAgentEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(accessByIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		agent, err := svc.ViewAgentByIDInternal(ctx, req.OwnerID, req.AgentID)
		if err != nil {
			return nil, err
		}
		res := agentRes{
			id:   agent.MFThingID,
			name: agent.Name.String(),
		}
		return res, nil
	}
}

func retrieveAgentGroupEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(accessAgByIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		agentGroup, err := svc.ViewAgentGroupByIDInternal(ctx, req.AgentGroupID, req.OwnerID)
		if err != nil {
			return nil, err
		}
		res := agentGroupRes{
			id:      agentGroup.ID,
			name:    agentGroup.Name.String(),
			channel: agentGroup.MFChannelID,
		}
		return res, nil
	}
}

func retrieveOwnerByChannelIDEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(accessOwnerByChannelIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		agent, err := svc.ViewAgentInfoByChannelIDInternal(ctx, req.ChannelID)
		if err != nil {
			return nil, err
		}
		res := ownerRes{ownerID: agent.MFOwnerID, agentName: agent.Name.String()}
		return res, nil
	}
}

func retrieveAgentInfoByChannelIDEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(accessAgentInfoByChannelIDReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		agent, err := svc.ViewAgentInfoByChannelIDInternal(ctx, req.ChannelID)
		if err != nil {
			return nil, err
		}

		matchingGroups, err := svc.ViewAgentMatchingGroupsByIDInternal(ctx, agent.MFThingID, agent.MFOwnerID)
		if err != nil {
			return nil, err
		}

		var groupIDs []string
		for _, group := range matchingGroups.Groups {
			groupIDs = append(groupIDs, group.GroupID)
		}

		res := agentInfoRes{
			ownerID:       agent.MFOwnerID,
			agentName:     agent.Name.String(),
			agentTags:     agent.AgentTags,
			orbTags:       agent.OrbTags,
			agentGroupIDs: groupIDs,
		}
		return res, nil
	}
}
