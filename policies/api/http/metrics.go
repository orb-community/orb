/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/ns1labs/orb/policies"
)

var _ policies.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     policies.Service
}

func (m metricsMiddleware) ViewDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) (policies.Dataset, error) {
	return m.svc.ViewDatasetByIDInternal(ctx, ownerID, datasetID)
}

func (m metricsMiddleware) RemoveDataset(ctx context.Context, token string, dsID string) error {
	return m.svc.RemoveDataset(ctx, token, dsID)
}

func (m metricsMiddleware) EditDataset(ctx context.Context, token string, ds policies.Dataset) (policies.Dataset, error) {
	return m.svc.EditDataset(ctx, token, ds)
}

func (m metricsMiddleware) RemovePolicy(ctx context.Context, token string, policyID string) error {
	return m.svc.RemovePolicy(ctx, token, policyID)
}

func (m metricsMiddleware) ListDatasetsByPolicyIDInternal(ctx context.Context, policyID string, owner string) ([]policies.Dataset, error) {
	return m.svc.ListDatasetsByPolicyIDInternal(ctx, policyID, owner)
}

func (m metricsMiddleware) ListDatasetsByPolicyID(ctx context.Context, policyID string, token string) ([]policies.Dataset, error) {
	return m.svc.ListDatasetsByPolicyIDInternal(ctx, policyID, token)
}

func (m metricsMiddleware) EditPolicy(ctx context.Context, token string, pol policies.Policy, format string, policyData string) (policies.Policy, error) {
	return m.svc.EditPolicy(ctx, token, pol, format, policyData)
}

func (m metricsMiddleware) ListPolicies(ctx context.Context, token string, pm policies.PageMetadata) (policies.Page, error) {
	return m.svc.ListPolicies(ctx, token, pm)
}

func (m metricsMiddleware) ViewPolicyByID(ctx context.Context, token string, policyID string) (policies.Policy, error) {
	return m.svc.ViewPolicyByID(ctx, token, policyID)
}

func (m metricsMiddleware) ListPoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) ([]policies.PolicyInDataset, error) {
	return m.svc.ListPoliciesByGroupIDInternal(ctx, groupIDs, ownerID)
}

func (m metricsMiddleware) ViewPolicyByIDInternal(ctx context.Context, policyID string, ownerID string) (policies.Policy, error) {
	return m.svc.ViewPolicyByIDInternal(ctx, policyID, ownerID)
}

func (m metricsMiddleware) AddDataset(ctx context.Context, token string, d policies.Dataset) (policies.Dataset, error) {
	return m.svc.AddDataset(ctx, token, d)
}

func (m metricsMiddleware) AddPolicy(ctx context.Context, token string, p policies.Policy, format string, policyData string) (policies.Policy, error) {
	return m.svc.AddPolicy(ctx, token, p, format, policyData)
}

func (m metricsMiddleware) InactivateDatasetByGroupID(ctx context.Context, groupID string, ownerID string) error {
	return m.svc.InactivateDatasetByGroupID(ctx, groupID, ownerID)
}

func (m metricsMiddleware) ValidatePolicy(ctx context.Context, token string, p policies.Policy, format string, policyData string) (policies.Policy, error) {
	return m.svc.ValidatePolicy(ctx, token, p, format, policyData)
}

func (m metricsMiddleware) ValidateDataset(ctx context.Context, token string, d policies.Dataset) (policies.Dataset, error) {
	return m.svc.ValidateDataset(ctx, token, d)
}

func (m metricsMiddleware) ViewDatasetByID(ctx context.Context, token string, datasetID string) (policies.Dataset, error) {
	return m.svc.ViewDatasetByID(ctx, token, datasetID)
}

func (m metricsMiddleware) ListDatasets(ctx context.Context, token string, pm policies.PageMetadata) (policies.PageDataset, error) {
	return m.svc.ListDatasets(ctx, token, pm)
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc policies.Service, counter metrics.Counter, latency metrics.Histogram) policies.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
