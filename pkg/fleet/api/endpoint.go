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
	"github.com/ns1labs/orb/pkg/fleet"
	"github.com/ns1labs/orb/pkg/types"
)

func addSelectorEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(c context.Context, request interface{}) (interface{}, error) {
		req := request.(addSelectorReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		selector := fleet.Selector{
			Name:     *nID,
			Metadata: req.Metadata,
		}
		saved, err := svc.CreateSelector(c, req.token, selector)
		if err != nil {
			return nil, err
		}

		res := selectorRes{
			name:    saved.Name.String(),
			created: true,
		}

		return res, nil
	}
}
