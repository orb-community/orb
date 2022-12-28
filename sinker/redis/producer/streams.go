package producer

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
)

const (
	streamID = "orb.sinker"
)

var _ config.ConfigRepo = (*eventStore)(nil)

type eventStore struct {
	sinkCache config.ConfigRepo
	client    *redis.Client
	logger    *zap.Logger
}

func (e eventStore) Exists(ownerID string, sinkID string) bool {
	return e.sinkCache.Exists(ownerID, sinkID)
}

func (e eventStore) Add(config config.SinkConfig) error {
	err := e.sinkCache.Add(config)
	if err != nil {
		return err
	}

	event := SinkerUpdateEvent{
		SinkID:    config.SinkID,
		Owner:     config.OwnerID,
		State:     config.State.String(),
		Msg:       config.Msg,
		Timestamp: time.Now(),
	}
	record := &redis.XAddArgs{
		Stream: streamID,
		Values: event.Encode(),
	}
	err = e.client.XAdd(context.Background(), record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
	}
	return nil
}

func (e eventStore) Remove(ownerID string, sinkID string) error {
	err := e.sinkCache.Remove(ownerID, sinkID)
	if err != nil {
		return err
	}

	event := SinkerUpdateEvent{
		SinkID:    sinkID,
		Owner:     ownerID,
		State:     config.Idle.String(),
		Timestamp: time.Now(),
	}
	record := &redis.XAddArgs{
		Stream: streamID,
		Values: event.Encode(),
	}
	err = e.client.XAdd(context.Background(), record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
	}
	return nil
}

func (e eventStore) Get(ownerID string, sinkID string) (config.SinkConfig, error) {
	return e.sinkCache.Get(ownerID, sinkID)
}

func (e eventStore) Edit(config config.SinkConfig) error {
	err := e.sinkCache.Edit(config)
	if err != nil {
		return err
	}

	event := SinkerUpdateEvent{
		SinkID:    config.SinkID,
		Owner:     config.OwnerID,
		State:     config.State.String(),
		Msg:       config.Msg,
		Timestamp: time.Now(),
	}
	record := &redis.XAddArgs{
		Stream: streamID,
		Values: event.Encode(),
	}
	err = e.client.XAdd(context.Background(), record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
	}
	return nil
}

func (e eventStore) GetAll(ownerID string) ([]config.SinkConfig, error) {
	return e.sinkCache.GetAll(ownerID)
}

func (e eventStore) GetAllOwners() ([]string, error) {
	return e.sinkCache.GetAllOwners()
}

func NewEventStoreMiddleware(repo config.ConfigRepo, client *redis.Client) config.ConfigRepo {
	return eventStore{
		sinkCache: repo,
		client:    client,
	}
}
