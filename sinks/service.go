// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"context"
	"github.com/mainflux/mainflux"
)

type Service interface {
	// CreateAgent creates new agent
	CreateAgent(ctx context.Context, token string, s Sink) (Sink, error)
}

var _ Service = (*sinkService)(nil)

type sinkService struct {
	auth mainflux.AuthServiceClient
	repo SinkRepository
}

func New(auth mainflux.AuthServiceClient, repo SinkRepository) Service {
	return &sinkService{
		auth: auth,
		repo: repo,
	}
}

func (s sinkService) CreateAgent(ctx context.Context, token string, sink Sink) (Sink, error) {
	panic("implement me")
}
