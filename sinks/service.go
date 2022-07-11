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
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks/backend/prometheus"
	"go.uber.org/zap"
	"time"
)

// PageMetadata contains page metadata that helps navigation
type PageMetadata struct {
	Total    uint64
	Offset   uint64         `json:"offset,omitempty"`
	Limit    uint64         `json:"limit,omitempty"`
	Name     string         `json:"name,omitempty"`
	Order    string         `json:"order,omitempty"`
	Dir      string         `json:"dir,omitempty"`
	Metadata types.Metadata `json:"metadata,omitempty"`
	Tags     types.Tags     `json:"tags,omitempty"`
}

var _ SinkService = (*sinkService)(nil)

type sinkService struct {
	logger *zap.Logger
	// for AuthN/AuthZ
	auth  mainflux.AuthServiceClient
	mfsdk mfsdk.SDK
	// Sinks
	sinkRepo SinkRepository
	// passwordService
	passwordService PasswordService
}

func (s sinkService) identify(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return res.GetId(), nil
}

func NewSinkService(logger *zap.Logger, auth mainflux.AuthServiceClient, sinkRepo SinkRepository, mfsdk mfsdk.SDK, services PasswordService) SinkService {

	prometheus.Register()

	return &sinkService{
		logger:          logger,
		auth:            auth,
		sinkRepo:        sinkRepo,
		mfsdk:           mfsdk,
		passwordService: services,
	}
}
