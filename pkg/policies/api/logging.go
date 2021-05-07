/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/policies"
	"go.uber.org/zap"
)

var _ policies.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    policies.Service
}

func (l loggingMiddleware) Add() (policies.Policy, error) {
	panic("implement me")
}

func NewLoggingMiddleware(svc policies.Service, logger *zap.Logger) policies.Service {
	return &loggingMiddleware{logger, svc}
}
