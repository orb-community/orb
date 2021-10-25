// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/backend"
	"go.uber.org/zap"
)

const (
	streamID  = "orb.sinks"
	streamLen = 1000
)

var _ sinks.SinkService = (*eventStore)(nil)

type eventStore struct {
	svc    sinks.SinkService
	client *redis.Client
	logger *zap.Logger
}

func (es eventStore) ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state sinks.State) error {
	return es.svc.ChangeSinkStateInternal(ctx, sinkID, msg, ownerID, state)
}

func (es eventStore) ViewSinkInternal(ctx context.Context, ownerID string, key string) (sinks.Sink, error) {
	return es.svc.ViewSinkInternal(ctx, ownerID, key)
}

func (es eventStore) CreateSink(ctx context.Context, token string, s sinks.Sink) (sinks.Sink, error) {
	return es.svc.CreateSink(ctx, token, s)
}

func (es eventStore) UpdateSink(ctx context.Context, token string, s sinks.Sink) (err error) {
	if err := es.svc.UpdateSink(ctx, token, s); err != nil {
		return err
	}

	event := updateSinkEvent{
		sinkID: s.ID,
		owner:  s.MFOwnerID,
		config: s.Config,
	}

	encode, err := event.Encode()
	if err != nil {
		es.logger.Error("error encoding object", zap.Error(err))
		return err
	}

	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       encode,
	}

	err = es.client.XAdd(ctx, record).Err()
	if err != nil {
		es.logger.Error("error sending event to event store", zap.Error(err))
		return err
	}
	return nil
}

func (es eventStore) ListSinks(ctx context.Context, token string, pm sinks.PageMetadata) (sinks.Page, error) {
	return es.svc.ListSinks(ctx, token, pm)
}

func (es eventStore) ListBackends(ctx context.Context, token string) (_ []string, err error) {
	return es.svc.ListBackends(ctx, token)
}

func (es eventStore) ViewBackend(ctx context.Context, token string, key string) (_ backend.Backend, err error) {
	return es.svc.ViewBackend(ctx, token, key)
}

func (es eventStore) ViewSink(ctx context.Context, token string, key string) (_ sinks.Sink, err error) {
	return es.svc.ViewSink(ctx, token, key)
}

func (es eventStore) DeleteSink(ctx context.Context, token, id string) (err error) {
	if err := es.svc.DeleteSink(ctx, token, id); err != nil {
		return err
	}

	event := deleteSinkEvent{
		id: id,
	}

	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}

	err = es.client.XAdd(ctx, record).Err()
	if err != nil {
		es.logger.Error("error sending event to event store", zap.Error(err))
		return err
	}
	return nil
}

func (es eventStore) ValidateSink(ctx context.Context, token string, sink sinks.Sink) (sinks.Sink, error) {
	return es.svc.ValidateSink(ctx, token, sink)
}

func (es eventStore) SinksStatistics(ctx context.Context, token string) (sinks.SinksStatistics, error) {
	return es.svc.SinksStatistics(ctx, token)
}

// NewEventStoreMiddleware returns wrapper around sinks service that sends
// events to event store.
func NewEventStoreMiddleware(svc sinks.SinkService, client *redis.Client) sinks.SinkService {
	return eventStore{
		svc:    svc,
		client: client,
	}
}
