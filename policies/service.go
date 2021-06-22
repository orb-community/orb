// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"context"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
	"time"
)

var (
	ErrCreatePolicy = errors.New("failed to create policy")
)

type Service interface {
	// CreatePolicy creates new data sink
	CreatePolicy(ctx context.Context, token string, p Policy) (Policy, error)
}

var _ Service = (*policiesService)(nil)

type policiesService struct {
	auth mainflux.AuthServiceClient
	repo Repository
}

func (s policiesService) identify(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return res.GetId(), nil
}

func (s policiesService) CreatePolicy(ctx context.Context, token string, p Policy) (Policy, error) {

	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Policy{}, err
	}

	p.MFOwnerID = mfOwnerID

	err = s.repo.Save(ctx, p)
	if err != nil {
		return Policy{}, errors.Wrap(ErrCreatePolicy, err)
	}
	return p, nil
}

func New(auth mainflux.AuthServiceClient, repo Repository) Service {
	return &policiesService{
		auth: auth,
		repo: repo,
	}
}
