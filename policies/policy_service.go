// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"context"
	"fmt"
	"google.golang.org/grpc/status"

	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies/backend"
	sinkpb "github.com/ns1labs/orb/sinks/pb"
)

var (
	ErrCreatePolicy            = errors.New("failed to create policy")
	ErrValidatePolicy          = errors.New("failed to validate policy")
	ErrCreateDataset           = errors.New("failed to create dataset")
	ErrInactivateDataset       = errors.New("failed to inactivate dataset")
	ErrUpdateEntity            = errors.New("failed to update entity")
	ErrRemoveEntity            = errors.New("failed to remove entity")
	ErrMalformedEntity         = errors.New("malformed entity")
	ErrNotFound                = errors.New("non-existent entity")
	ErrUnauthorizedAccess      = errors.New("missing or invalid credentials provided")
	ErrNotifyAgentGroupChannel = errors.New("failed to notify agent group channel")
)

func (s policiesService) ListDatasetsByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) ([]Dataset, error) {
	if len(groupIDs) == 0 || ownerID == "" {
		return nil, ErrMalformedEntity
	}
	return s.repo.RetrieveDatasetsByGroupID(ctx, groupIDs, ownerID)
}

func (s policiesService) ListPolicies(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	ownerID, err := s.identify(token)
	if err != nil {
		return Page{}, err
	}
	return s.repo.RetrieveAll(ctx, ownerID, pm)
}

func (s policiesService) ListPoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) ([]PolicyInDataset, error) {
	if len(groupIDs) == 0 || ownerID == "" {
		return nil, ErrMalformedEntity
	}
	return s.repo.RetrievePoliciesByGroupID(ctx, groupIDs, ownerID)
}

func (s policiesService) ViewPolicyByIDInternal(ctx context.Context, policyID string, ownerID string) (Policy, error) {
	if policyID == "" || ownerID == "" {
		return Policy{}, ErrMalformedEntity
	}
	return s.repo.RetrievePolicyByID(ctx, policyID, ownerID)
}

func (s policiesService) AddDataset(ctx context.Context, token string, d Dataset) (Dataset, error) {
	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Dataset{}, err
	}

	d.MFOwnerID = mfOwnerID

	id, err := s.repo.SaveDataset(ctx, d)
	if err != nil {
		return Dataset{}, errors.Wrap(ErrCreateDataset, err)
	}
	d.ID = id
	return d, nil
}

func (s policiesService) InactivateDatasetByGroupID(ctx context.Context, groupID string, token string) error {
	ownerID, err := s.identify(token)
	if err != nil {
		return err
	}

	if groupID == "" {
		return ErrMalformedEntity
	}
	return s.repo.InactivateDatasetByGroupID(ctx, groupID, ownerID)
}

func (s policiesService) AddPolicy(ctx context.Context, token string, p Policy) (Policy, error) {

	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Policy{}, err
	}

	err = validatePolicyBackend(&p)
	if err != nil {
		return Policy{}, err
	}

	p.MFOwnerID = mfOwnerID

	id, err := s.repo.SavePolicy(ctx, p)
	if err != nil {
		return Policy{}, errors.Wrap(ErrCreatePolicy, err)
	}
	p.ID = id
	return p, nil
}

func (s policiesService) ViewPolicyByID(ctx context.Context, token string, policyID string) (Policy, error) {
	ownerID, err := s.identify(token)
	if err != nil {
		return Policy{}, err
	}

	res, err := s.repo.RetrievePolicyByID(ctx, policyID, ownerID)
	if err != nil {
		return Policy{}, err
	}
	return res, nil
}

func (s policiesService) EditPolicy(ctx context.Context, token string, pol Policy) (Policy, error) {
	ownerID, err := s.identify(token)
	if err != nil {
		return Policy{}, err
	}

	// Used to get the policy backend and validate it
	currentPol, err := s.repo.RetrievePolicyByID(ctx, pol.ID, ownerID)
	if err != nil {
		return Policy{}, err
	}
	pol.Backend = currentPol.Backend
	pol.MFOwnerID = ownerID
	pol.Version = currentPol.Version

	// If backend policy is not being edited, retrieve saved one
	if pol.PolicyData == "" && pol.Policy == nil {
		pol.Policy = currentPol.Policy
		pol.PolicyData = currentPol.PolicyData
		pol.Format = currentPol.Format
	}

	err = validatePolicyBackend(&pol)
	if err != nil {
		return Policy{}, err
	}

	// If policy name is not being edited, retrieve saved one
	if pol.Name.String() == "" {
		pol.Name = currentPol.Name
	}

	pol.Version++
	err = s.repo.UpdatePolicy(ctx, ownerID, pol)
	if err != nil {
		return Policy{}, err
	}

	// Used to return the updated policy
	res, err := s.repo.RetrievePolicyByID(ctx, pol.ID, ownerID)
	if err != nil {
		return Policy{}, err
	}

	return res, nil
}

func (s policiesService) ListDatasetsByPolicyIDInternal(ctx context.Context, policyID string, token string) ([]Dataset, error) {
	ownerID, err := s.identify(token)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.RetrieveDatasetsByPolicyID(ctx, policyID, ownerID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s policiesService) RemovePolicy(ctx context.Context, token string, policyID string) error {
	ownerID, err := s.identify(token)
	if err != nil {
		return err
	}

	err = s.repo.DeletePolicy(ctx, ownerID, policyID)
	if err != nil {
		return err
	}

	return nil
}

func (s policiesService) RemoveAllDatasetsByPolicyIDInternal(ctx context.Context, token string, policyID string) error {
	ownerID, err := s.identify(token)
	if err != nil {
		return err
	}

	err = s.repo.DeleteAllDatasetsPolicy(ctx, policyID, ownerID)
	if err != nil {
		return err
	}

	return nil
}

func (s policiesService) ViewDatasetByID(ctx context.Context, token string, datasetID string) (Dataset, error) {
	ownerID, err := s.identify(token)
	if err != nil {
		return Dataset{}, err
	}

	res, err := s.repo.RetrieveDatasetByID(ctx, datasetID, ownerID)
	if err != nil {
		return Dataset{}, err
	}
	return res, nil
}

func (s policiesService) ViewDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) (Dataset, error) {
	res, err := s.repo.RetrieveDatasetByID(ctx, datasetID, ownerID)
	if err != nil {
		return Dataset{}, err
	}
	return res, nil
}
func validatePolicyBackend(p *Policy) (err error) {
	if !backend.HaveBackend(p.Backend) {
		return errors.Wrap(ErrValidatePolicy, errors.New(fmt.Sprintf("unsupported backend: '%s'", p.Backend)))
	}

	if p.Policy == nil {
		// if not already in json, make sure the back end can convert it
		if !backend.GetBackend(p.Backend).SupportsFormat(p.Format) {
			return errors.Wrap(ErrValidatePolicy,
				errors.New(fmt.Sprintf("unsupported policy format '%s' for given backend '%s'", p.Format, p.Backend)))
		}

		p.Policy, err = backend.GetBackend(p.Backend).ConvertFromFormat(p.Format, p.PolicyData)
		if err != nil {
			return errors.Wrap(ErrValidatePolicy, err)
		}
	} else {
		// policy was already received as a json
		p.Format = "json"
	}

	err = backend.GetBackend(p.Backend).Validate(p.Policy)
	if err != nil {
		return errors.Wrap(ErrCreatePolicy, err)
	}
	return nil
}

func (s policiesService) ValidatePolicy(ctx context.Context, token string, p Policy) (Policy, error) {

	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Policy{}, err
	}

	err = validatePolicyBackend(&p)
	if err != nil {
		return p, errors.Wrap(ErrCreatePolicy, err)
	}

	p.MFOwnerID = mfOwnerID

	return p, nil
}

func (s policiesService) EditDataset(ctx context.Context, token string, ds Dataset) (Dataset, error) {
	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Dataset{}, err
	}
	ds.MFOwnerID = mfOwnerID

	err = s.validateDatasetSink(ctx, ds.MFOwnerID, ds.SinkIDs)
	if err != nil {
		return Dataset{}, err
	}

	err = s.repo.UpdateDataset(ctx, mfOwnerID, ds)
	if err != nil {
		return Dataset{}, err
	}

	datasetEdited, err := s.repo.RetrieveDatasetByID(ctx, ds.ID, ds.MFOwnerID)
	if err != nil {
		return Dataset{}, err
	}

	errValidatePolicy := s.validateDatasetPolicy(ctx, datasetEdited.MFOwnerID, datasetEdited.PolicyID)
	errValidateAGroup := s.validateDatasetAgentGroup(ctx, datasetEdited.MFOwnerID, datasetEdited.AgentGroupID)

	if errValidatePolicy == nil && errValidateAGroup == nil {
		err = s.repo.ActivateDatasetByID(ctx, datasetEdited.ID, datasetEdited.MFOwnerID)
		if err != nil {
			return Dataset{}, err
		}
		datasetEdited.Valid = true
	}

	return datasetEdited, nil
}

func (s policiesService) RemoveDataset(ctx context.Context, token string, dsID string) error {
	mfOwnerID, err := s.identify(token)
	if err != nil {
		return err
	}
	err = s.repo.DeleteDataset(ctx, mfOwnerID, dsID)
	if err != nil {
		return err
	}
	return nil
}

func (s policiesService) ValidateDataset(ctx context.Context, token string, d Dataset) (Dataset, error) {
	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Dataset{}, err
	}

	d.MFOwnerID = mfOwnerID

	err = s.validateDatasetSink(ctx, d.MFOwnerID, d.SinkIDs)
	if err != nil {
		return Dataset{}, err
	}

	err = s.validateDatasetPolicy(ctx, d.MFOwnerID, d.PolicyID)
	if err != nil {
		return Dataset{}, err
	}

	err = s.validateDatasetAgentGroup(ctx, d.MFOwnerID, d.AgentGroupID)
	if err != nil {
		return Dataset{}, err
	}

	return d, nil
}

func (s policiesService) ListDatasets(ctx context.Context, token string, pm PageMetadata) (PageDataset, error) {
	ownerID, err := s.identify(token)
	if err != nil {
		return PageDataset{}, err
	}
	return s.repo.RetrieveAllDatasetsByOwner(ctx, ownerID, pm)
}

func (s policiesService) DeleteSinkFromAllDatasetsInternal(ctx context.Context, sinkID string, ownerID string) ([]Dataset, error) {
	if sinkID == "" || ownerID == "" {
		return []Dataset{}, ErrMalformedEntity
	}

	datasets, err := s.repo.DeleteSinkFromAllDatasets(ctx, sinkID, ownerID)
	if err != nil {
		return []Dataset{}, err
	}

	return datasets, nil
}

func (s policiesService) InactivateDatasetByIDInternal(ctx context.Context, ownerID string, datasetID string) error {
	if datasetID == "" || ownerID == "" {
		return ErrMalformedEntity
	}

	err := s.repo.InactivateDatasetByID(ctx, datasetID, ownerID)
	if err != nil {
		return errors.Wrap(ErrInactivateDataset, err)
	}

	return nil
}

func (s policiesService) validateDatasetSink(ctx context.Context, ownerID string, sinkIDs []string) error {

	if len(sinkIDs) == 0 {
		return errors.Wrap(ErrMalformedEntity, errors.New("empty sink IDs"))
	}
	for _, sinkID := range sinkIDs {
		_, err := uuid.FromString(sinkID)
		if err != nil {
			return errors.Wrap(errors.New("invalid sink id"), ErrMalformedEntity)
		}

		_, err = s.sinksGrpcClient.RetrieveSink(ctx, &sinkpb.SinkByIDReq{
			SinkID:  sinkID,
			OwnerID: ownerID,
		})
		if err != nil {
			return errors.Wrap(errors.New("sink id does not exist"), err)
		}
	}
	return nil
}

func (s policiesService) validateDatasetPolicy(ctx context.Context, ownerID string, policyID string) error {
	_, err := uuid.FromString(policyID)
	if err != nil {
		return errors.Wrap(errors.New("invalid policy id"), ErrMalformedEntity)
	}

	_, err = s.repo.RetrievePolicyByID(ctx, policyID, ownerID)
	if err != nil {
		return errors.Wrap(errors.New("policy id does not exist"), err)
	}
	return nil
}

func (s policiesService) validateDatasetAgentGroup(ctx context.Context, ownerID string, aGroupID string) error {
	_, err := uuid.FromString(aGroupID)
	if err != nil {
		return errors.Wrap(errors.New("invalid agent group id"), ErrMalformedEntity)
	}

	_, err = s.fleetGrpcClient.RetrieveAgentGroup(ctx, &pb.AgentGroupByIDReq{
		AgentGroupID: aGroupID,
		OwnerID:      ownerID,
	})
	if err != nil {
		return errors.Wrap(errors.New("agent group id does not exist"), err)
	}
	return nil
}

func (s policiesService) DeleteAgentGroupFromAllDatasets(ctx context.Context, groupID string, token string) error {
	ownerID, err := s.identify(token)
	if err != nil {
		return err
	}

	if groupID == "" {
		return ErrMalformedEntity
	}

	err = s.repo.DeleteAgentGroupFromAllDatasets(ctx, groupID, ownerID)
	if err != nil {
		return err
	}

	return nil
}

func (s policiesService) DuplicatePolicy(ctx context.Context, token string, policyID string, name string) (Policy, error) {

	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Policy{}, err
	}

	existingPolicy, err := s.repo.RetrievePolicyByID(ctx, policyID, mfOwnerID)
	if err != nil {
		return Policy{}, err
	}

	policy := existingPolicy
	policy.Version = 0

	var nameSuffix string
	var id string
	var errCreate error

	if name != "" {
		policyName, err := types.NewIdentifier(name)
		if err != nil {
			return Policy{}, err
		}
		policy.Name = policyName
		id, err = s.repo.SavePolicy(ctx, policy)
		if err != nil {
			return Policy{}, errors.Wrap(ErrCreatePolicy, err)
		}
	} else {
		nameSuffix = fmt.Sprintf("_copy")
		var i = 1
		for {
			policyName, err := types.NewIdentifier(existingPolicy.Name.String() + nameSuffix)
			if err != nil {
				return Policy{}, err
			}

			policy.Name = policyName
			id, errCreate = s.repo.SavePolicy(ctx, policy)
			if errCreate != nil && status.Code(errCreate) == status.Code(errors.ErrConflict) {
				if i < 3 {
					i++
					nameSuffix = fmt.Sprintf("_copy%d", i)
					continue
				} else {
					return Policy{}, errors.Wrap(errors.New("limit of copies of a single policy exceeded"), errors.ErrConflict)
				}
			}
			break
		}
		if errCreate != nil {
			return Policy{}, errCreate
		}
	}
	policy.ID = id

	return policy, nil
}
