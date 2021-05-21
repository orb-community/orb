// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"github.com/go-redis/redis"
	"github.com/ns1labs/orb/policies"
)

const (
	streamID  = "orb.policies"
	streamLen = 1000
)

var _ policies.Service = (*eventStore)(nil)

type eventStore struct {
	svc    policies.Service
	client *redis.Client
}

// NewEventStoreMiddleware returns wrapper around policies service that sends
// events to event store.
func NewEventStoreMiddleware(svc policies.Service, client *redis.Client) policies.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}

func (es eventStore) Add() (policies.Policy, error) {
	return policies.Policy{}, nil
}
