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
	"github.com/mainflux/mainflux/things"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/backend"
	"strings"
	"sync"
)

var _ sinks.SinkRepository = (*sinkRepositoryMock)(nil)
var _ sinks.SinkService = (*sinkServiceMock)(nil)
var _ backend.Backend = (*backendMock)(nil)

type sinkServiceMock struct {
	Backends map[string]backendMock
}

//TODO check if it's really necessary this mock
func NewSinkServiceMock() sinks.SinkService {
	return &sinkServiceMock{
		map[string]backendMock{
			"prometheus": {
				Name:        "prometheus",
				Description: "prometheus backend",
				Config:      map[string]interface{}{"title": "Remote Host", "type": "string", "name": "remote_host"},
			},
		},
	}
}

func (s *sinkServiceMock) CreateSink(ctx context.Context, token string, sink sinks.Sink) (sinks.Sink, error) {
	return sinks.Sink{}, nil
}

func (s *sinkServiceMock) UpdateSink(ctx context.Context, token string, sink sinks.Sink) (err error) {
	return nil
}

func (s *sinkServiceMock) ListSinks(ctx context.Context, token string, pm sinks.PageMetadata) (sinks.Page, error) {
	return sinks.Page{}, nil
}

func (s *sinkServiceMock) ListBackends(ctx context.Context, token string) ([]string, error) {
	keys := make([]string, 0, len(s.Backends))
	for k := range s.Backends {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *sinkServiceMock) ViewBackend(ctx context.Context, token string, key string) (backend.Backend, error) {
	return s.Backends[key], nil
}

func (s *sinkServiceMock) ViewSink(ctx context.Context, token string, key string) (sinks.Sink, error) {
	return sinks.Sink{}, nil
}

func (s *sinkServiceMock) DeleteSink(ctx context.Context, token string, id string) error {
	return nil
}

// Backend Mock
type backendMock struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Config      types.Metadata `json:"config"`
}

func (p backendMock) Validate(config types.Metadata) error {
	return nil
}

func (p backendMock) Metadata() interface{} {
	return p.Metadata()
}

func (p backendMock) GetName() string {
	return p.Name
}

func (p backendMock) GetDescription() string {
	return p.Description
}

func (p backendMock) GetConfig() types.Metadata {
	return p.Config
}

// Mock Repository
type sinkRepositoryMock struct {
	mu        sync.Mutex
	counter   uint64
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
	s.sinksMock[key(sink.MFOwnerID, sink.ID)] = sink

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

	prefix := fmt.Sprintf("%s", owner)
	id := uint64(0)
	for _, v := range s.sinksMock {
		id++
		if strings.HasPrefix(v.MFOwnerID, prefix) && id >= first && id < last {
			sks = append(sks, v)
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

	return sinks.Sink{}, things.ErrNotFound
}

func (s *sinkRepositoryMock) Remove(ctx context.Context, owner string, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sinksMock, key(owner, id))
	return nil
}