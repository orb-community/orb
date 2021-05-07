// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package fleet

import (
	"errors"
	"github.com/mainflux/mainflux"
)

var (
	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrMalformedEntity indicates malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")
)

type Service interface {
	Add() (Agent, error)
}

var _ Service = (*fleetService)(nil)

type fleetService struct {
	auth mainflux.AuthServiceClient
	repo FleetRepository
}

func New(auth mainflux.AuthServiceClient, repo FleetRepository) Service {
	return &fleetService{
		auth: auth,
		repo: repo,
	}
}

func (s fleetService) Add() (Agent, error) {
	panic("implement me")
}
