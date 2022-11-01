package consumer

import (
	"context"
	"github.com/ns1labs/orb/maestro/kubecontrol"
	"github.com/ns1labs/orb/pkg/types"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	streamSinker = "orb.sinker"
	streamSinks  = "orb.sinks"
	group        = "orb.collectors"

	sinkerPrefix = "sinker."
	sinkerUpdate = sinkerPrefix + "update"

	sinksPrefix = "sinks."
	sinksUpdate = sinksPrefix + "update"
	sinksCreate = sinksPrefix + "create"
	sinksDelete = sinksPrefix + "remove"

	exists = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	CreateDeploymentEntry(ctx context.Context, sinkId, sinkUrl, sinkUsername, sinkPassword string) error
	GetDeploymentEntryFromSinkId(ctx context.Context, sinkId string) (string, error)
	SubscribeSinks(context context.Context) error
	SubscribeSinker(context context.Context) error
}

type eventStore struct {
	kubecontrol kubecontrol.Service
	sinksClient sinkspb.SinkServiceClient
	client      *redis.Client
	esconsumer  string
	logger      *zap.Logger
}

func NewEventStore(client *redis.Client, kubecontrol kubecontrol.Service, esconsumer string, sinksClient sinkspb.SinkServiceClient, logger *zap.Logger) Subscriber {
	return eventStore{
		kubecontrol: kubecontrol,
		client:      client,
		sinksClient: sinksClient,
		esconsumer:  esconsumer,
		logger:      logger,
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

			es.logger.Info("debugging event", zap.Any("sink_event", event))
			rte, err := decodeSinksEvent(event, event["operation"].(string))
			if err != nil {
				es.logger.Error("error decoding sinks event", zap.Any("operation", event["operation"]), zap.Any("sink_event", event), zap.Error(err))
				break
			}
			switch event["operation"] {
			case sinksCreate:

				if v, ok := rte.config["opentelemetry"]; ok && v.(string) == "enabled" {
					err = es.handleSinksCreateCollector(context, rte) //should create collector
				}

			case sinksUpdate:
				if v, ok := rte.config["opentelemetry"]; ok && v.(string) == "enabled" {
					err = es.handleSinksUpdateCollector(context, rte) //should create collector
				}

			case sinksDelete:
				err = es.handleSinksDeleteCollector(context, rte) //should delete collector

			}
			if err != nil {
				es.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
				break
			}
			es.client.XAck(context, streamSinks, group, msg.ID)
		}
	}
}

// handleSinkerDeleteCollector Delete collector
func (es eventStore) handleSinkerDeleteCollector(ctx context.Context, event sinkerUpdateEvent) error {
	es.logger.Info("Received maestro DELETE event from sinker, sink state=" + event.state + ", , Sink ID=" + event.sinkID + ", Owner ID=" + event.ownerID)
	deployment, err := es.GetDeploymentEntryFromSinkId(ctx, event.sinkID)
	if err != nil {
		return err
	}
	err = es.kubecontrol.DeleteOtelCollector(ctx, event.sinkID, deployment)
	if err != nil {
		return err
	}
	return nil
}

// handleSinkerCreateCollector Create collector
func (es eventStore) handleSinkerCreateCollector(ctx context.Context, event sinkerUpdateEvent) error {
	es.logger.Info("Received maestro CREATE event from sinker, sink state=" + event.state + ", Sink ID=" + event.sinkID + ", Owner ID=" + event.ownerID)
	deploymentEntry, err := es.GetDeploymentEntryFromSinkId(ctx, event.sinkID)
	if err != nil {
		return err
	}
	err = es.kubecontrol.CreateOtelCollector(ctx, event.sinkID, deploymentEntry)
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

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}

func readMetadata(event map[string]interface{}, key string) types.Metadata {
	val, ok := event[key].(types.Metadata)
	if !ok {
		return types.Metadata{}
	}

	return val
}
