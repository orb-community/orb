package consumer

import (
	"context"
	"github.com/ns1labs/orb/maestro/config"
	"github.com/ns1labs/orb/pkg/errors"
	"time"

	"github.com/ns1labs/orb/maestro/kubecontrol"
	maestroredis "github.com/ns1labs/orb/maestro/redis"
	"github.com/ns1labs/orb/pkg/types"
	sinkspb "github.com/ns1labs/orb/sinks/pb"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	streamSinks  = "orb.sinks"
	streamSinker = "orb.sinker"
	groupMaestro = "orb.maestro"

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

	UpdateSinkCache(ctx context.Context, data config.SinkData) (err error)
	PublishSinkStateChange(sink *sinkspb.SinkRes, status string, logsErr error, err error)

	GetActivity(sinkID string) (int64, error)
	RemoveSinkActivity(ctx context.Context, sinkId string) error

	SubscribeSinksEvents(context context.Context) error
	SubscribeSinkerEvents(context context.Context) error
}

type eventStore struct {
	kafkaUrl             string
	kubecontrol          kubecontrol.Service
	sinksClient          sinkspb.SinkServiceClient
	streamRedisClient    *redis.Client
	sinkerKeyRedisClient *redis.Client
	esconsumer           string
	logger               *zap.Logger
}

func NewEventStore(streamRedisClient, sinkerKeyRedisClient *redis.Client, kafkaUrl string, kubecontrol kubecontrol.Service, esconsumer string, sinksClient sinkspb.SinkServiceClient, logger *zap.Logger) Subscriber {
	return eventStore{
		kafkaUrl:             kafkaUrl,
		kubecontrol:          kubecontrol,
		streamRedisClient:    streamRedisClient,
		sinkerKeyRedisClient: sinkerKeyRedisClient,
		sinksClient:          sinksClient,
		esconsumer:           esconsumer,
		logger:               logger,
	}
}

// SubscribeSinkerEvents Subscribe to listen events from sinker to maestro
func (es eventStore) SubscribeSinkerEvents(ctx context.Context) error {
	err := es.streamRedisClient.XGroupCreateMkStream(ctx, streamSinker, groupMaestro, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.streamRedisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupMaestro,
			Consumer: "orb_maestro-es-consumer",
			Streams:  []string{streamSinker, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}
		for _, msg := range streams[0].Messages {
			event := msg.Values
			switch event["operation"] {
			case sinkerUpdate:
				rte := decodeSinkerStateUpdate(event)
				if rte.State == "active" {
					err = es.handleSinkerCreateCollector(ctx, rte) //sinker request create collector
				}
				es.streamRedisClient.XAck(ctx, streamSinker, groupMaestro, msg.ID)
			case <-ctx.Done():
				return errors.New("stopped listening to sinks, due to context cancellation")
			}
			if err != nil {
				es.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
				continue
			}
		}
	}
}

// SubscribeSinksEvents Subscribe to listen events from sinks to maestro
func (es eventStore) SubscribeSinksEvents(context context.Context) error {
	//listening sinker events
	err := es.streamRedisClient.XGroupCreateMkStream(context, streamSinks, groupMaestro, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.streamRedisClient.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    groupMaestro,
			Consumer: "orb_maestro-es-consumer",
			Streams:  []string{streamSinks, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}
		for _, msg := range streams[0].Messages {
			event := msg.Values
			rte, err := decodeSinksEvent(event, event["operation"].(string))
			if err != nil {
				es.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
				break
			}
			switch event["operation"] {
			case sinksCreate:
				if v, ok := rte.Config["opentelemetry"]; ok && v.(string) == "enabled" {
					err = es.handleSinksCreateCollector(context, rte) //should create collector
				}
				es.streamRedisClient.XAck(context, streamSinks, groupMaestro, msg.ID)
			case sinksUpdate:
				err = es.handleSinksUpdateCollector(context, rte) //should create collector
				es.streamRedisClient.XAck(context, streamSinks, groupMaestro, msg.ID)
			case sinksDelete:
				err = es.handleSinksDeleteCollector(context, rte) //should delete collector
				es.streamRedisClient.XAck(context, streamSinks, groupMaestro, msg.ID)
			case <-context.Done():
				return errors.New("stopped listening to sinks, due to context cancellation")
			}
			if err != nil {
				es.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
				continue
			}
		}
	}
}

// handleSinkerDeleteCollector Delete collector
func (es eventStore) handleSinkerDeleteCollector(ctx context.Context, event maestroredis.SinkerUpdateEvent) error {
	es.logger.Info("Received maestro DELETE event from sinker, sink state", zap.String("state", event.State), zap.String("sinkdID", event.SinkID), zap.String("ownerID", event.Owner))
	_, err := es.GetDeploymentEntryFromSinkId(ctx, event.SinkID)
	if err != nil {
		return err
	}
	err = es.kubecontrol.DeleteOtelCollector(ctx, event.Owner, event.SinkID)
	if err != nil {
		return err
	}
	return nil
}

// handleSinkerCreateCollector Create collector
func (es eventStore) handleSinkerCreateCollector(ctx context.Context, event maestroredis.SinkerUpdateEvent) error {
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

func decodeSinkerStateUpdate(event map[string]interface{}) maestroredis.SinkerUpdateEvent {
	val := maestroredis.SinkerUpdateEvent{
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
