// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"context"
	"github.com/orb-community/orb/sinks/authentication_type"

	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/sinks"
	"github.com/orb-community/orb/sinks/backend"
	"go.uber.org/zap"
)

const (
	streamID  = "orb.sinks"
	streamLen = 1000
)

var _ sinks.SinkService = (*sinksStreamProducer)(nil)

type sinksStreamProducer struct {
	svc    sinks.SinkService
	client *redis.Client
	logger *zap.Logger
}

// ListSinksInternal will only call following service
func (es sinksStreamProducer) ListSinksInternal(ctx context.Context, filter sinks.Filter) ([]sinks.Sink, error) {
	return es.svc.ListSinksInternal(ctx, filter)
}

func (es sinksStreamProducer) ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state sinks.State) error {
	return es.svc.ChangeSinkStateInternal(ctx, sinkID, msg, ownerID, state)
}

func (es sinksStreamProducer) ViewSinkInternal(ctx context.Context, ownerID string, key string) (sinks.Sink, error) {
	return es.svc.ViewSinkInternal(ctx, ownerID, key)
}

func (es sinksStreamProducer) CreateSink(ctx context.Context, token string, s sinks.Sink) (sink sinks.Sink, err error) {
	defer func() {
		event := createSinkEvent{
			sinkID:  sink.ID,
			owner:   sink.MFOwnerID,
			config:  sink.Config,
			backend: sink.Backend,
		}

		encode, err := event.Encode()
		if err != nil {
			es.logger.Error("error encoding object", zap.Error(err))
		}

		record := &redis.XAddArgs{
			Stream: streamID,
			MaxLen: streamLen,
			Approx: true,
			Values: encode,
		}

		err = es.client.XAdd(ctx, record).Err()
		if err != nil {
			es.logger.Error("error sending event to sinks event store", zap.Error(err))
		}

	}()

	return es.svc.CreateSink(ctx, token, s)
}

func (es sinksStreamProducer) UpdateSinkInternal(ctx context.Context, s sinks.Sink) (sink sinks.Sink, err error) {
	defer func() {
		event := updateSinkEvent{
			sinkID:  sink.ID,
			owner:   sink.MFOwnerID,
			config:  sink.Config,
			backend: sink.Backend,
		}

		encode, err := event.Encode()
		if err != nil {
			es.logger.Error("error encoding object", zap.Error(err))
		}

		record := &redis.XAddArgs{
			Stream: streamID,
			MaxLen: streamLen,
			Approx: true,
			Values: encode,
		}

		err = es.client.XAdd(ctx, record).Err()
		if err != nil {
			es.logger.Error("error sending event to sinks event store", zap.Error(err))
		}
	}()
	return es.svc.UpdateSinkInternal(ctx, s)
}

func (es sinksStreamProducer) UpdateSink(ctx context.Context, token string, s sinks.Sink) (sink sinks.Sink, err error) {
	defer func() {
		event := updateSinkEvent{
			sinkID: sink.ID,
			owner:  sink.MFOwnerID,
			config: sink.Config,
		}

		encode, err := event.Encode()
		if err != nil {
			es.logger.Error("error encoding object", zap.Error(err))
		}

		record := &redis.XAddArgs{
			Stream: streamID,
			MaxLen: streamLen,
			Approx: true,
			Values: encode,
		}

		err = es.client.XAdd(ctx, record).Err()
		if err != nil {
			es.logger.Error("error sending event to sinks event store", zap.Error(err))
		}
	}()
	return es.svc.UpdateSink(ctx, token, s)
}

func (es sinksStreamProducer) ListSinks(ctx context.Context, token string, pm sinks.PageMetadata) (sinks.Page, error) {
	return es.svc.ListSinks(ctx, token, pm)
}

func (es sinksStreamProducer) ListAuthenticationTypes(ctx context.Context, token string) ([]authentication_type.AuthenticationTypeConfig, error) {
	return es.svc.ListAuthenticationTypes(ctx, token)
}

func (es sinksStreamProducer) ViewAuthenticationType(ctx context.Context, token string, key string) (authentication_type.AuthenticationTypeConfig, error) {
	return es.svc.ViewAuthenticationType(ctx, token, key)
}

func (es sinksStreamProducer) ListBackends(ctx context.Context, token string) (_ []string, err error) {
	return es.svc.ListBackends(ctx, token)
}

func (es sinksStreamProducer) ViewBackend(ctx context.Context, token string, key string) (_ backend.Backend, err error) {
	return es.svc.ViewBackend(ctx, token, key)
}

func (es sinksStreamProducer) ViewSink(ctx context.Context, token string, key string) (_ sinks.Sink, err error) {
	return es.svc.ViewSink(ctx, token, key)
}

func (es sinksStreamProducer) GetLogger() *zap.Logger {
	return es.logger
}

func (es sinksStreamProducer) DeleteSink(ctx context.Context, token, id string) (err error) {
	sink, err := es.svc.ViewSink(ctx, token, id)
	if err != nil {
		return err
	}

	if err := es.svc.DeleteSink(ctx, token, id); err != nil {
		return err
	}

	event := deleteSinkEvent{
		sinkID:  id,
		ownerID: sink.MFOwnerID,
	}

	encode, err := event.Encode()
	if err != nil {
		es.logger.Error("error encoding object", zap.Error(err))
	}

	record := &redis.XAddArgs{
		Stream: streamID,
		MaxLen: streamLen,
		Approx: true,
		Values: encode,
	}

	err = es.client.XAdd(ctx, record).Err()
	if err != nil {
		es.logger.Error("error sending event to sinks event store", zap.Error(err))
		return err
	}
	return nil
}

func (es sinksStreamProducer) ValidateSink(ctx context.Context, token string, sink sinks.Sink) (sinks.Sink, error) {
	return es.svc.ValidateSink(ctx, token, sink)
}

// NewSinkStreamProducerMiddleware returns wrapper around sinks service that sends
// events to event store.
func NewSinkStreamProducerMiddleware(svc sinks.SinkService, client *redis.Client) sinks.SinkService {
	return sinksStreamProducer{
		svc:    svc,
		client: client,
	}
}
