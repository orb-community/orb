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

func (l loggingMiddleware) InactivateDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: inactivate_dataset_by_id_internal",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: inactivate_dataset_by_id_internal",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.InactivateDatasetByIDInternal(ctx, ownerID, datasetID)
}

func (l loggingMiddleware) ViewDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) (_ policies.Dataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_dataset_by_id_internal",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: view_dataset_by_id_internal",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewDatasetByIDInternal(ctx, ownerID, datasetID)
}

func (l loggingMiddleware) RemoveDataset(ctx context.Context, token string, dsID string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: remove_dataset",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: remove_dataset",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.RemoveDataset(ctx, token, dsID)
}

func (l loggingMiddleware) EditDataset(ctx context.Context, token string, ds policies.Dataset) (_ policies.Dataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: edit_dataset",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: edit_dataset",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.EditDataset(ctx, token, ds)
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

func (l loggingMiddleware) EditPolicy(ctx context.Context, token string, pol policies.Policy) (_ policies.Policy, err error) {
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
	return l.svc.EditPolicy(ctx, token, pol)
}

func (l loggingMiddleware) AddPolicy(ctx context.Context, token string, p policies.Policy) (_ policies.Policy, err error) {
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
	return l.svc.AddPolicy(ctx, token, p)
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

func (l loggingMiddleware) ListPoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) (_ []policies.PolicyInDataset, err error) {
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

func (l loggingMiddleware) ValidatePolicy(ctx context.Context, token string, p policies.Policy) (_ policies.Policy, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: validate_policy",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: validate_policy",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ValidatePolicy(ctx, token, p)
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

func (l loggingMiddleware) ViewDatasetByID(ctx context.Context, token string, datasetID string) (_ policies.Dataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_dataset_by_id",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: view_dataset_by_id",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewDatasetByID(ctx, token, datasetID)
}

func (l loggingMiddleware) ListDatasets(ctx context.Context, token string, pm policies.PageMetadata) (_ policies.PageDataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_dataset",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_dataset",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListDatasets(ctx, token, pm)
}

func (l loggingMiddleware) DeleteSinkFromAllDatasetsInternal(ctx context.Context, sinkID string, ownerID string) (ds []policies.Dataset, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: delete_sink_from_all_datasets",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: delete_sink_from_all_datasets",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.DeleteSinkFromAllDatasetsInternal(ctx, sinkID, ownerID)
}

func (l loggingMiddleware) DeleteAgentGroupFromAllDatasets(ctx context.Context, groupID string, token string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: delete_agent_group_from_all_datasets",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: delete_agent_group_from_all_datasets",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.DeleteAgentGroupFromAllDatasets(ctx, groupID, token)
}

func (l loggingMiddleware) DuplicatePolicy(ctx context.Context, token string, policyID string, name string) (_ policies.Policy, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: duplicate_policy",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: duplicate_policy",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.DuplicatePolicy(ctx, token, policyID, name)
}

func NewLoggingMiddleware(svc policies.Service, logger *zap.Logger) policies.Service {
	return &loggingMiddleware{logger, svc}
}
