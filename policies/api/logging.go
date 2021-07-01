/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/ns1labs/orb/policies"
	"go.uber.org/zap"
	"time"
)

var _ policies.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    policies.Service
}

func (l loggingMiddleware) CreateDataset(ctx context.Context, token string, d policies.Dataset) (_ policies.Dataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: create_dataset",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: create_dataset",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.CreateDataset(ctx, token, d)
}

func (l loggingMiddleware) CreatePolicy(ctx context.Context, token string, p policies.Policy, format string, policyData string) (_ policies.Policy, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: create_policy",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: create_policy",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.CreatePolicy(ctx, token, p, format, policyData)
}

func NewLoggingMiddleware(svc policies.Service, logger *zap.Logger) policies.Service {
	return &loggingMiddleware{logger, svc}
}
