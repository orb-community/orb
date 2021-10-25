package producer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
	"time"
)

const (
	streamID  = "orb.sinker"
	streamLen = 1000
)

var _ config.ConfigRepo = (*eventStore)(nil)

type eventStore struct {
	repo   config.ConfigRepo
	client *redis.Client
	logger *zap.Logger
}

func (e eventStore) Exists(sinkID string) bool {
	return e.repo.Exists(sinkID)
}

func (e eventStore) Add(config config.SinkConfig) error {
	err := e.repo.Add(config)
	if err != nil {
		return err
	}

	event := ChangeSinkerStateEvent{
		SinkID:    config.SinkID,
		Owner:     config.OwnerID,
		State:     config.State.String(),
		Msg:       config.Msg,
		Timestamp: time.Now(),
	}
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}
	err = e.client.XAdd(context.Background(), record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
	}
	return nil
}

func (e eventStore) Get(sinkID string) (config.SinkConfig, error) {
	return e.repo.Get(sinkID)
}

func (e eventStore) Edit(config config.SinkConfig) error {
	err := e.repo.Edit(config)
	if err != nil {
		return err
	}

	event := ChangeSinkerStateEvent{
		SinkID:    config.SinkID,
		Owner:     config.OwnerID,
		State:     config.State.String(),
		Msg:       config.Msg,
		Timestamp: time.Now(),
	}
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}
	err = e.client.XAdd(context.Background(), record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
	}
	return nil
}

func (e eventStore) GetAll() ([]config.SinkConfig, error) {
	return e.repo.GetAll()
}

func NewEventStoreMiddleware(repo config.ConfigRepo, client *redis.Client) config.ConfigRepo {
	return eventStore{
		repo:   repo,
		client: client,
	}
}
