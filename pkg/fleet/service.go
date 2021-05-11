// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"errors"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
)

var (
	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrMalformedEntity indicates malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = errors.New("failed to scan metadata")
)

// A flat kv pair object
type Tags map[string]interface{}

// Maybe a full object hierarchy
type Metadata map[string]interface{}

type Service interface {
	AgentService
	SelectorService
}

var _ Service = (*fleetService)(nil)

type fleetService struct {
	auth         mainflux.AuthServiceClient
	agentRepo    AgentRepository
	selectorRepo SelectorRepository
	mfsdk        mfsdk.SDK
}

func (f fleetService) CreateSelector(ctx context.Context, token string, s Selector) (Selector, error) {
	panic("implement me")
}

func (f fleetService) CreateAgent(ctx context.Context, token string, a Agent) (Agent, error) {
	panic("implement me")
}

func NewFleetService(auth mainflux.AuthServiceClient, agentRepo AgentRepository, selectorRepo SelectorRepository, mfsdk mfsdk.SDK) Service {
	return &fleetService{
		auth:         auth,
		agentRepo:    agentRepo,
		selectorRepo: selectorRepo,
		mfsdk:        mfsdk,
	}
}
