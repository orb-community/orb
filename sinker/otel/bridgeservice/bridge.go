package bridgeservice

import (
	"context"
	"github.com/ns1labs/orb/fleet/pb"
	fleetpb "github.com/ns1labs/orb/fleet/pb"
	policiespb "github.com/ns1labs/orb/policies/pb"
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
)

type BridgeService interface {
	ExtractAgent(ctx context.Context, channelID string) (*pb.AgentInfoRes, error)
	GetSinksFromDataSet(ownerID, datasetId string) (map[string]bool, error)
}

type SinkerOtelBridgeService struct {
	logger         *zap.Logger
	sinkerCache    config.ConfigRepo
	policiesClient policiespb.PolicyServiceClient
	fleetClient    fleetpb.FleetServiceClient
}

func (bs *SinkerOtelBridgeService) ExtractAgent(ctx context.Context, channelID string) (*pb.AgentInfoRes, error) {
	agentPb, err := bs.fleetClient.RetrieveAgentInfoByChannelID(ctx, &pb.AgentInfoByChannelIDReq{Channel: channelID})
	if err != nil {
		return nil, err
	}
	return agentPb, nil
}

func (bs *SinkerOtelBridgeService) GetSinksFromDataSet(ownerID, datasetId string) (map[string]bool, error) {

	_, err := bs.policiesClient.RetrieveDataset(context.Background(), &policiespb.DatasetByIDReq{
		DatasetID: datasetId,
		OwnerID:   ownerID,
	})
	if err != nil {
		bs.logger.Error("unable to retrieve dataset", zap.String("dataset_id", datasetId), zap.String("owner_id", ownerID), zap.Error(err))
		return nil, err
	}

	return nil, nil

}
