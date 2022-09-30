// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/sinks"
	"reflect"
	"sync"
)

var _ sinks.SinkRepository = (*sinkRepositoryMock)(nil)

// Mock Repository
type sinkRepositoryMock struct {
	mu        sync.Mutex
	counter   uint64
	sinksMock map[string]sinks.Sink
}

func (s *sinkRepositoryMock) UpdateSinkState(ctx context.Context, sinkID string, msg string, ownerID string, state sinks.State) error {
	return nil
}

func (s *sinkRepositoryMock) RetrieveByOwnerAndId(ctx context.Context, ownerID string, key string) (sinks.Sink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.sinksMock[key]; ok {
		if s.sinksMock[key].MFOwnerID == ownerID {
			return c, nil
		} else {
			return sinks.Sink{}, sinks.ErrNotFound
		}
	}

	return sinks.Sink{}, sinks.ErrNotFound
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

func (s *sinkRepositoryMock) Update(ctx context.Context, sink sinks.Sink) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sinksMock[sink.ID]; ok {
		if s.sinksMock[sink.ID].MFOwnerID != sink.MFOwnerID {
			return errors.ErrUpdateEntity
		}
		s.sinksMock[sink.ID] = sink
		return nil
	}
	return sinks.ErrNotFound
}

func (s *sinkRepositoryMock) RetrieveAll(ctx context.Context, owner string, pm sinks.PageMetadata) (sinks.Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	first := uint64(pm.Offset) + 1
	last := first + uint64(pm.Limit)

	var sks []sinks.Sink

	id := uint64(0)
	for _, v := range s.sinksMock {
		id++
		if v.MFOwnerID == owner && id >= first && id < last {
			if reflect.DeepEqual(pm.Tags, v.Tags) || pm.Tags == nil {
				sks = append(sks, v)
			}
		}
	}

	sks = sortSinks(pm, sks)

	page := sinks.Page{
		Sinks: sks,
		PageMetadata: sinks.PageMetadata{
			Total:  s.counter,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}
	return page, nil
}

func (s *sinkRepositoryMock) RetrieveById(ctx context.Context, key string) (sinks.Sink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.sinksMock[key]; ok {
		return c, nil
	}

	return sinks.Sink{}, sinks.ErrNotFound
}

func (s *sinkRepositoryMock) Remove(ctx context.Context, owner string, key string) error {
	return nil
}
