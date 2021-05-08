// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"github.com/go-redis/redis"
	"github.com/ns1labs/orb/pkg/sinks"
)

const (
	streamID  = "orb.sinks"
	streamLen = 1000
)

var _ sinks.Service = (*eventStore)(nil)

type eventStore struct {
	svc    sinks.Service
	client *redis.Client
}

// NewEventStoreMiddleware returns wrapper around sinks service that sends
// events to event store.
func NewEventStoreMiddleware(svc sinks.Service, client *redis.Client) sinks.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}

func (es eventStore) Add() (sinks.Sink, error) {
	return sinks.Sink{}, nil
}
