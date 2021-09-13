/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

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

func (l loggingMiddleware) RemovePolicy(ctx context.Context, token string, policyID string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: remove_policy",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: remove_policy",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.RemovePolicy(ctx, token, policyID)
}

func (l loggingMiddleware) ListDatasetsByPolicyIDInternal(ctx context.Context, policyID string, token string) (_ []policies.Dataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_dataset_by_policy_id",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_dataset_by_policy_id",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListDatasetsByPolicyIDInternal(ctx, policyID, token)
}

func (l loggingMiddleware) EditPolicy(ctx context.Context, token string, pol policies.Policy, format string, policyData string) (_ policies.Policy, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: edit_policy",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: edit_policy",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.EditPolicy(ctx, token, pol, format, policyData)
}

func (l loggingMiddleware) AddPolicy(ctx context.Context, token string, p policies.Policy, format string, policyData string) (_ policies.Policy, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: add_policy",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: add_policy",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.AddPolicy(ctx, token, p, format, policyData)
}

func (l loggingMiddleware) ViewPolicyByID(ctx context.Context, token string, policyID string) (_ policies.Policy, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_policy_by_id",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: view_policy_by_id",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewPolicyByID(ctx, token, policyID)
}

func (l loggingMiddleware) ListPolicies(ctx context.Context, token string, pm policies.PageMetadata) (_ policies.Page, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_policies",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_policies",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListPolicies(ctx, token, pm)
}

func (l loggingMiddleware) ViewPolicyByIDInternal(ctx context.Context, policyID string, ownerID string) (_ policies.Policy, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_policy_by_id_internal",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: view_policy_by_id_internal",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewPolicyByIDInternal(ctx, policyID, ownerID)
}

func (l loggingMiddleware) ListPoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) (_ []policies.Policy, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_policies_by_groups",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_policies_by_groups",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListPoliciesByGroupIDInternal(ctx, groupIDs, ownerID)
}

func (l loggingMiddleware) AddDataset(ctx context.Context, token string, d policies.Dataset) (_ policies.Dataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: add_dataset",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: add_dataset",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.AddDataset(ctx, token, d)
}

func (l loggingMiddleware) InactivateDatasetByGroupID(ctx context.Context, groupID string, token string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: inactivate_dataset",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: inactivate_dataset",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.InactivateDatasetByGroupID(ctx, groupID, token)
}

func (l loggingMiddleware) ValidateDataset(ctx context.Context, token string, d policies.Dataset) (_ policies.Dataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: validate_dataset",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: validate_dataset",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ValidateDataset(ctx, token, d)
}

func NewLoggingMiddleware(svc policies.Service, logger *zap.Logger) policies.Service {
	return &loggingMiddleware{logger, svc}
}
