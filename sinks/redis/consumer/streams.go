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

	sinkerPrefix = "sinker."
	sinkerUpdate = sinkerPrefix + "update"

	otelYamlPrefix = "otel.yaml.sinker."

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
			case sinkerUpdate:
				rte := decodeSinkerStateUpdate(event)
				err = es.handleSinkerStateUpdate(context, rte)
			}
			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, stream, group, msg.ID)
		}
	}
}

func (es eventStore) handleSinkerStateUpdate(ctx context.Context, event stateUpdateEvent) error {
	err := es.sinkService.ChangeSinkStateInternal(ctx, event.sinkID, event.msg, event.ownerID, event.state)
	if err != nil {
		return err
	}
	return nil
}

func decodeSinkerStateUpdate(event map[string]interface{}) stateUpdateEvent {
	val := stateUpdateEvent{
		ownerID:   read(event, "owner", ""),
		sinkID:    read(event, "sink_id", ""),
		msg:       read(event, "msg", ""),
		timestamp: time.Time{},
	}
	val.state.Scan(event["state"])
	return val
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}
