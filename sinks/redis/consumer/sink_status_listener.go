package consumer

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/sinks"
	"go.uber.org/zap"
)

const exists2 = "BUSYGROUP Consumer Group name already exists"

type SinkStatusListener interface {
	SubscribeToMaestroSinkStatus(ctx context.Context) error
	ReceiveMessage(ctx context.Context, message redis.XMessage) error
}

type sinkUpdateStatusEvent struct {
	ownerId      string
	sinkId       string
	status       string
	errorMessage string
}

type sinkStatusListener struct {
	logger       *zap.Logger
	streamClient *redis.Client
	sinkService  sinks.SinkService
}

func NewSinkStatusListener(l *zap.Logger, streamClient *redis.Client, sinkService sinks.SinkService) SinkStatusListener {
	logger := l.Named("sink_status_listener")
	return &sinkStatusListener{
		logger:       logger,
		streamClient: streamClient,
		sinkService:  sinkService,
	}
}

func (s *sinkStatusListener) SubscribeToMaestroSinkStatus(ctx context.Context) error {
	// First will create consumer group
	groupName := "orb.sinks"
	streamName := "orb.maestro.sink_status"
	consumerName := "sinks_consumer"
	err := s.streamClient.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err()
	if err != nil && err.Error() != exists2 {
		s.logger.Error("failed to create group", zap.Error(err))
		return err
	}
	go func(rLogger *zap.Logger) {
		for {
			select {
			case <-ctx.Done():
				rLogger.Info("closing sink_status_listener routine")
				return
			default:
				streams, err := s.streamClient.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    groupName,
					Consumer: consumerName,
					Streams:  []string{streamName, ">"},
					Count:    1000,
				}).Result()
				if err != nil || len(streams) == 0 {
					continue
				}
				for _, msg := range streams[0].Messages {
					err = s.ReceiveMessage(ctx, msg)
					if err != nil {
						rLogger.Error("failed to process message", zap.Error(err))
					}
				}
			}
		}
	}(s.logger.Named("goroutine_sink_status_listener"))
	return nil
}

func (s *sinkStatusListener) ReceiveMessage(ctx context.Context, message redis.XMessage) error {
	logger := s.logger.Named(fmt.Sprintf("sink_status_msg:%s", message.ID))
	go func(ctx context.Context, logger *zap.Logger, message redis.XMessage) {
		event := s.decodeMessage(message.Values)
		gotSink, err := s.sinkService.ViewSinkInternal(ctx, event.ownerId, event.sinkId)
		if err != nil {
			logger.Error("failed to get sink for sink_id from message", zap.String("owner_id", event.ownerId),
				zap.String("sink_id", event.sinkId), zap.Error(err))
			return
		}
		newState := sinks.NewStateFromString(event.status)
		if newState == sinks.Error || newState == sinks.ProvisioningError || newState == sinks.Warning {
			gotSink.Error = event.errorMessage
		}
		_, err = s.sinkService.UpdateSinkInternal(ctx, gotSink)
		if err != nil {
			logger.Error("failed to update sink", zap.String("owner_id", event.ownerId),
				zap.String("sink_id", event.sinkId), zap.Error(err))
			return
		}
	}(ctx, logger, message)
	return nil
}

// func (es eventStore) decodeSinkerStateUpdate(event map[string]interface{}) *sinks.SinkerStateUpdate {
func (s *sinkStatusListener) decodeMessage(content map[string]interface{}) sinkUpdateStatusEvent {
	return sinkUpdateStatusEvent{
		ownerId:      content["owner_id"].(string),
		sinkId:       content["sink_id"].(string),
		status:       content["status"].(string),
		errorMessage: content["error_message"].(string),
	}
}
