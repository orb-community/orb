package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/orb-community/orb/pkg/errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinker"
	"github.com/orb-community/orb/sinker/config"
	"go.uber.org/zap"
)

const (
	stream = "orb.sinks"
	group  = "orb.sinker"

	sinksPrefix = "sinks."
	sinksUpdate = sinksPrefix + "update"
	sinksCreate = sinksPrefix + "create"
	sinksDelete = sinksPrefix + "remove"
	exists      = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	Subscribe(context context.Context) error
}

type eventStore struct {
	otelEnabled   bool
	sinkerService sinker.Service
	configRepo    config.ConfigRepo
	client        *redis.Client
	esconsumer    string
	logger        *zap.Logger
}

func (es eventStore) Subscribe(context context.Context) error {
	subGroup := group
	if es.otelEnabled {
		subGroup = group + ".otel"
	}
	err := es.client.XGroupCreateMkStream(context, stream, subGroup, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    subGroup,
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
			case sinksCreate:
				rte, derr := decodeSinksCreate(event)
				if derr != nil {
					err = derr
					break
				}
				err = es.handleSinksCreate(context, rte)
			case sinksUpdate:
				rte, derr := decodeSinksUpdate(event)
				if derr != nil {
					err = derr
					break
				}
				err = es.handleSinksUpdate(context, rte)
			case sinksDelete:
				rte, derr := decodeSinksRemove(event)
				if derr != nil {
					err = derr
					break
				}
				err = es.handleSinksRemove(context, rte)
			}
			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				continue
			}
			es.client.XAck(context, stream, subGroup, msg.ID)
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

func decodeSinksCreate(event map[string]interface{}) (UpdateSinkEvent, error) {
	val := UpdateSinkEvent{
		SinkID:    read(event, "sink_id", ""),
		Owner:     read(event, "owner", ""),
		Config:    readMetadata(event, "config"),
		Timestamp: time.Now(),
	}
	var metadata types.Metadata
	if err := json.Unmarshal([]byte(read(event, "config", "")), &metadata); err != nil {
		return UpdateSinkEvent{}, err
	}
	val.Config = metadata
	return val, nil
}

func decodeSinksUpdate(event map[string]interface{}) (UpdateSinkEvent, error) {
	val := UpdateSinkEvent{
		SinkID:    read(event, "sink_id", ""),
		Owner:     read(event, "owner", ""),
		Timestamp: time.Now(),
	}
	var metadata types.Metadata
	if err := json.Unmarshal([]byte(read(event, "config", "")), &metadata); err != nil {
		return UpdateSinkEvent{}, err
	}
	val.Config = metadata
	return val, nil
}

func decodeSinksRemove(event map[string]interface{}) (UpdateSinkEvent, error) {
	val := UpdateSinkEvent{
		SinkID:    read(event, "sink_id", ""),
		Owner:     read(event, "owner", ""),
		Timestamp: time.Now(),
	}
	return val, nil
}

func (es eventStore) handleSinksRemove(_ context.Context, e UpdateSinkEvent) error {
	if ok := es.configRepo.Exists(e.Owner, e.SinkID); ok {
		err := es.configRepo.Remove(e.Owner, e.SinkID)
		if err != nil {
			es.logger.Error("error during remove sinker cache entry", zap.Error(err))
			return err
		}
	} else {
		es.logger.Error("did not found any sinker cache entry for removal",
			zap.String("key", fmt.Sprintf("sinker_key-%s-%s", e.Owner, e.SinkID)))
		return errors.New("did not found any sinker cache entry for removal")
	}
	return nil
}

func (es eventStore) handleSinksUpdate(_ context.Context, e UpdateSinkEvent) error {
	var cfg config.SinkConfig
	cfg.Config = types.FromMap(e.Config)
	cfg.SinkID = e.SinkID
	cfg.OwnerID = e.Owner
	cfg.State = config.Unknown
	if ok := es.configRepo.Exists(e.Owner, e.SinkID); ok {
		sinkConfig, err := es.configRepo.Get(e.Owner, e.SinkID)
		if err != nil {
			return err
		}
		sinkConfig.Config = cfg.Config
		if sinkConfig.OwnerID == "" {
			sinkConfig.OwnerID = e.Owner
		}
		if sinkConfig.SinkID == "" {
			sinkConfig.SinkID = e.SinkID
		}
		err = es.configRepo.Edit(sinkConfig)
		if err != nil {
			return err
		}
	} else {
		err := es.configRepo.Add(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es eventStore) handleSinksCreate(_ context.Context, e UpdateSinkEvent) error {
	var cfg config.SinkConfig
	cfg.Config = types.FromMap(e.Config)
	cfg.SinkID = e.SinkID
	cfg.OwnerID = e.Owner
	cfg.State = config.Unknown
	err := es.configRepo.Add(cfg)
	if err != nil {
		return err
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

func readMetadata(event map[string]interface{}, key string) types.Metadata {
	val, ok := event[key].(types.Metadata)
	if !ok {
		return types.Metadata{}
	}

	return val
}
