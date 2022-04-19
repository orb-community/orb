// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/policies"
	"github.com/ns1labs/orb/policies/backend"
	"go.uber.org/zap"
	"strings"
)

const (
	streamID  = "orb.policies"
	streamLen = 1000
)

var (
	ErrValidatePolicy = errors.New("failed to validate policy")
)

var _ policies.Service = (*eventStore)(nil)

type eventStore struct {
	svc    policies.Service
	client *redis.Client
	logger *zap.Logger
}

func (e eventStore) ViewDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) (policies.Dataset, error) {
	return e.svc.ViewDatasetByIDInternal(ctx, ownerID, datasetID)
}

func (e eventStore) RemoveDataset(ctx context.Context, token string, dsID string) (err error) {
	ds, err := e.svc.ViewDatasetByID(ctx, token, dsID)
	if err != nil {
		return err
	}

	if err := e.svc.RemoveDataset(ctx, token, dsID); err != nil {
		return err
	}

	event := removeDatasetEvent{
		id:           dsID,
		ownerID:      ds.MFOwnerID,
		agentGroupID: ds.AgentGroupID,
		policyID:     ds.PolicyID,
		datasetID:    ds.ID,
	}
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}

	err = e.client.XAdd(ctx, record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
		return err
	}

	return nil
}

func (e eventStore) EditDataset(ctx context.Context, token string, ds policies.Dataset) (policies.Dataset, error) {
	return e.svc.EditDataset(ctx, token, ds)
}

func (e eventStore) RemovePolicy(ctx context.Context, token string, policyID string) error {
	policy, err := e.svc.ViewPolicyByID(ctx, token, policyID)
	if err != nil {
		return err
	}
	if err := e.svc.RemovePolicy(ctx, token, policyID); err != nil {
		return err
	}

	datasets, err := e.svc.ListDatasetsByPolicyIDInternal(ctx, policyID, token)
	if err != nil {
		return err
	}

	if len(datasets) == 0 {
		return nil
	}

	var groupsIDs []string
	var ownerID string
	for _, ds := range datasets {
		ownerID = ds.MFOwnerID
		groupsIDs = append(groupsIDs, ds.AgentGroupID)
	}

	event := removePolicyEvent{
		id:       policyID,
		ownerID:  ownerID,
		name:     policy.Name.String(),
		backend:  policy.Backend,
		groupIDs: strings.Join(groupsIDs, ","),
	}
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}
	err = e.client.XAdd(ctx, record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
		return err
	}

	return nil
}

func (e eventStore) ListDatasetsByPolicyIDInternal(ctx context.Context, policyID string, token string) ([]policies.Dataset, error) {
	return e.svc.ListDatasetsByPolicyIDInternal(ctx, policyID, token)
}

func (e eventStore) EditPolicy(ctx context.Context, token string, pol policies.Policy) (policies.Policy, error) {
	res, err := e.svc.EditPolicy(ctx, token, pol)
	if err != nil {
		return policies.Policy{}, err
	}

	datasets, err := e.svc.ListDatasetsByPolicyIDInternal(ctx, res.ID, token)
	if err != nil {
		return policies.Policy{}, err
	}

	var groupsIDs []string
	for _, ds := range datasets {
		groupsIDs = append(groupsIDs, ds.AgentGroupID)
	}

	p, err := e.svc.ViewPolicyByID(ctx, token, pol.ID)
	if err != nil {
		return policies.Policy{}, err
	}
	pol.Backend = p.Backend
	pol.MFOwnerID = p.MFOwnerID

	err = validatePolicyBackend(&pol, pol.Format, pol.PolicyData)
	if err != nil {
		return policies.Policy{}, err
	}

	event := updatePolicyEvent{
		id:       pol.ID,
		ownerID:  pol.MFOwnerID,
		groupIDs: strings.Join(groupsIDs, ","),
	}
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}
	err = e.client.XAdd(ctx, record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
		return res, err
	}

	return res, nil
}

func (e eventStore) AddPolicy(ctx context.Context, token string, p policies.Policy) (policies.Policy, error) {
	return e.svc.AddPolicy(ctx, token, p)
}

func (e eventStore) ViewPolicyByID(ctx context.Context, token string, policyID string) (policies.Policy, error) {
	return e.svc.ViewPolicyByID(ctx, token, policyID)
}

func (e eventStore) ListPolicies(ctx context.Context, token string, pm policies.PageMetadata) (policies.Page, error) {
	return e.svc.ListPolicies(ctx, token, pm)
}

func (e eventStore) ViewPolicyByIDInternal(ctx context.Context, policyID string, ownerID string) (policies.Policy, error) {
	return e.svc.ViewPolicyByIDInternal(ctx, policyID, ownerID)
}

func (e eventStore) ListPoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) ([]policies.PolicyInDataset, error) {
	return e.svc.ListPoliciesByGroupIDInternal(ctx, groupIDs, ownerID)
}

func (e eventStore) AddDataset(ctx context.Context, token string, d policies.Dataset) (policies.Dataset, error) {
	ds, err := e.svc.AddDataset(ctx, token, d)
	if err != nil {
		return ds, err
	}

	event := createDatasetEvent{
		id:           ds.ID,
		ownerID:      ds.MFOwnerID,
		name:         ds.Name.String(),
		agentGroupID: ds.AgentGroupID,
		policyID:     ds.PolicyID,
		sinkIDs:      strings.Join(ds.SinkIDs, ","),
	}
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}
	err = e.client.XAdd(ctx, record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
		return ds, err
	}

	return ds, nil
}

func (e eventStore) InactivateDatasetByGroupID(ctx context.Context, groupID string, ownerID string) error {
	return e.svc.InactivateDatasetByGroupID(ctx, groupID, ownerID)
}

func (e eventStore) ValidatePolicy(ctx context.Context, token string, p policies.Policy) (policies.Policy, error) {
	return e.svc.ValidatePolicy(ctx, token, p)
}

func (e eventStore) DeleteSinkFromAllDatasetsInternal(ctx context.Context, sinkID string, token string) ([]policies.Dataset, error) {
	return e.svc.DeleteSinkFromAllDatasetsInternal(ctx, sinkID, token)
}

func (e eventStore) InactivateDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) error {
	ds, err := e.svc.ViewDatasetByIDInternal(ctx, ownerID, datasetID)
	if err != nil {
		return err
	}

	if err := e.svc.InactivateDatasetByIDInternal(ctx, ownerID, datasetID); err != nil {
		return err
	}

	event := removeDatasetEvent{
		id:           datasetID,
		ownerID:      ds.MFOwnerID,
		agentGroupID: ds.AgentGroupID,
		policyID:     ds.PolicyID,
		datasetID:    ds.ID,
	}
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}

	err = e.client.XAdd(ctx, record).Err()
	if err != nil {
		e.logger.Error("error sending event to event store", zap.Error(err))
		return err
	}

	return nil
}

func (e eventStore) DeleteAgentGroupFromAllDatasets(ctx context.Context, groupID string, token string) error {
	return e.svc.DeleteAgentGroupFromAllDatasets(ctx, groupID, token)
}

func (e eventStore) DuplicatePolicy(ctx context.Context, token string, policyID string, name string) (policies.Policy, error) {
	return e.svc.DuplicatePolicy(ctx, token, policyID, name)
}

// NewEventStoreMiddleware returns wrapper around policies service that sends
// events to event store.
func NewEventStoreMiddleware(svc policies.Service, client *redis.Client, logger *zap.Logger) policies.Service {
	return eventStore{
		logger: logger,
		svc:    svc,
		client: client,
	}
}

func validatePolicyBackend(p *policies.Policy, format string, policyData string) (err error) {
	if !backend.HaveBackend(p.Backend) {
		return errors.Wrap(ErrValidatePolicy, errors.New(fmt.Sprintf("unsupported backend: '%s'", p.Backend)))
	}

	if p.Policy == nil {
		// if not already in json, make sure the back end can convert it
		if !backend.GetBackend(p.Backend).SupportsFormat(format) {
			return errors.Wrap(ErrValidatePolicy,
				errors.New(fmt.Sprintf("unsupported policy format '%s' for given backend '%s'", format, p.Backend)))
		}

		p.Policy, err = backend.GetBackend(p.Backend).ConvertFromFormat(format, policyData)
		if err != nil {
			return errors.Wrap(ErrValidatePolicy, err)
		}
	}

	err = backend.GetBackend(p.Backend).Validate(p.Policy)
	if err != nil {
		return errors.Wrap(ErrValidatePolicy, err)
	}
	return nil
}

func (e eventStore) ValidateDataset(ctx context.Context, token string, d policies.Dataset) (policies.Dataset, error) {
	return e.svc.ValidateDataset(ctx, token, d)
}
func (e eventStore) ListDatasets(ctx context.Context, token string, pm policies.PageMetadata) (policies.PageDataset, error) {
	return e.svc.ListDatasets(ctx, token, pm)
}

func (e eventStore) ViewDatasetByID(ctx context.Context, token string, datasetID string) (policies.Dataset, error) {
	return e.svc.ViewDatasetByID(ctx, token, datasetID)
}
