package mocks

import (
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/nats"
)

var _ nats.PubSub = (*pubSubMock)(nil)

type pubSubMock struct {

}

func NewPubSub(url, queue string, logger logger.Logger) (nats.PubSub, error) {
	return &pubSubMock{}, nil
}

func (p pubSubMock) Publish(topic string, msg messaging.Message) error {
	return nil
}

func (p pubSubMock) Subscribe(topic string, handler messaging.MessageHandler) error {
	return nil
}

func (p pubSubMock) Unsubscribe(topic string) error {
	return nil
}

func (p pubSubMock) Close() {
}
