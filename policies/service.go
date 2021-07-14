// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"context"
	"fmt"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/policies/backend"
	"github.com/ns1labs/orb/policies/backend/orb"
	"github.com/ns1labs/orb/policies/backend/pktvisor"
	"time"
)

var (
	ErrCreatePolicy    = errors.New("failed to create policy")
	ErrCreateDataset   = errors.New("failed to create dataset")
	ErrMalformedEntity = errors.New("malformed entity")
)

type Service interface {
	// CreatePolicy creates new agent Policy
	CreatePolicy(ctx context.Context, token string, p Policy, format string, policyData string) (Policy, error)

	// RetrievePolicyByIDInternal gRPC version of retrieving policy by id with no token
	RetrievePolicyByIDInternal(ctx context.Context, policyID string, ownerID string) (Policy, error)

	// RetrievePoliciesByGroupIDInternal gRPC version of retrieving list of policies belonging to specified agent group with no token
	RetrievePoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) ([]Policy, error)

	// CreateDataset creates new Dataset
	CreateDataset(ctx context.Context, token string, d Dataset) (Dataset, error)
}

var _ Service = (*policiesService)(nil)

type policiesService struct {
	auth mainflux.AuthServiceClient
	repo Repository
}

func (s policiesService) RetrievePoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) ([]Policy, error) {
	if len(groupIDs) == 0 || ownerID == "" {
		return nil, ErrMalformedEntity
	}
	return s.repo.RetrievePoliciesByGroupID(ctx, groupIDs, ownerID)
}

func (s policiesService) RetrievePolicyByIDInternal(ctx context.Context, policyID string, ownerID string) (Policy, error) {
	if policyID == "" || ownerID == "" {
		return Policy{}, ErrMalformedEntity
	}
	return s.repo.RetrievePolicyByID(ctx, policyID, ownerID)
}

func (s policiesService) CreateDataset(ctx context.Context, token string, d Dataset) (Dataset, error) {
	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Dataset{}, err
	}

	d.MFOwnerID = mfOwnerID

	id, err := s.repo.SaveDataset(ctx, d)
	if err != nil {
		return Dataset{}, errors.Wrap(ErrCreateDataset, err)
	}
	d.ID = id
	return d, nil
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

func (s policiesService) CreatePolicy(ctx context.Context, token string, p Policy, format string, policyData string) (Policy, error) {

	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Policy{}, err
	}

	if !backend.HaveBackend(p.Backend) {
		return Policy{}, errors.Wrap(ErrCreatePolicy, errors.New(fmt.Sprintf("unsupported backend: '%s'", p.Backend)))
	}

	if p.Policy == nil {
		// if not already in json, make sure the back end can convert it
		if !backend.GetBackend(p.Backend).SupportsFormat(format) {
			return Policy{}, errors.Wrap(ErrCreatePolicy,
				errors.New(fmt.Sprintf("unsupported policy format '%s' for given backend '%s'", format, p.Backend)))
		}

		p.Policy, err = backend.GetBackend(p.Backend).ConvertFromFormat(format, policyData)
		if err != nil {
			return Policy{}, errors.Wrap(ErrCreatePolicy, err)
		}
	}

	err = backend.GetBackend(p.Backend).Validate(p.Policy)
	if err != nil {
		return Policy{}, errors.Wrap(ErrCreatePolicy, err)
	}

	p.MFOwnerID = mfOwnerID

	id, err := s.repo.SavePolicy(ctx, p)
	if err != nil {
		return Policy{}, errors.Wrap(ErrCreatePolicy, err)
	}
	p.ID = id
	return p, nil
}

func New(auth mainflux.AuthServiceClient, repo Repository) Service {

	orb.Register()
	pktvisor.Register()

	return &policiesService{
		auth: auth,
		repo: repo,
	}
}
