// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package maestro

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var _ MaestroService = (*maestroService)(nil)

type maestroService struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

func NewMaestroService(logger *zap.Logger, redisClient *redis.Client) MaestroService {
	return &maestroService{
		logger:      logger,
		redisClient: redisClient,
	}
}
