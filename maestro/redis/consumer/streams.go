package consumer

import (
	"context"
	"time"

	"github.com/ns1labs/orb/maestro/kubecontrol"
	redis2 "github.com/ns1labs/orb/maestro/redis"
	"github.com/ns1labs/orb/pkg/types"
	sinkspb "github.com/ns1labs/orb/sinks/pb"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	streamMaestro = "orb.maestro"
	groupSinker   = "orb.sinkerCollectors"
	groupSinks    = "orb.sinksCollectors"

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
	kafkaUrl    string
	kubecontrol kubecontrol.Service
	sinksClient sinkspb.SinkServiceClient
	client      *redis.Client
	esconsumer  string
	logger      *zap.Logger
}

func NewEventStore(client *redis.Client, kafkaUrl string, kubecontrol kubecontrol.Service, esconsumer string, sinksClient sinkspb.SinkServiceClient, logger *zap.Logger) Subscriber {
	return eventStore{
		kafkaUrl:    kafkaUrl,
		kubecontrol: kubecontrol,
		client:      client,
		sinksClient: sinksClient,
		esconsumer:  esconsumer,
		logger:      logger,
	}
}

// to listen events from sinker to maestro
func (es eventStore) SubscribeSinker(context context.Context) error {
	//listening sinker events
	err := es.client.XGroupCreateMkStream(context, streamMaestro, groupSinker, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    groupSinker,
			Consumer: "sinker.maestro",
			Streams:  []string{streamMaestro, ">"},
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
				if rte.State == "active" {
					err = es.handleSinkerCreateCollector(context, rte) //sinker request create collector
				}
			}
			if err != nil {
				es.logger.Error("Failed to handle sinker event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, streamMaestro, groupSinker, msg.ID)
		}
	}
}

func (es eventStore) SubscribeSinks(context context.Context) error {
	err := es.client.XGroupCreateMkStream(context, streamMaestro, groupSinks, "$").Err()
	if err != nil && err.Error() != exists {
		return nil
	}
	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    groupSinks,
			Consumer: "sinks.maestro",
			Streams:  []string{streamMaestro, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			event := msg.Values

			rte, err := decodeSinksEvent(event, event["operation"].(string))
			if err != nil {
				es.logger.Error("error decoding sinks event", zap.Any("operation", event["operation"]), zap.Any("sink_event", event), zap.Error(err))
				break
			}
			es.logger.Info("Decoded sinks event", zap.Any("event", event))
			switch event["operation"] {
			case sinksCreate:
				es.logger.Info("Received Sinks create event from sinks, first step", zap.Any("event", event))
				if v, ok := rte.Config["opentelemetry"]; ok && v.(string) == "enabled" {
					es.logger.Info("Received Sinks create event from sinks, second step", zap.Any("event", event))
					err = es.handleSinksCreateCollector(context, rte) //should create collector
				}

			case sinksUpdate:
				if v, ok := rte.Config["opentelemetry"]; ok && v.(string) == "enabled" {
					err = es.handleSinksUpdateCollector(context, rte) //should create collector
				}

			case sinksDelete:
				if v, ok := rte.Config["opentelemetry"]; ok && v.(string) == "enabled" {
					err = es.handleSinksDeleteCollector(context, rte) //should delete collector
				}

			}
			if err != nil {
				es.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
				break
			}
			es.client.XAck(context, streamMaestro, groupSinks, msg.ID)
		}
	}
}

// handleSinkerDeleteCollector Delete collector
func (es eventStore) handleSinkerDeleteCollector(ctx context.Context, event redis2.SinkerUpdateEvent) error {
	es.logger.Info("Received maestro DELETE event from sinker, sink state", zap.String("state", event.State), zap.String("sinkdID", event.SinkID), zap.String("ownerID", event.Owner))
	deployment, err := es.GetDeploymentEntryFromSinkId(ctx, event.SinkID)
	if err != nil {
		return err
	}
	err = es.kubecontrol.DeleteOtelCollector(ctx, event.Owner, event.SinkID, deployment)
	if err != nil {
		return err
	}
	return nil
}

// handleSinkerCreateCollector Create collector
func (es eventStore) handleSinkerCreateCollector(ctx context.Context, event redis2.SinkerUpdateEvent) error {
	es.logger.Info("Received maestro CREATE event from sinker, sink state", zap.String("state", event.State), zap.String("sinkdID", event.SinkID), zap.String("ownerID", event.Owner))
	deploymentEntry, err := es.GetDeploymentEntryFromSinkId(ctx, event.SinkID)
	if err != nil {
		es.logger.Error("could not find deployment entry from sink-id", zap.String("sinkID", event.SinkID), zap.Error(err))
		return err
	}
	err = es.kubecontrol.CreateOtelCollector(ctx, event.Owner, event.SinkID, deploymentEntry)
	if err != nil {
		es.logger.Error("could not find deployment entry from sink-id", zap.String("sinkID", event.SinkID), zap.Error(err))
		return err
	}
	return nil
}

func decodeSinkerStateUpdate(event map[string]interface{}) redis2.SinkerUpdateEvent {
	val := redis2.SinkerUpdateEvent{
		Owner:     read(event, "owner", ""),
		SinkID:    read(event, "sink_id", ""),
		State:     read(event, "state", ""),
		Timestamp: time.Time{},
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
