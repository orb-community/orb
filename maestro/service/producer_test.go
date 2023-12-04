package service

import (
	"context"
	"github.com/orb-community/orb/maestro/redis/producer"
	"go.uber.org/zap"
)

type testProducer struct {
	logger *zap.Logger
}

func NewTestProducer(logger *zap.Logger) producer.Producer {
	return &testProducer{logger: logger}
}

func (t *testProducer) PublishSinkStatus(_ context.Context, _ string, _ string, _ string, _ string) error {
	return nil
}
