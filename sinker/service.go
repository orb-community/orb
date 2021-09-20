/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mainflux/mainflux/pkg/messaging"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"go.uber.org/zap"
)

const (
	BackendMetricsTopic = "be.*.m.*"
	MaxMsgPayloadSize   = 1024 * 100
)

var (
	ErrPayloadTooBig = errors.New("payload too big")
)

type Service interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error
}

type sinkerService struct {
	pubSub mfnats.PubSub

	esclient *redis.Client

	logger *zap.Logger
}

func (svc sinkerService) handleMetrics(thingID string, channelID string, subtopic string, payload []byte) error {
	svc.logger.Debug("received metrics",
		zap.Any("payload", payload),
		zap.String("subtopic", subtopic),
		zap.String("channel_id", channelID),
		zap.String("thing_id", thingID))
	return nil
}

func (svc sinkerService) handleMsgFromAgent(msg messaging.Message) error {

	// NOTE: we need to consider ALL input from the agent as untrusted, the same as untrusted HTTP API would be

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}

	svc.logger.Debug("received agent message",
		zap.Any("payload", payload),
		zap.String("subtopic", msg.Subtopic),
		zap.String("channel", msg.Channel),
		zap.String("protocol", msg.Protocol),
		zap.Int64("created", msg.Created),
		zap.String("publisher", msg.Publisher))

	if len(msg.Payload) > MaxMsgPayloadSize {
		return ErrPayloadTooBig
	}

	// dispatch
	switch msg.Subtopic {
	case BackendMetricsTopic:
		if err := svc.handleMetrics(msg.Publisher, msg.Channel, msg.Subtopic, msg.Payload); err != nil {
			svc.logger.Error("metrics processing failure", zap.Error(err))
			return err
		}
	default:
		svc.logger.Warn("unsupported/unhandled agent subtopic, ignoring",
			zap.String("subtopic", msg.Subtopic),
			zap.String("thing_id", msg.Publisher),
			zap.String("channel_id", msg.Channel))
	}

	return nil
}

func (svc sinkerService) Start() error {
	topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
	if err := svc.pubSub.Subscribe(topic, svc.handleMsgFromAgent); err != nil {
		return err
	}
	svc.logger.Info("started metrics consumer", zap.String("topic", topic))
	return nil
}

func (svc sinkerService) Stop() error {
	topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
	if err := svc.pubSub.Unsubscribe(topic); err != nil {
		return err
	}
	svc.logger.Info("unsubscribed from agent metrics")
	return nil
}

// New instantiates the sinker service implementation.
func New(logger *zap.Logger, pubSub mfnats.PubSub, esclient *redis.Client) Service {
	return &sinkerService{
		logger:   logger,
		pubSub:   pubSub,
		esclient: esclient,
	}
}
