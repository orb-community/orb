package consumer

import (
	"context"
	"github.com/go-redis/redis/v8"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/maestro/service"
	redis2 "github.com/orb-community/orb/sinks/redis"
	"go.uber.org/zap"
)

type SinkerActivityListener interface {
	// SubscribeSinksEvents - listen to sink_activity, sink_idle because of state management and deployments start or stop
	SubscribeSinksEvents(ctx context.Context) error
	// ListenToActivity - go routine to handle the sink activity stream
	ListenToActivity(ctx context.Context) error
	// ListenToIdle - go routine to handle the sink idle stream
	ListenToIdle(ctx context.Context) error
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
	err := s.redisClient.XGroupCreateMkStream(ctx, redis2.StreamSinks, redis2.GroupMaestro, "$").Err()
	if err != nil && err.Error() != redis2.Exists {
		return err
	}

	for {
		streams, err := s.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    redis2.GroupMaestro,
			Consumer: "orb_maestro-es-consumer",
			Streams:  []string{"orb.sink_activity", "orb.sink_idle", ">"},
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}
		for _, stream := range streams {
			go func() {
				if stream.Stream == "orb.sink_activity" {
					for _, message := range stream.Messages {
						event := maestroredis.SinkerUpdateEvent{}
						event.Decode(message.Values)
						err := s.eventService.HandleSinkActivity(ctx, message)
						if err != nil {
							s.logger.Error("error receiving message", zap.Error(err))
						}
					}
				} else if stream.Stream == "orb.sink_idle" {
					for _, message := range stream.Messages {
						err := s.ReceiveIdleMessage(ctx, message)
						if err != nil {
							s.logger.Error("error receiving message", zap.Error(err))
						}
					}
				}
			}()
		}

	}
}
