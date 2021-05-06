/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"github.com/mainflux/mainflux"
)

type Service interface {
	Add() error
}

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

func (s sinkService) Add() error {
	panic("implement me")
}
