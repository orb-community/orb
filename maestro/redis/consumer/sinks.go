package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/maestro/config"
	"github.com/orb-community/orb/maestro/deployment"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/pkg/types"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
	"time"
)

type SinksListenerController interface {
	// SubscribeSinksEvents - listen to sinks.create, sinks.update, sinks.delete to handle the deployment creation
	SubscribeSinksEvents(context context.Context) error
}

type sinksListenerService struct {
	logger            *zap.Logger
	deploymentService deployment.Service
	redisClient       *redis.Client
	sinksClient       sinkspb.SinkServiceClient
}

// SubscribeSinksEvents Subscribe to listen events from sinks to maestro
func (ls *sinksListenerService) SubscribeSinksEvents(ctx context.Context) error {
	//listening sinker events
	err := ls.redisClient.XGroupCreateMkStream(ctx, streamSinks, groupMaestro, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := ls.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupMaestro,
			Consumer: "orb_maestro-es-consumer",
			Streams:  []string{streamSinks, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}
		for _, msg := range streams[0].Messages {
			event := msg.Values
			rte, err := decodeSinksEvent(event, event["operation"].(string))
			if err != nil {
				ls.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
				break
			}
			ls.logger.Info("received message in sinks event bus", zap.Any("operation", event["operation"]))
			switch event["operation"] {
			case sinksCreate:
				go func() {
					err = ls.handleSinksCreateCollector(ctx, rte) //should create collector
					if err != nil {
						ls.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
					} else {
						ls.redisClient.XAck(ctx, streamSinks, groupMaestro, msg.ID)
					}
				}()
			case sinksUpdate:
				go func() {
					err = ls.handleSinksUpdateCollector(ctx, rte) //should create collector
					if err != nil {
						ls.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
					} else {
						ls.redisClient.XAck(ctx, streamSinks, groupMaestro, msg.ID)
					}
				}()
			case sinksDelete:
				go func() {
					err = ls.handleSinksDeleteCollector(ctx, rte) //should delete collector
					if err != nil {
						ls.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
					} else {
						ls.redisClient.XAck(ctx, streamSinks, groupMaestro, msg.ID)
					}
				}()
			case <-ctx.Done():
				return errors.New("stopped listening to sinks, due to context cancellation")
			}
		}
	}
}

// handleSinksUpdateCollector This will move to DeploymentService
func (ls *sinksListenerService) handleSinksUpdateCollector(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	var metadata types.Metadata
	if err := json.Unmarshal(event.Config, &metadata); err != nil {
		return err
	}
	data := config.SinkData{
		SinkID:  sinkData.Id,
		OwnerID: sinkData.OwnerID,
		Backend: sinkData.Backend,
		Config:  metadata,
	}
	_ = data.State.SetFromString(sinkData.State)

	deploy, err := config.BuildDeploymentJson(es.kafkaUrl, data)

	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", event.SinkID), zap.Error(err))
		return err
	}
	err = es.sinkerKeyRedisClient.HSet(ctx, deploymentKey, event.SinkID, deploy).Err()
	if err != nil {
		es.logger.Error("error trying to update deployment json for sink ID", zap.String("sinkId", event.SinkID), zap.Error(err))
		return err
	}
	err = es.kubecontrol.UpdateOtelCollector(ctx, event.Owner, event.SinkID, deploy)
	if err != nil {
		return err
	}
	return nil
}

// handleSinksDeleteCollector will delete Deployment Entry and force delete otel collector
func (ls *sinksListenerService) handleSinksDeleteCollector(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	ls.logger.Info("Received maestro DELETE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))

	deploymentEntry, _, err := ls.deploymentService.GetDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		ls.logger.Warn("did not find collector entry for sink", zap.String("sink-id", event.SinkID))
		return err
	}
	if deploymentEntry.LastCollectorDeployTime != nil || deploymentEntry.LastCollectorDeployTime.Before(time.Now()) {
		if deploymentEntry.LastCollectorStopTime != nil || deploymentEntry.LastCollectorStopTime.Before(time.Now()) {
			ls.logger.Warn("collector is not running, skipping")
		} else {

		}
	}
	err = ls.deploymentService.RemoveDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		return err
	}

	return nil
}

// handleSinksCreateCollector will create Deployment Entry in Redis
func (ls *sinksListenerService) handleSinksCreateCollector(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	ls.logger.Info("Received event to Create DeploymentEntry from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	sinkData, err := ls.sinksClient.RetrieveSink(ctx, &sinkspb.SinkByIDReq{
		SinkID:  event.SinkID,
		OwnerID: event.Owner,
	})
	if err != nil || (sinkData != nil && sinkData.Config == nil) {
		ls.logger.Error("could not fetch info for sink", zap.String("sink-id", event.SinkID), zap.Error(err))
		return err
	}
	var metadata types.Metadata
	if err := json.Unmarshal(sinkData.Config, &metadata); err != nil {
		return err
	}
	data := config.SinkData{
		SinkID:  sinkData.Id,
		OwnerID: sinkData.OwnerID,
		Backend: sinkData.Backend,
		Config:  metadata,
	}
	deploymentEntry := deployment.NewDeployment(sinkData.OwnerID, sinkData.Id, metadata)
	err2 := ls.deploymentService.CreateDeployment(ctx, data)
	if err2 != nil {
		return err2
	}

	return nil
}
