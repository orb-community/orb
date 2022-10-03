package consumer

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinker"
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
	"time"

	"github.com/ns1labs/orb/sinker/sinkOrchestrator"
)

const (
	stream = "orb.sinks"
	group  = "orb.sinker"

	sinksPrefix = "sinks."
	sinksUpdate = sinksPrefix + "update"

	exists = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	Subscribe(context context.Context) error
}

type eventStore struct {
	sinkerService sinker.Service
	configRepo    config.ConfigRepo
	client        *redis.Client
	esconsumer    string
	logger        *zap.Logger
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
			case sinksUpdate:
				rte, derr := decodeSinksUpdate(event)
				if derr != nil {
					err = derr
					break
				}
				err = es.handleSinksUpdate(context, rte)
			}
			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, stream, group, msg.ID)
		}
	}
}

// NewEventStore returns new event store instance.
func NewEventStore(sinkerService sinker.Service, configRepo config.ConfigRepo, client *redis.Client, esconsumer string, log *zap.Logger) Subscriber {
	return eventStore{
		sinkerService: sinkerService,
		configRepo:    configRepo,
		client:        client,
		esconsumer:    esconsumer,
		logger:        log,
	}
}

func decodeSinksUpdate(event map[string]interface{}) (updateSinkEvent, error) {
	val := updateSinkEvent{
		sinkID:    read(event, "sink_id", ""),
		owner:     read(event, "owner", ""),
		timestamp: time.Time{},
	}

	var config types.Metadata
	if err := json.Unmarshal([]byte(read(event, "config", "")), &config); err != nil {
		return updateSinkEvent{}, err
	}
	val.config = config
	return val, nil
}

func (es eventStore) handleSinksUpdate(ctx context.Context, e updateSinkEvent) error {
	data, err := json.Marshal(e.config)
	if err != nil {
		return err
	}
	var cfg config.SinkConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	if ok := es.configRepo.Exists(e.owner, e.sinkID); ok {
		sinkConfig, err := es.configRepo.Get(e.owner, e.sinkID)
		if err != nil {
			return err
		}
		sinkConfig.Url = cfg.Url
		sinkConfig.User = cfg.User
		sinkConfig.Password = cfg.Password
		if sinkConfig.OwnerID == "" {
			sinkConfig.OwnerID = e.owner
		}

		es.configRepo.Edit(sinkConfig)
		sinkOrchestrator.UpdateSink(sinkConfig)
	} else {
		cfg.SinkID = e.sinkID
		cfg.OwnerID = e.owner
		es.configRepo.Add(cfg)
	}
	return nil
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}
	return val
}
