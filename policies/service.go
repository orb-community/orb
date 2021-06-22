// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"github.com/mainflux/mainflux"
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
