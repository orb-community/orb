package bridgeservice

import (
	"context"
	fleetpb "github.com/ns1labs/orb/fleet/pb"
	policiespb "github.com/ns1labs/orb/policies/pb"
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
)

type BridgeService interface {
	ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error)
	GetDataSetsFromAgentGroups(ctx context.Context, mfOwnerId string, agentGroupIds []string) (map[string]string, error)
}

type SinkerOtelBridgeService struct {
	logger         *zap.Logger
	sinkerCache    config.ConfigRepo
	policiesClient policiespb.PolicyServiceClient
	fleetClient    fleetpb.FleetServiceClient
}

func (bs *SinkerOtelBridgeService) ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error) {
	agentPb, err := bs.fleetClient.RetrieveAgentInfoByChannelID(ctx, &fleetpb.AgentInfoByChannelIDReq{Channel: channelID})
	if err != nil {
		return nil, err
	}
	return agentPb, nil
}

func (bs *SinkerOtelBridgeService) GetSinkIdsFromAgentGroups(ctx context.Context, mfOwnerId string, agentGroupIds []string) (map[string]string, error) {
	policiesRes, err := bs.policiesClient.RetrievePoliciesByGroups(ctx, &policiespb.PoliciesByGroupsReq{
		GroupIDs: agentGroupIds,
		OwnerID:  mfOwnerId,
	})
	if err != nil {
		bs.logger.Error("unable to retrieve policies from agent groups", zap.Error(err))
		return nil, err
	}
	mapSinkIdPolicy := make(map[string]string)
	for _, policy := range policiesRes.Policies {
		datasetRes, err := bs.policiesClient.RetrieveDataset(ctx, &policiespb.DatasetByIDReq{
			DatasetID: policy.DatasetId,
			OwnerID:   mfOwnerId,
		})
		if err != nil {
			bs.logger.Error("unable to retrieve datasets from policy", zap.String("policy", policy.Name), zap.Error(err))
			continue
		}
		for _, sinkId := range datasetRes.SinkIds {
			mapSinkIdPolicy[sinkId] = "active"
		}
	}
	return mapSinkIdPolicy, nil
}
