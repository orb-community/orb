package bridgeservice

import (
	"context"
	fleetpb "github.com/ns1labs/orb/fleet/pb"
	policiespb "github.com/ns1labs/orb/policies/pb"
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
	"strings"
	"time"
)

type BridgeService interface {
	ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error)
	GetDataSetsFromAgentGroups(ctx context.Context, mfOwnerId string, agentGroupIds []string) (map[string]string, error)
	NotifyActiveSink(ctx context.Context, mfOwnerId, sinkId, state, message string) error
}

func NewBridgeService(logger *zap.Logger,
	sinkerCache config.ConfigRepo,
	policiesClient policiespb.PolicyServiceClient,
	fleetClient fleetpb.FleetServiceClient) SinkerOtelBridgeService {
	return SinkerOtelBridgeService{
		logger:         logger,
		sinkerCache:    sinkerCache,
		policiesClient: policiesClient,
		fleetClient:    fleetClient,
	}
}

type SinkerOtelBridgeService struct {
	logger         *zap.Logger
	sinkerCache    config.ConfigRepo
	policiesClient policiespb.PolicyServiceClient
	fleetClient    fleetpb.FleetServiceClient
}

func (bs *SinkerOtelBridgeService) NotifyActiveSink(_ context.Context, mfOwnerId, sinkId, newState, message string) error {
	cfgRepo, err := bs.sinkerCache.Get(mfOwnerId, sinkId)
	if err != nil {
		bs.logger.Error("unable to retrieve the sink config", zap.Error(err))
		return err
	}
	err = cfgRepo.State.SetFromString(newState)
	if err != nil {
		bs.logger.Error("unable to set state", zap.String("new_state", newState), zap.Error(err))
		return err
	}
	if cfgRepo.State == config.Error {
		cfgRepo.Msg = message
	} else if cfgRepo.State == config.Active {
		cfgRepo.LastRemoteWrite = time.Now()
	}
	err = bs.sinkerCache.Edit(cfgRepo)
	if err != nil {
		bs.logger.Error("error during update sink cache", zap.String("sinkId", sinkId), zap.Error(err))
		return err
	}

	return nil
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
		bs.logger.Info("unable to retrieve policies from agent groups ID=" + strings.Join(agentGroupIds, ", "))
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
