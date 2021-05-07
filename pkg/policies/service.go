// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package policies

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
	Add() (Policy, error)
}

var _ Service = (*policiesService)(nil)

type policiesService struct {
	auth mainflux.AuthServiceClient
	repo PoliciesRepository
}

func New(auth mainflux.AuthServiceClient, repo PoliciesRepository) Service {
	return &policiesService{
		auth: auth,
		repo: repo,
	}
}

func (s policiesService) Add() (Policy, error) {
	panic("implement me")
}
