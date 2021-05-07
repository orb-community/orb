// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package sinks

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
	Add() (Sink, error)
}

var _ Service = (*sinkService)(nil)

type sinkService struct {
	auth mainflux.AuthServiceClient
	repo SinksRepository
}

func New(auth mainflux.AuthServiceClient, repo SinksRepository) Service {
	return &sinkService{
		auth: auth,
		repo: repo,
	}
}

func (s sinkService) Add() (Sink, error) {
	panic("implement me")
}
