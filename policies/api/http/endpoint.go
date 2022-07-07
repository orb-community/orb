// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
)

func addPolicyEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addPolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		policy := policies.Policy{
			Name:          nID,
			Backend:       req.Backend,
			SchemaVersion: req.SchemaVersion,
			Policy:        req.Policy,
			Description:   req.Description,
			OrbTags:       req.Tags,
			PolicyData:    req.PolicyData,
			Format:        req.Format,
		}

		saved, err := svc.AddPolicy(ctx, req.token, policy)
		if err != nil {
			if e, ok := err.(errors.Error); ok && errors.Contains(err, policies.ErrValidatePolicy) {
				err = errors.New(fmt.Sprintf("%s : %s", e.Msg(), e.Err()))
			}
			return nil, err
		}

		res := policyRes{
			ID:            saved.ID,
			Name:          saved.Name.String(),
			Description:   saved.Description,
			Tags:          saved.OrbTags,
			Backend:       saved.Backend,
			SchemaVersion: saved.SchemaVersion,
			Policy:        saved.Policy,
			Format:        saved.Format,
			PolicyData:    saved.PolicyData,
			created:       true,
		}

		return res, nil
	}
}

func viewPolicyEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		policy, err := svc.ViewPolicyByID(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := policyRes{
			ID:            policy.ID,
			Name:          policy.Name.String(),
			Description:   policy.Description,
			Tags:          policy.OrbTags,
			Backend:       policy.Backend,
			SchemaVersion: policy.SchemaVersion,
			Policy:        policy.Policy,
			Version:       policy.Version,
			LastModified:  policy.LastModified,
			PolicyData:    policy.PolicyData,
			Format:        policy.Format,
			Created:       policy.Created,
		}
		return res, nil
	}
}

func listPoliciesEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listResourcesReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListPolicies(ctx, req.token, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		res := policiesPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
				Order:  page.Order,
				Dir:    page.Dir,
			},
			Policies: []policyRes{},
		}
		for _, ag := range page.Policies {
			view := policyRes{
				ID:            ag.ID,
				Name:          ag.Name.String(),
				Description:   ag.Description,
				Version:       ag.Version,
				Backend:       ag.Backend,
				SchemaVersion: ag.SchemaVersion,
				LastModified:  ag.LastModified,
				Created:       ag.Created,
			}
			res.Policies = append(res.Policies, view)
		}
		return res, nil
	}
}

func editPoliciyEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updatePolicyReq)
		if err := req.validate(); err != nil {
			return policyUpdateRes{}, err
		}

		var nameID types.Identifier
		if req.Name != "" {
			nameID, err = types.NewIdentifier(req.Name)
			if err != nil {
				return policyUpdateRes{}, errors.Wrap(errors.ErrMalformedEntity, err)
			}
		}
		plcy := policies.Policy{
			ID:          req.id,
			Name:        nameID,
			Description: req.Description,
			OrbTags:     req.Tags,
			Policy:      req.Policy,
			PolicyData:  req.PolicyData,
			Format:      req.Format,
		}

		res, err := svc.EditPolicy(ctx, req.token, plcy)
		if err != nil {
			return policyUpdateRes{}, err
		}

		plcyRes := policyUpdateRes{
			ID:          res.ID,
			Name:        res.Name.String(),
			Description: res.Description,
			Tags:        res.OrbTags,
			Policy:      res.Policy,
			Format:      res.Format,
			PolicyData:  res.PolicyData,
			Version:     res.Version,
		}

		return plcyRes, nil
	}
}

func removePolicyEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return removeRes{}, err
		}
		err = svc.RemovePolicy(ctx, req.token, req.id)
		if err != nil {
			return removeRes{}, err
		}
		return removeRes{}, nil
	}
}

func addDatasetEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addDatasetReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		d := policies.Dataset{
			Name:         nID,
			AgentGroupID: req.AgentGroupID,
			PolicyID:     req.PolicyID,
			SinkIDs:      req.SinkIDs,
		}

		saved, err := svc.AddDataset(ctx, req.token, d)
		if err != nil {
			return nil, err
		}

		res := datasetRes{
			ID:           saved.ID,
			Name:         saved.Name.String(),
			Valid:        saved.Valid,
			AgentGroupID: saved.AgentGroupID,
			PolicyID:     saved.PolicyID,
			SinkIDs:      saved.SinkIDs,
			Metadata:     saved.Metadata,
			TsCreated:    saved.Created,
			Tags:         saved.Tags,
			created:      true,
		}

		return res, nil
	}
}

func editDatasetEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateDatasetReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		dataset := policies.Dataset{
			Name:    nID,
			ID:      req.id,
			Tags:    req.Tags,
			SinkIDs: req.SinkIDs,
		}

		ds, err := svc.EditDataset(ctx, req.token, dataset)
		if err != nil {
			return nil, err
		}

		res := datasetRes{
			ID:           ds.ID,
			Name:         ds.Name.String(),
			Valid:        ds.Valid,
			AgentGroupID: ds.AgentGroupID,
			PolicyID:     ds.PolicyID,
			SinkIDs:      ds.SinkIDs,
			Metadata:     ds.Metadata,
			TsCreated:    ds.Created,
			Tags:         ds.Tags,
		}

		return res, nil
	}
}

func validatePolicyEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addPolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		policy := policies.Policy{
			Name:        nID,
			Backend:     req.Backend,
			Policy:      req.Policy,
			OrbTags:     req.Tags,
			Description: req.Description,
			Format:      req.Format,
			PolicyData:  req.PolicyData,
		}

		validated, err := svc.ValidatePolicy(ctx, req.token, policy)
		if err != nil {
			return nil, err
		}

		res := policyValidateRes{
			Name:        validated.Name.String(),
			Backend:     validated.Backend,
			Tags:        validated.OrbTags,
			Policy:      validated.Policy,
			PolicyData:  validated.PolicyData,
			Format:      validated.Format,
			Description: validated.Description,
		}

		return res, nil
	}
}

func removeDatasetEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return removeRes{}, err
		}
		if err := svc.RemoveDataset(ctx, req.token, req.id); err != nil {
			return removeRes{}, err
		}
		return removeRes{}, nil
	}
}

func validateDatasetEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addDatasetReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		d := policies.Dataset{
			Name:         nID,
			AgentGroupID: req.AgentGroupID,
			PolicyID:     req.PolicyID,
			SinkIDs:      req.SinkIDs,
			Tags:         req.Tags,
		}

		validated, err := svc.ValidateDataset(ctx, req.token, d)
		if err != nil {
			return nil, err
		}

		res := validateDatasetRes{
			Name:         validated.Name.String(),
			Valid:        true,
			Tags:         validated.Tags,
			AgentGroupID: validated.AgentGroupID,
			PolicyID:     validated.PolicyID,
			SinkIDs:      validated.SinkIDs,
		}

		return res, nil
	}
}

func viewDatasetEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		dataset, err := svc.ViewDatasetByID(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := datasetRes{
			ID:           dataset.ID,
			Name:         dataset.Name.String(),
			PolicyID:     dataset.PolicyID,
			SinkIDs:      dataset.SinkIDs,
			AgentGroupID: dataset.AgentGroupID,
			Valid:        dataset.Valid,
			TsCreated:    dataset.Created,
		}
		return res, nil
	}
}

func listDatasetEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listResourcesReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListDatasets(ctx, req.token, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		res := datasetPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
				Order:  page.Order,
				Dir:    page.Dir,
			},
			Datasets: []datasetRes{},
		}
		for _, dataset := range page.Datasets {
			view := datasetRes{
				ID:           dataset.ID,
				Name:         dataset.Name.String(),
				PolicyID:     dataset.PolicyID,
				SinkIDs:      dataset.SinkIDs,
				AgentGroupID: dataset.AgentGroupID,
				TsCreated:    dataset.Created,
				Valid:        dataset.Valid,
			}
			res.Datasets = append(res.Datasets, view)
		}
		return res, nil
	}
}

func duplicatePolicyEndpoint(svc policies.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(duplicatePolicyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		duplicatedPolicy, err := svc.DuplicatePolicy(ctx, req.token, req.id, req.Name)
		if err != nil {
			return nil, err
		}

		res := policyRes{
			ID:            duplicatedPolicy.ID,
			Name:          duplicatedPolicy.Name.String(),
			Description:   duplicatedPolicy.Description,
			Tags:          duplicatedPolicy.OrbTags,
			Backend:       duplicatedPolicy.Backend,
			SchemaVersion: duplicatedPolicy.SchemaVersion,
			Policy:        duplicatedPolicy.Policy,
			Version:       duplicatedPolicy.Version,
			created:       true,
		}

		return res, nil
	}
}
