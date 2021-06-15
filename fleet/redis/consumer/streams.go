// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ns1labs/orb/sinks/writer"
	"go.uber.org/zap"
)

const (
	stream = "mainflux.things"
	group  = "orb.prom-sink"

	thingPrefix     = "thing."
	thingRemove     = thingPrefix + "remove"
	thingDisconnect = thingPrefix + "disconnect"

	channelPrefix = "channel."
	channelUpdate = channelPrefix + "update"
	channelRemove = channelPrefix + "remove"

	exists = "BUSYGROUP Consumer Group name already exists"
)

// Subscriber represents event source for things and channels provisioning.
type Subscriber interface {
	// Subscribes to given subject and receives events.
	Subscribe(string) error
}

type eventStore struct {
	svc        writer.Service
	client     *redis.Client
	esconsumer string
	logger     *zap.Logger
}

// NewEventStore returns new event store instance.
func NewEventStore(svc writer.Service, client *redis.Client, esconsumer string, log *zap.Logger) Subscriber {
	return eventStore{
		svc:        svc,
		client:     client,
		esconsumer: esconsumer,
		logger:     log,
	}
}

func (es eventStore) Subscribe(subject string) error {
	err := es.client.XGroupCreateMkStream(stream, group, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(&redis.XReadGroupArgs{
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

			fmt.Printf("promsink consume event: %+v", event)

			var err error
			switch event["operation"] {
			case thingRemove:
				rte := decodeRemoveThing(event)
				err = es.handleRemoveThing(rte)
			}
			if err != nil {
				es.logger.Warn("Failed to handle event sourcing")
				break
			}
			es.client.XAck(stream, group, msg.ID)
		}
	}
}

func decodeRemoveThing(event map[string]interface{}) removeEvent {
	return removeEvent{
		id: read(event, "id", ""),
	}
}

func decodeUpdateChannel(event map[string]interface{}) updateChannelEvent {
	strmeta := read(event, "metadata", "{}")
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(strmeta), metadata); err != nil {
		metadata = map[string]interface{}{}
	}

	return updateChannelEvent{
		id:       read(event, "id", ""),
		name:     read(event, "name", ""),
		metadata: metadata,
	}
}

func decodeRemoveChannel(event map[string]interface{}) removeEvent {
	return removeEvent{
		id: read(event, "id", ""),
	}
}

func decodeDisconnectThing(event map[string]interface{}) disconnectEvent {
	return disconnectEvent{
		channelID: read(event, "chan_id", ""),
		thingID:   read(event, "thing_id", ""),
	}
}

func (es eventStore) handleRemoveThing(rte removeEvent) error {
	// return es.svc.RemoveConfigHandler(rte.id)
	return nil
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}
