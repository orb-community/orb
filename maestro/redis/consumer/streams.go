package consumer

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/maestro"
	"go.uber.org/zap"
)

const (
	streamID  = "orb.collectors"
	streamLen = 1000

	streamSinker = "orb.sinker"
	streamSinks  = "orb.sinks"
	group        = "orb.collectors"

	sinkerPrefix = "sinker."
	sinkerUpdate = sinkerPrefix + "update"
	sinkerCreate = sinkerPrefix + "create"
	sinkerDelete = sinkerPrefix + "remove"

	sinksPrefix = "sinks."
	sinksUpdate = sinksPrefix + "update"
	sinksCreate = sinksPrefix + "create"
	sinksDelete = sinksPrefix + "remove"

	exists = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	SubscribeSinks(context context.Context) error
	SubscribeSinker(context context.Context) error
}

type eventStore struct {
	maestroService maestro.MaestroService
	client         *redis.Client
	esconsumer     string
	logger         *zap.Logger
}

func NewEventStore(maestroService maestro.MaestroService, client *redis.Client, esconsumer string, logger *zap.Logger) Subscriber {
	return eventStore{
		maestroService: maestroService,
		client:         client,
		esconsumer:     esconsumer,
		logger:         logger,
	}
}

func (es eventStore) SubscribeSinker(context context.Context) error {
	//listening sinker events
	err := es.client.XGroupCreateMkStream(context, streamSinker, group, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: es.esconsumer,
			Streams:  []string{streamSinker, ">"},
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
				if rte.state == "idle" {
					err = es.handleSinkerDeleteCollector(context, rte) //sinker request delete collector
				} else if rte.state == "active" {
					err = es.handleSinkerCreateCollector(context, rte) //sinker request create collector
				}
			}
			if err != nil {
				es.logger.Error("Failed to handle sinker event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, streamSinker, group, msg.ID)
		}
	}
}

func (es eventStore) SubscribeSinks(context context.Context) error {
	err := es.client.XGroupCreateMkStream(context, streamSinks, group, "$").Err()
	if err != nil && err.Error() != exists {
		return nil
	}
	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: es.esconsumer,
			Streams:  []string{streamSinks, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			event := msg.Values

			var err error
			switch event["operation"] {
			case sinksCreate:
				rte := decodeSinksUpdate(event)
				err = es.handleSinksCreateCollector(context, rte) //should create collector

			case sinksUpdate:
				rte := decodeSinksUpdate(event)
				err = es.handleSinksUpdateCollector(context, rte) //should create collector

			case sinksDelete:
				rte := decodeSinksUpdate(event)
				err = es.handleSinksDeleteCollector(context, rte) //should delete collector

			}
			if err != nil {
				es.logger.Error("Failed to handle sinks event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, streamSinks, group, msg.ID)
		}
	}
}

// Delete collector
func (es eventStore) handleSinksDeleteCollector(ctx context.Context, event sinksUpdateEvent) error {
	es.logger.Info("Received maestro DELETE event from sinks ID=" + event.sinkID + ", Owner ID=" + event.ownerID)
	err := es.maestroService.DeleteOtelCollector(ctx, event.sinkID, event.config, event.ownerID)
	if err != nil {
		return err
	}
	return nil
}

// handleSinksCreateCollector will create Deployment Entry in Redis
func (es eventStore) handleSinksCreateCollector(ctx context.Context, event sinksUpdateEvent) error {
	es.logger.Info("Received maestro CREATE event from sinks ID=" + event.sinkID + ", Owner ID=" + event.ownerID)

	es.client.HSet(ctx, event.sinkID)

	return nil
}

// handleSinksUpdateCollector will update Deployment Entry in Redis
func (es eventStore) handleSinksUpdateCollector(ctx context.Context, event sinksUpdateEvent) error {
	es.logger.Info("Received maestro UPDATE event from sinks ID=" + event.sinkID + ", Owner ID=" + event.ownerID)

	return nil
}

// Delete collector
func (es eventStore) handleSinkerDeleteCollector(ctx context.Context, event sinkerUpdateEvent) error {
	es.logger.Info("Received maestro DELETE event from sinker, sink state=" + event.state + ", , Sink ID=" + event.sinkID + ", Owner ID=" + event.ownerID)
	err := es.maestroService.DeleteOtelCollector(ctx, event.sinkID, event.state, event.ownerID)
	if err != nil {
		return err
	}
	return nil
}

// Create collector
func (es eventStore) handleSinkerCreateCollector(ctx context.Context, event sinkerUpdateEvent) error {
	es.logger.Info("Received maestro CREATE event from sinker, sink state=" + event.state + ", Sink ID=" + event.sinkID + ", Owner ID=" + event.ownerID)

	err := es.maestroService.CreateOtelCollector(ctx, event.sinkID, event.state, event.ownerID)
	if err != nil {
		return err
	}
	return nil
}

func decodeSinkerStateUpdate(event map[string]interface{}) sinkerUpdateEvent {
	val := sinkerUpdateEvent{
		ownerID:   read(event, "owner", ""),
		sinkID:    read(event, "sink_id", ""),
		state:     read(event, "state", ""),
		timestamp: time.Time{},
	}
	return val
}

func decodeSinksUpdate(event map[string]interface{}) sinksUpdateEvent {
	val := sinksUpdateEvent{
		ownerID:   read(event, "owner", ""),
		sinkID:    read(event, "sink_id", ""),
		config:    read(event, "config", ""),
		timestamp: time.Time{},
	}
	return val
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}
