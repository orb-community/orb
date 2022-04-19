/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/policies"
	"time"
)

var _ policies.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	auth    mainflux.AuthServiceClient
	counter metrics.Counter
	latency metrics.Histogram
	svc     policies.Service
}

func (m metricsMiddleware) InactivateDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "inactivateDatasetByIDInternal",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", datasetID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.InactivateDatasetByIDInternal(ctx, ownerID, datasetID)
}

func (m metricsMiddleware) ViewDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) (policies.Dataset, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewDatasetByIDInternal",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", datasetID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewDatasetByIDInternal(ctx, ownerID, datasetID)
}

func (m metricsMiddleware) RemoveDataset(ctx context.Context, token string, dsID string) error {
	ownerID, err := m.identify(token)
	if err != nil {
		return err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "removeDataset",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", dsID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.RemoveDataset(ctx, token, dsID)
}

func (m metricsMiddleware) EditDataset(ctx context.Context, token string, ds policies.Dataset) (policies.Dataset, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Dataset{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "editDataset",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", ds.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.EditDataset(ctx, token, ds)
}

func (m metricsMiddleware) RemovePolicy(ctx context.Context, token string, policyID string) error {
	ownerID, err := m.identify(token)
	if err != nil {
		return err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "removePolicy",
			"owner_id", ownerID,
			"policy_id", policyID,
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.RemovePolicy(ctx, token, policyID)
}

func (m metricsMiddleware) ListDatasetsByPolicyIDInternal(ctx context.Context, policyID string, token string) ([]policies.Dataset, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return nil, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "listDatasetsByPolicyIDInternal",
			"owner_id", ownerID,
			"policy_id", policyID,
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListDatasetsByPolicyIDInternal(ctx, policyID, token)
}

func (m metricsMiddleware) EditPolicy(ctx context.Context, token string, pol policies.Policy) (policies.Policy, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Policy{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "editPolicy",
			"owner_id", ownerID,
			"policy_id", pol.ID,
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.EditPolicy(ctx, token, pol)
}

func (m metricsMiddleware) ListPolicies(ctx context.Context, token string, pm policies.PageMetadata) (policies.Page, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Page{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "listPolicies",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListPolicies(ctx, token, pm)
}

func (m metricsMiddleware) ViewPolicyByID(ctx context.Context, token string, policyID string) (policy policies.Policy, _ error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Policy{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "viewPolicyByID",
			"owner_id", ownerID,
			"policy_id", policyID,
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewPolicyByID(ctx, token, policyID)
}

func (m metricsMiddleware) ListPoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) ([]policies.PolicyInDataset, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "listPoliciesByGroupIDInternal",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListPoliciesByGroupIDInternal(ctx, groupIDs, ownerID)
}

func (m metricsMiddleware) ViewPolicyByIDInternal(ctx context.Context, policyID string, ownerID string) (policies.Policy, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewPolicyByIDInternal",
			"owner_id", ownerID,
			"policy_id", policyID,
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewPolicyByIDInternal(ctx, policyID, ownerID)
}

func (m metricsMiddleware) AddDataset(ctx context.Context, token string, d policies.Dataset) (dataset policies.Dataset, _ error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Dataset{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "addDataset",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", dataset.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.AddDataset(ctx, token, d)
}

func (m metricsMiddleware) AddPolicy(ctx context.Context, token string, p policies.Policy) (policy policies.Policy, _ error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Policy{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "addPolicy",
			"owner_id", ownerID,
			"policy_id", policy.ID,
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.AddPolicy(ctx, token, p)
}

func (m metricsMiddleware) InactivateDatasetByGroupID(ctx context.Context, groupID string, ownerID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "inactivateDatasetByGroupID",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.InactivateDatasetByGroupID(ctx, groupID, ownerID)
}

func (m metricsMiddleware) ValidatePolicy(ctx context.Context, token string, p policies.Policy) (policies.Policy, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Policy{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "validatePolicy",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ValidatePolicy(ctx, token, p)
}

func (m metricsMiddleware) ValidateDataset(ctx context.Context, token string, d policies.Dataset) (policies.Dataset, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Dataset{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "validateDataset",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ValidateDataset(ctx, token, d)
}

func (m metricsMiddleware) ViewDatasetByID(ctx context.Context, token string, datasetID string) (policies.Dataset, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Dataset{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "viewDatasetByID",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", datasetID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewDatasetByID(ctx, token, datasetID)
}

func (m metricsMiddleware) ListDatasets(ctx context.Context, token string, pm policies.PageMetadata) (policies.PageDataset, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.PageDataset{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "listDatasets",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListDatasets(ctx, token, pm)
}

func (m metricsMiddleware) DeleteSinkFromAllDatasetsInternal(ctx context.Context, sinkID string, ownerID string) ([]policies.Dataset, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "deleteSinkFromAllDatasetsInternal",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.DeleteSinkFromAllDatasetsInternal(ctx, sinkID, ownerID)
}

func (m metricsMiddleware) DeleteAgentGroupFromAllDatasets(ctx context.Context, groupID string, token string) error {
	ownerID, err := m.identify(token)
	if err != nil {
		return err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "deleteAgentGroupFromAllDatasets",
			"owner_id", ownerID,
			"policy_id", "",
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.DeleteAgentGroupFromAllDatasets(ctx, groupID, token)
}

func (m metricsMiddleware) DuplicatePolicy(ctx context.Context, token string, policyID string, name string) (policies.Policy, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return policies.Policy{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "duplicatePolicy",
			"owner_id", ownerID,
			"policy_id", policyID,
			"dataset_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.DuplicatePolicy(ctx, token, policyID, name)
}

func (m metricsMiddleware) identify(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return res.GetId(), nil
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(auth mainflux.AuthServiceClient, svc policies.Service, counter metrics.Counter, latency metrics.Histogram) policies.Service {
	return &metricsMiddleware{
		auth:    auth,
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
