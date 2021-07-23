// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/sinks"
	"strconv"
	"strings"
	"sync"
)

var _ sinks.SinkRepository = (*sinkRepositoryMock)(nil)

type sinkRepositoryMock struct {
	mu      sync.Mutex
	counter uint64
	sinksMock map[string]sinks.Sink
}

func NewSinkRepository() sinks.SinkRepository {
	return &sinkRepositoryMock{
		sinksMock: make(map[string]sinks.Sink),
	}
}

func (s *sinkRepositoryMock) Save(ctx context.Context, sink sinks.Sink) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sk := range s.sinksMock {
		if sk.Name == sink.Name {
			return "", sinks.ErrConflictSink
		}
	}

	s.counter++
	ID, _ := uuid.NewV4()
	sink.ID = ID.String()
	s.sinksMock[sink.ID] = sink

	return sink.ID, nil
}

func (s* sinkRepositoryMock) RetrieveAll(ctx context.Context, owner string, pm sinks.PageMetadata) (sinks.Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pm.Limit < 0 {
		return sinks.Page{}, nil
	}

	first := uint64(pm.Offset) + 1
	last := first + uint64(pm.Limit)

	var sks[]sinks.Sink

	prefix := fmt.Sprintf("%s", owner)
	for k, v := range s.sinksMock {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if strings.HasPrefix(k, prefix) && id >= first && id < last {
			sks = append(sks, v)
		}
	}

	page := sinks.Page{
		Sinks: sks,
		PageMetadata: sinks.PageMetadata{
			Total: s.counter,
			Offset: pm.Offset,
			Limit: pm.Limit,
		},
	}
	return page, nil
}

func (s *sinkRepositoryMock) Remove(ctx context.Context, owner string, id string) error  {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sinksMock, key(owner, id))
	return nil
}

// Since mocks will store data in map, and they need to resemble the real
// identifiers as much as possible, a key will be created as combination of
// owner and their own identifiers. This will allow searching either by
// prefix or suffix.
func key(owner string, id string) string {
	return fmt.Sprintf("%s-%s", owner, id)
}