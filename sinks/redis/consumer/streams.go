package consumer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinks"
	"go.uber.org/zap"
	"time"
)

const (
	stream = "orb.sinker"
	group  = "orb.sinks"

	sinkerPrefix         = "sinker."
	sinkerConnectionLost = sinkerPrefix + "connection_lost"

	exists = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	Subscribe(context context.Context) error
}

type eventStore struct {
	sinkService sinks.SinkService
	client      *redis.Client
	esconsumer  string
	logger      *zap.Logger
}

func NewEventStore(sinkService sinks.SinkService, client *redis.Client, esconsumer string, logger *zap.Logger) Subscriber {
	return eventStore{
		sinkService: sinkService,
		client:      client,
		esconsumer:  esconsumer,
		logger:      logger,
	}
}

func (es eventStore) Subscribe(context context.Context) error {
	err := es.client.XGroupCreateMkStream(context, stream, group, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: es.esconsumer,
			Streams:  []string{stream, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			event := msg.Values

			var err error
			switch event["operation"] {
			case sinkerConnectionLost:
				rte := decodeSinkerConnectionLost(event)
				err = es.handleSinkerConnectionLost(context, rte)
			}
			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, stream, group, msg.ID)
		}
	}
}

func (es eventStore) handleSinkerConnectionLost(ctx context.Context, event connectionLostEvent) error {
	err := es.sinkService.ChangeSinkStateInternal(ctx, event.sinkID, event.error, event.ownerID, sinks.Error)
	if err != nil {
		return err
	}
	return nil
}

func decodeSinkerConnectionLost(event map[string]interface{}) connectionLostEvent {
	return connectionLostEvent{
		ownerID:   read(event, "owner_id", ""),
		sinkID:    read(event, "sink_id", ""),
		timestamp: time.Time{},
	}
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}
