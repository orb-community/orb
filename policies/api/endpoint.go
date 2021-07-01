// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/go-kit/kit/endpoint"
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
			Name:    nID,
			Backend: req.Backend,
			Policy:  req.Policy,
		}

		saved, err := svc.CreatePolicy(ctx, req.token, policy, req.Format, req.PolicyData)
		if err != nil {
			return nil, err
		}

		res := policyRes{
			ID:      saved.ID,
			Name:    saved.Name.String(),
			Backend: saved.Backend,
			created: true,
		}

		return res, nil
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
			SinkID:       req.SinkID,
		}

		saved, err := svc.CreateDataset(ctx, req.token, d)
		if err != nil {
			return nil, err
		}

		res := datasetRes{
			ID:      saved.ID,
			Name:    saved.Name.String(),
			created: true,
		}

		return res, nil
	}
}
