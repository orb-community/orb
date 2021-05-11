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
)

func addEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(c context.Context, request interface{}) (interface{}, error) {
		req := request.(addAgentReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		// TODO
		saved, err := svc.CreateAgent(c, req.token, fleet.Agent{})
		if err != nil {
			return nil, err
		}

		res := fleetRes{
			id:      saved.MFThingID,
			created: true,
		}

		return res, nil
	}
}
