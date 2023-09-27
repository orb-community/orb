package consumer

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/maestro/service"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	redis2 "github.com/orb-community/orb/sinks/redis"
	"go.uber.org/zap"
)

type SinksListener interface {
	// SubscribeSinksEvents - listen to sinks.create, sinks.update, sinks.delete to handle the deployment creation
	SubscribeSinksEvents(context context.Context) error
}

type sinksListenerService struct {
	logger            *zap.Logger
	deploymentService service.EventService
	redisClient       *redis.Client
	sinksClient       sinkspb.SinkServiceClient
}

func NewSinksListenerController(l *zap.Logger, eventService service.EventService, redisClient *redis.Client,
	sinksClient sinkspb.SinkServiceClient) SinksListener {
	logger := l.Named("sinks_listener")
	return &sinksListenerService{
		logger:            logger,
		deploymentService: eventService,
		redisClient:       redisClient,
		sinksClient:       sinksClient,
	}
}

// SubscribeSinksEvents Subscribe to listen events from sinks to maestro
func (ls *sinksListenerService) SubscribeSinksEvents(ctx context.Context) error {
	//listening sinker events
	err := ls.redisClient.XGroupCreateMkStream(ctx, redis2.StreamSinks, redis2.GroupMaestro, "$").Err()
	if err != nil && err.Error() != redis2.Exists {
		return err
	}
	ls.logger.Info("Reading Sinks Events", zap.String("stream", redis2.StreamSinks))
	for {
		streams, err := ls.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    redis2.GroupMaestro,
			Consumer: "orb_maestro-es-consumer",
			Streams:  []string{redis2.StreamSinks, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}
		for _, msg := range streams[0].Messages {
			err := ls.ReceiveMessage(ctx, msg)
			if err != nil {
				return err
			}
		}
	}
}

func (ls *sinksListenerService) ReceiveMessage(ctx context.Context, msg redis.XMessage) error {
	logger := ls.logger.Named("sinks_listener:" + msg.ID)
	event := msg.Values
	rte, err := redis2.DecodeSinksEvent(event, event["operation"].(string))
	if err != nil {
		logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
		return err
	}
	logger.Info("received message in sinks event bus", zap.Any("operation", event["operation"]))
	switch event["operation"] {
	case redis2.SinkCreate:
		go func() {
			err = ls.handleSinksCreate(ctx, rte) //should create deployment
			if err != nil {
				logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
			} else {
				ls.redisClient.XAck(ctx, redis2.StreamSinks, redis2.GroupMaestro, msg.ID)
			}
		}()
	case redis2.SinkUpdate:
		go func() {
			err = ls.handleSinksUpdate(ctx, rte) //should create collector
			if err != nil {
				logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
			} else {
				ls.redisClient.XAck(ctx, redis2.StreamSinks, redis2.GroupMaestro, msg.ID)
			}
		}()
	case redis2.SinkDelete:
		go func() {
			err = ls.handleSinksDelete(ctx, rte) //should delete collector
			if err != nil {
				logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
			} else {
				ls.redisClient.XAck(ctx, redis2.StreamSinks, redis2.GroupMaestro, msg.ID)
			}
		}()
	case <-ctx.Done():
		return errors.New("stopped listening to sinks, due to context cancellation")
	}
	return nil
}

// handleSinksUpdate logic moved to deployment.EventService
func (ls *sinksListenerService) handleSinksUpdate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	ls.logger.Info("Received sinks UPDATE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	err := ls.deploymentService.HandleSinkUpdate(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

// handleSinksDelete logic moved to deployment.EventService
func (ls *sinksListenerService) handleSinksDelete(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	ls.logger.Info("Received sinks DELETE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	err := ls.deploymentService.HandleSinkDelete(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

// handleSinksCreate logic moved to deployment.EventService
func (ls *sinksListenerService) handleSinksCreate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	ls.logger.Info("Received sinks to CREATE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	err := ls.deploymentService.HandleSinkCreate(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
