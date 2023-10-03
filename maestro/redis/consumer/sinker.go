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

func (s *sinkerActivityListenerService) ReadSinksActivity(ctx context.Context) error {
	const activityStream = "orb.sink_activity"
	err := s.redisClient.XGroupCreateMkStream(ctx, activityStream, maestroredis.GroupMaestro, "$").Err()
	if err != nil && err.Error() != maestroredis.Exists {
		return err
	}
	go func() {
		for {
			streams, err := s.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    maestroredis.GroupMaestro,
				Consumer: "orb_maestro-es-consumer",
				Streams:  []string{activityStream, ">"},
			}).Result()
			if err != nil || len(streams) == 0 {
				if err != nil {
					s.logger.Error("error reading activity stream", zap.Error(err))
				}
				continue
			}
			for _, msg := range streams[0].Messages {
				event := maestroredis.SinkerUpdateEvent{}
				s.logger.Debug("Debug Message", zap.Any("message", msg.Values))
				event.Decode(msg.Values)
				s.logger.Debug("Reading message from activity stream",
					zap.String("message_id", msg.ID),
					zap.String("sink_id", event.SinkID),
					zap.String("owner_id", event.OwnerID))
				err := s.eventService.HandleSinkActivity(ctx, event)
				if err != nil {
					s.logger.Error("error receiving message", zap.Error(err))
					return
				}
			}
		}
	}()
	return nil
}

func (s *sinkerActivityListenerService) ReadSinksIdle(ctx context.Context) error {
	const idleStream = "orb.sink_idle"
	err := s.redisClient.XGroupCreateMkStream(ctx, idleStream, maestroredis.GroupMaestro, "$").Err()
	if err != nil && err.Error() != maestroredis.Exists {
		return err
	}
	go func() {
		for {
			streams, err := s.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    maestroredis.GroupMaestro,
				Consumer: "orb_maestro-es-consumer",
				Streams:  []string{idleStream, ">"},
			}).Result()
			if err != nil || len(streams) == 0 {
				if err != nil {
					s.logger.Error("error reading idle stream", zap.Error(err))
				}
				continue
			}
			for _, msg := range streams[0].Messages {
				event := maestroredis.SinkerUpdateEvent{}
				event.Decode(msg.Values)
				s.logger.Debug("Debug Message", zap.Any("message", msg.Values), zap.Any("event", event))
				s.logger.Debug("Reading message from idle stream",
					zap.String("message_id", msg.ID),
					zap.String("sink_id", event.SinkID),
					zap.String("owner_id", event.OwnerID))
				err := s.eventService.HandleSinkIdle(ctx, event)
				if err != nil {
					s.logger.Error("error receiving message", zap.Error(err))
					return
				}
			}
		}
	}()
	return nil
}

func (s *sinkerActivityListenerService) SubscribeSinksEvents(ctx context.Context) error {
	go func() {
		err := s.ReadSinksActivity(ctx)
		if err != nil {
			s.logger.Error("error reading activity stream", zap.Error(err))
		}
	}()
	go func() {
		err := s.ReadSinksIdle(ctx)
		if err != nil {
			s.logger.Error("error reading idle stream", zap.Error(err))
		}
	}()
	return nil
}

func (s *sinkerActivityListenerService) processActivity(ctx context.Context, stream redis.XStream) {
	for _, message := range stream.Messages {
		event := maestroredis.SinkerUpdateEvent{}
		event.Decode(message.Values)
		s.logger.Debug("Reading message from activity stream",
			zap.String("message_id", message.ID),
			zap.String("sink_id", event.SinkID),
			zap.String("owner_id", event.OwnerID))
		err := s.eventService.HandleSinkActivity(ctx, event)
		if err != nil {
			s.logger.Error("error receiving message", zap.Error(err))
		}
	}
}

func (s *sinkerActivityListenerService) processIdle(ctx context.Context, stream redis.XStream) {
	for _, message := range stream.Messages {
		event := maestroredis.SinkerUpdateEvent{}
		event.Decode(message.Values)
		s.logger.Debug("Reading message from activity stream",
			zap.String("message_id", message.ID),
			zap.String("sink_id", event.SinkID),
			zap.String("owner_id", event.OwnerID))
		err := s.eventService.HandleSinkIdle(ctx, event)
		if err != nil {
			s.logger.Error("error receiving message", zap.Error(err))
		}
	}
}
