package bridgeservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/orb-community/orb/pkg/types"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"sort"
	"time"

	"github.com/go-kit/kit/metrics"
	fleetpb "github.com/orb-community/orb/fleet/pb"
	policiespb "github.com/orb-community/orb/policies/pb"
	"github.com/orb-community/orb/sinker/config"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type BridgeService interface {
	ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error)
	GetPolicyName(ctx context.Context, policyId, ownerId string) (*policiespb.PolicyRes, error)
	GetDataSetsFromAgentGroups(ctx context.Context, mfOwnerId string, agentGroupIds []string) (map[string]string, error)
	NotifyActiveSink(ctx context.Context, mfOwnerId, sinkId, state, message string) error
	GetSinkIdsFromPolicyID(ctx context.Context, mfOwnerId string, policyID string) (map[string]string, error)
	IncreamentMessageCounter(publisher, subtopic, channel, protocol string)
}

func NewBridgeService(logger *zap.Logger,
	defaultCacheExpiration time.Duration,
	sinkerCache config.ConfigRepo,
	policiesClient policiespb.PolicyServiceClient,
	sinksClient sinkspb.SinkServiceClient,
	fleetClient fleetpb.FleetServiceClient, messageInputCounter metrics.Counter) SinkerOtelBridgeService {
	return SinkerOtelBridgeService{
		defaultCacheExpiration: defaultCacheExpiration,
		inMemoryCache:          *cache.New(defaultCacheExpiration, defaultCacheExpiration*2),
		logger:                 logger,
		sinkerCache:            sinkerCache,
		policiesClient:         policiesClient,
		fleetClient:            fleetClient,
		sinksClient:            sinksClient,
		messageInputCounter:    messageInputCounter,
	}
}

type SinkerOtelBridgeService struct {
	inMemoryCache          cache.Cache
	defaultCacheExpiration time.Duration
	logger                 *zap.Logger
	sinkerCache            config.ConfigRepo
	policiesClient         policiespb.PolicyServiceClient
	fleetClient            fleetpb.FleetServiceClient
	sinksClient            sinkspb.SinkServiceClient
	messageInputCounter    metrics.Counter
}

// Implementar nova funcao
func (bs *SinkerOtelBridgeService) IncrementMessageCounter(publisher, subtopic, channel, protocol string) {
	labels := []string{
		"method", "handleMsgFromAgent",
		"agent_id", publisher,
		"subtopic", subtopic,
		"channel", channel,
		"protocol", protocol,
	}
	bs.messageInputCounter.With(labels...).Add(1)
}

func (bs *SinkerOtelBridgeService) NotifyActiveSink(ctx context.Context, mfOwnerId, sinkId, newState, message string) error {
	cfgRepo, err := bs.sinkerCache.Get(mfOwnerId, sinkId)
	if err != nil {
		bs.logger.Error("unable to retrieve the sink config", zap.Error(err))
		sinkData, _ := bs.sinksClient.RetrieveSink(ctx, &sinkspb.SinkByIDReq{
			SinkID:  sinkId,
			OwnerID: mfOwnerId,
		})
		var metadata types.Metadata
		_ = json.Unmarshal(sinkData.Config, &metadata)
		cfgRepo = config.SinkConfig{
			SinkID:  sinkId,
			OwnerID: mfOwnerId,
			Config:  metadata,
			State:   config.Active,
			Msg:     "",
		}
		err = bs.sinkerCache.DeployCollector(ctx, cfgRepo)
		if err != nil {
			bs.logger.Error("error during update sink cache", zap.String("sinkId", sinkId), zap.Error(err))
			return err
		}
	}

	// only updates sink state if status Idle or Unknown
	if cfgRepo.State == config.Idle || cfgRepo.State == config.Unknown {
		cfgRepo.LastRemoteWrite = time.Now()
		// only deploy collector if new state is "active" and current state "not active"
		if newState == "active" && cfgRepo.State != config.Active {
			err = cfgRepo.State.SetFromString(newState)
			if err != nil {
				bs.logger.Error("unable to set state", zap.String("new_state", newState), zap.Error(err))
				return err
			}
			err = bs.sinkerCache.AddActivity(mfOwnerId, sinkId)
			if err != nil {
				bs.logger.Error("error during update last remote write", zap.String("sinkId", sinkId), zap.Error(err))
				return err
			}
			err = bs.sinkerCache.DeployCollector(ctx, cfgRepo)
			if err != nil {
				bs.logger.Error("error during update sink cache", zap.String("sinkId", sinkId), zap.Error(err))
				return err
			}
			bs.logger.Info("waking up sink to active", zap.String("sinkID", sinkId), zap.String("newState", newState), zap.Any("currentState", cfgRepo.State))
		} else {
			err = bs.sinkerCache.AddActivity(mfOwnerId, sinkId)
			if err != nil {
				bs.logger.Error("error during update last remote write", zap.String("sinkId", sinkId), zap.Error(err))
				return err
			}
			bs.logger.Info("registering sink activity", zap.String("sinkID", sinkId), zap.String("newState", newState), zap.Any("currentState", cfgRepo.State))
		}
	} else {
		err = bs.sinkerCache.AddActivity(mfOwnerId, sinkId)
		if err != nil {
			bs.logger.Error("error during update last remote write", zap.String("sinkId", sinkId), zap.Error(err))
			return err
		}
		bs.logger.Info("registering sink activity", zap.String("sinkID", sinkId), zap.String("newState", newState), zap.Any("currentState", cfgRepo.State))
	}

	return nil
}

func (bs *SinkerOtelBridgeService) ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error) {
	cacheKey := fmt.Sprintf("agent-%s", channelID)
	value, found := bs.inMemoryCache.Get(cacheKey)
	if !found {
		agentPb, err := bs.fleetClient.RetrieveAgentInfoByChannelID(ctx, &fleetpb.AgentInfoByChannelIDReq{Channel: channelID})
		if err != nil {
			return nil, err
		}
		bs.inMemoryCache.Set(cacheKey, agentPb, cache.DefaultExpiration)
		return agentPb, nil
	}
	return value.(*fleetpb.AgentInfoRes), nil
}

func (bs *SinkerOtelBridgeService) GetPolicyName(ctx context.Context, policyId, ownerID string) (*policiespb.PolicyRes, error) {
	cacheKey := fmt.Sprintf("policy-%s", policyId)
	value, found := bs.inMemoryCache.Get(cacheKey)
	if !found {
		policyPb, err := bs.policiesClient.RetrievePolicy(ctx, &policiespb.PolicyByIDReq{PolicyID: policyId, OwnerID: ownerID})
		if err != nil {
			return nil, err
		}
		bs.inMemoryCache.Set(cacheKey, policyPb, cache.DefaultExpiration)
		return policyPb, nil
	}
	return value.(*policiespb.PolicyRes), nil
}

func (bs *SinkerOtelBridgeService) GetSinkIdsFromDatasetIDs(ctx context.Context, mfOwnerId string, datasetIDs []string) (map[string]string, error) {
	// Here needs to retrieve datasets
	mapSinkIdPolicy := make(map[string]string)
	sort.Strings(datasetIDs)
	for i := 0; i < len(datasetIDs); i++ {
		datasetID := datasetIDs[i]
		cacheKey := fmt.Sprintf("ds-%s-%s", mfOwnerId, datasetID)
		value, found := bs.inMemoryCache.Get(cacheKey)
		if !found {
			datasetRes, err := bs.policiesClient.RetrieveDataset(ctx, &policiespb.DatasetByIDReq{
				DatasetID: datasetID,
				OwnerID:   mfOwnerId,
			})
			if err != nil {
				bs.logger.Info("unable to retrieve datasets from policy")
				return nil, err
			}
			value = datasetRes.SinkIds
			bs.inMemoryCache.Set(cacheKey, value, cache.DefaultExpiration)
		}
		for _, sinkId := range value.([]string) {
			mapSinkIdPolicy[sinkId] = "active"
		}
	}
	return mapSinkIdPolicy, nil
}
