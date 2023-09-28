package consumer

import (
	"context"
	"github.com/go-redis/redis/v8"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/maestro/service"
	"go.uber.org/zap"
)

type SinkerActivityListener interface {
	// SubscribeSinksEvents - listen to sink_activity, sink_idle because of state management and deployments start or stop
	SubscribeSinksEvents(ctx context.Context) error
}

type sinkerActivityListenerService struct {
	logger       *zap.Logger
	redisClient  *redis.Client
	eventService service.EventService
}

func NewSinkerActivityListener(l *zap.Logger, eventService service.EventService, redisClient *redis.Client) SinkerActivityListener {
	logger := l.Named("sinker-activity-listener")
	return &sinkerActivityListenerService{
		logger:       logger,
		redisClient:  redisClient,
		eventService: eventService,
	}
}

func (s *sinkerActivityListenerService) SubscribeSinksEvents(ctx context.Context) error {
	//listening sinker events
	err := s.redisClient.XGroupCreateMkStream(ctx, maestroredis.SinksActivityStream, maestroredis.GroupMaestro, "$").Err()
	if err != nil && err.Error() != maestroredis.Exists {
		return err
	}

	err = s.redisClient.XGroupCreateMkStream(ctx, maestroredis.SinksIdleStream, maestroredis.GroupMaestro, "$").Err()
	if err != nil && err.Error() != maestroredis.Exists {
		return err
	}
	s.logger.Debug("Reading Sinker Events", zap.String("stream", maestroredis.SinksIdleStream), zap.String("stream", maestroredis.SinksActivityStream))
	for {
		streams, err := s.readStreams(ctx)
		if err != nil || len(streams) == 0 {
			continue
		}
		for _, str := range streams {
			go s.processStream(ctx, str)
		}
	}
}

// createStreamIfNotExists - create stream if not exists
func (s *sinkerActivityListenerService) createStreamIfNotExists(ctx context.Context, streamName string) error {
	err := s.redisClient.XGroupCreateMkStream(ctx, streamName, maestroredis.GroupMaestro, "$").Err()
	if err != nil && err.Error() != maestroredis.Exists {
		return err
	}
	return nil
}

// readStreams - read streams
func (s *sinkerActivityListenerService) readStreams(ctx context.Context) ([]redis.XStream, error) {
	const activityStream = "orb.sink_activity"
	const idleStream = "orb.sink_idle"
	streams, err := s.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    maestroredis.GroupMaestro,
		Consumer: "orb_maestro-es-consumer",
		Streams:  []string{activityStream, idleStream, ">"},
	}).Result()
	if err != nil {
		return nil, err
	}
	return streams, nil
}

// processStream - process stream
func (s *sinkerActivityListenerService) processStream(ctx context.Context, stream redis.XStream) {
	eventType := ""
	if stream.Stream == "orb.sink_activity" {
		eventType = "activity"
	} else if stream.Stream == "orb.sink_idle" {
		eventType = "idle"
	}
	for _, message := range stream.Messages {
		event := maestroredis.SinkerUpdateEvent{}
		event.Decode(message.Values)
		switch eventType {
		case "activity":
			s.logger.Debug("Reading message from activity stream",
				zap.String("message_id", message.ID),
				zap.String("sink_id", event.SinkID),
				zap.String("owner_id", event.OwnerID))
			err := s.eventService.HandleSinkActivity(ctx, event)
			if err != nil {
				s.logger.Error("error receiving message", zap.Error(err))
			}
		case "idle":
			s.logger.Debug("Reading message from idle stream",
				zap.String("message_id", message.ID),
				zap.String("sink_id", event.SinkID),
				zap.String("owner_id", event.OwnerID))
			err := s.eventService.HandleSinkIdle(ctx, event)
			if err != nil {
				s.logger.Error("error receiving message", zap.Error(err))
			}
		}
	}
}
