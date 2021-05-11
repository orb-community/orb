// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
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

type Service interface {
	Add() (Agent, error)
}

var _ Service = (*fleetService)(nil)

type fleetService struct {
	auth         mainflux.AuthServiceClient
	agentRepo    AgentRepository
	selectorRepo SelectorRepository
	mfsdk        mfsdk.SDK
}

func New(auth mainflux.AuthServiceClient, agentRepo AgentRepository, selectorRepo SelectorRepository, mfsdk mfsdk.SDK) Service {
	return &fleetService{
		auth:         auth,
		agentRepo:    agentRepo,
		selectorRepo: selectorRepo,
		mfsdk:        mfsdk,
	}
}

func (s fleetService) Add() (Agent, error) {
	panic("implement me")
}
