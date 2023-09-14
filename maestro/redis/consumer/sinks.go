package consumer

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/maestro/deployment"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
)

type SinksListenerController interface {
	// SubscribeSinksEvents - listen to sinks.create, sinks.update, sinks.delete to handle the deployment creation
	SubscribeSinksEvents(context context.Context) error
}

type sinksListenerService struct {
	logger            *zap.Logger
	deploymentService deployment.DeployService
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
					err = ls.handleSinksCreate(ctx, rte) //should create deployment
					if err != nil {
						ls.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
					} else {
						ls.redisClient.XAck(ctx, streamSinks, groupMaestro, msg.ID)
					}
				}()
			case sinksUpdate:
				go func() {
					err = ls.handleSinksUpdate(ctx, rte) //should create collector
					if err != nil {
						ls.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
					} else {
						ls.redisClient.XAck(ctx, streamSinks, groupMaestro, msg.ID)
					}
				}()
			case sinksDelete:
				go func() {
					err = ls.handleSinksDelete(ctx, rte) //should delete collector
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

// handleSinksUpdate logic moved to deployment.DeployService
func (ls *sinksListenerService) handleSinksUpdate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	ls.logger.Info("Received maestro UPDATE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	err := ls.deploymentService.HandleSinkUpdate(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

// handleSinksDelete logic moved to deployment.DeployService
func (ls *sinksListenerService) handleSinksDelete(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	ls.logger.Info("Received maestro DELETE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	err := ls.deploymentService.HandleSinkDelete(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

// handleSinksCreate logic moved to deployment.DeployService
func (ls *sinksListenerService) handleSinksCreate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	ls.logger.Info("Received event to CREATE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	err := ls.deploymentService.HandleSinkCreate(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
