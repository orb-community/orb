// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package consumer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/policies"
	"go.uber.org/zap"
)

const (
	stream     = "orb.fleet"
	streamSink = "orb.sinks"
	group      = "orb.policies"

	agentGroupPrefix = "agent_group."
	agentGroupRemove = agentGroupPrefix + "remove"
	sinkPrefix       = "sinks."
	sinkRemove       = sinkPrefix + "remove"

	exists = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	SubscribeToFleet(context context.Context) error
	SubscribeToSink(context context.Context) error
}

type eventStore struct {
	policiesService policies.Service
	client          *redis.Client
	esconsumer      string
	logger          *zap.Logger
}

// NewEventStore returns new event store instance.
func NewEventStore(policiesService policies.Service, client *redis.Client, esconsumer string, log *zap.Logger) Subscriber {
	return eventStore{
		policiesService: policiesService,
		client:          client,
		esconsumer:      esconsumer,
		logger:          log,
	}
}

func (es eventStore) SubscribeToFleet(context context.Context) error {
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
			case agentGroupRemove:
				rte := decodeAgentGroupRemove(event)
				err = es.handleAgentGroupRemove(context, rte.groupID, rte.token)
			}
			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, stream, group, msg.ID)
		}
	}
}

func (es eventStore) SubscribeToSink(context context.Context) error {
	err := es.client.XGroupCreateMkStream(context, streamSink, group, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: es.esconsumer,
			Streams:  []string{streamSink, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			event := msg.Values

			var err error
			switch event["operation"] {
			case sinkRemove:
				rte := decodeSinkRemove(event)
				err = es.handleSinkRemove(context, rte.sinkID, rte.token)
			}

			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, streamSink, group, msg.ID)
		}
	}
}

func decodeAgentGroupRemove(event map[string]interface{}) removeAgentGroupEvent {
	return removeAgentGroupEvent{
		groupID: read(event, "group_id", ""),
		token:   read(event, "token", ""),
	}
}

func decodeSinkRemove(event map[string]interface{}) removeSinkEvent {
	return removeSinkEvent{
		sinkID: read(event, "sink_id", ""),
		token:  read(event, "token", ""),
	}
}

// Inactivate a Dataset after AgentGroup deletion
func (es eventStore) handleAgentGroupRemove(ctx context.Context, groupID string, token string) error {

	err := es.policiesService.InactivateDatasetByGroupID(ctx, groupID, token)
	if err != nil {
		return err
	}
	return nil
}

func (es eventStore) handleSinkRemove(ctx context.Context, sinkID string, token string) error {

	datasets, err := es.policiesService.DeleteSinkFromDataset(ctx, sinkID, token)
	if err != nil {
		return err
	}

	for _, ds := range datasets{
		if len(ds.SinkIDs) == 0{
			err = es.policiesService.InactivateDatasetByID(ctx, ds.ID, token)
			if err != nil {
				return err
			}
		}
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
