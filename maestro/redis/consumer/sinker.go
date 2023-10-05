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
	SubscribeSinkerIdleEvents(ctx context.Context) error

	// SubscribeSinksEvents - listen to sink_activity
	SubscribeSinkerActivityEvents(ctx context.Context) error
}

type sinkerActivityListenerService struct {
	logger       *zap.Logger
	redisClient  *redis.Client
	eventService service.EventService
}

const (
	idleStream = "orb.sink_idle"
	activityStream = "orb.sink_activity"
)

func NewSinkerActivityListener(l *zap.Logger, eventService service.EventService, redisClient *redis.Client) SinkerActivityListener {
	logger := l.Named("sinker-activity-listener")
	return &sinkerActivityListenerService{
		logger:       logger,
		redisClient:  redisClient,
		eventService: eventService,
	}
}

func (s *sinkerActivityListenerService) SubscribeSinksActivity(ctx context.Context) error {
	err := s.redisClient.XGroupCreateMkStream(ctx, activityStream, maestroredis.GroupMaestro, "$").Err()
	if err != nil && err.Error() != maestroredis.Exists {
		return err
	}
	s.logger.Debug("Reading Sinker Activity Events", zap.String("stream", activityStream))
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("closing sinker_activity_listener routine")
			return nil
		default:
			streams, err := s.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    maestroredis.GroupMaestro,
				Consumer: "orb_maestro-es-consumer",
				Streams:  []string{activityStream, ">"},
				Count:    1000,
			}).Result()
			if err != nil || len(streams) == 0 {
				if err != nil {
					s.logger.Error("error reading activity stream", zap.Error(err))
				}
				continue
			}
			for _, msg := range streams[0].Messages {
				event := maestroredis.SinkerUpdateEvent{}
				event.Decode(msg.Values)
				s.logger.Debug("Reading message from activity stream",
					zap.String("message_id", msg.ID),
					zap.String("sink_id", event.SinkID),
					zap.String("owner_id", event.OwnerID))
				go func() {
					err := s.eventService.HandleSinkActivity(ctx, event)
					if err != nil {
						s.logger.Error("Failed to handle sinks event", zap.Error(err))
					} else {
						s.redisClient.XAck(ctx, activityStream, maestroredis.GroupMaestro, msg.ID)
					}
				}()
				if err != nil {
					s.logger.Error("error receiving message", zap.Error(err))
					return err
				}
			}
		}
	}
}

func (s *sinkerActivityListenerService) SubscribeSinksIdle(ctx context.Context) error {
	err := s.redisClient.XGroupCreateMkStream(ctx, idleStream, maestroredis.GroupMaestro, "$").Err()
	if err != nil && err.Error() != maestroredis.Exists {
		return err
	}
	s.logger.Debug("Reading Sinker Idle Events", zap.String("stream", idleStream))
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("closing sinker_idle_listener routine")
			return nil
		default:
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
				s.logger.Debug("Reading message from idle stream",
					zap.String("message_id", msg.ID),
					zap.String("sink_id", event.SinkID),
					zap.String("owner_id", event.OwnerID))
				go func() {
					err := s.eventService.HandleSinkIdle(ctx, event)
					if err != nil {
						s.logger.Error("Failed to handle sinks event", zap.Error(err))
					} else {
						s.redisClient.XAck(ctx, idleStream, maestroredis.GroupMaestro, msg.ID)
					}
				}()
				if err != nil {
					s.logger.Error("error receiving message", zap.Error(err))
					return err
				}
			}
		}
	}
}

func (s *sinkerActivityListenerService) SubscribeSinkerActivityEvents(ctx context.Context) error {
	err := s.SubscribeSinksActivity(ctx)
	if err != nil {
		s.logger.Error("error reading activity stream", zap.Error(err))
		return err
	}
	return nil
}

func (s *sinkerActivityListenerService) SubscribeSinkerIdleEvents(ctx context.Context) error {
	err := s.SubscribeSinksIdle(ctx)
	if err != nil {
		s.logger.Error("error reading idle stream", zap.Error(err))
		return err
	}
	return nil
}
