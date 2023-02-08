package consumer

import (
	"context"
	"github.com/ns1labs/orb/maestro/config"
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

	UpdateSinkCache(ctx context.Context, data config.SinkData) (err error)
	PublishSinkStateChange(sink *sinkspb.SinkRes, status string, logsErr error, err error)

	GetActivity(sinkID string) (int64, error)
	RemoveSinkActivity(ctx context.Context, sinkId string) error

	Subscribe(context context.Context) error
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

// to listen events from sinker to maestro
func (es eventStore) Subscribe(context context.Context) error {
	//listening sinker events
	err := es.streamRedisClient.XGroupCreateMkStream(context, streamMaestro, groupSinker, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.streamRedisClient.XReadGroup(context, &redis.XReadGroupArgs{
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
			switch event["operation"] {
			case sinkerUpdate:
				rte := decodeSinkerStateUpdate(event)
				if rte.State == "active" {
					err = es.handleSinkerCreateCollector(context, rte) //sinker request create collector
				}
			case sinksCreate:
				rte, err := decodeSinksEvent(event, event["operation"].(string))
				if err != nil {
					es.logger.Error("error decoding sinks event", zap.Any("operation", event["operation"]), zap.Any("sink_event", event), zap.Error(err))
					break
				}
				if v, ok := rte.Config["opentelemetry"]; ok && v.(string) == "enabled" {
					err = es.handleSinksCreateCollector(context, rte) //should create collector
				}
			case sinksUpdate:
				rte, err := decodeSinksEvent(event, event["operation"].(string))
				if err != nil {
					es.logger.Error("error decoding sinks event", zap.Any("operation", event["operation"]), zap.Any("sink_event", event), zap.Error(err))
					break
				}
				err = es.handleSinksUpdateCollector(context, rte) //should create collector

			case sinksDelete:
				rte, err := decodeSinksEvent(event, event["operation"].(string))
				if err != nil {
					es.logger.Error("error decoding sinks event", zap.Any("operation", event["operation"]), zap.Any("sink_event", event), zap.Error(err))
					break
				}
				err = es.handleSinksDeleteCollector(context, rte) //should delete collector
			}
			if err != nil {
				es.logger.Error("Failed to handle sinks event", zap.Any("operation", event["operation"]), zap.Error(err))
				break
			}
			es.streamRedisClient.XAck(context, streamMaestro, groupSinks, msg.ID)
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
