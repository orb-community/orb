// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/ns1labs/orb/sinks"
)

func retrieveSinkEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(accessByIDReq)
		sink, err := svc.ViewSinkInternal(ctx, req.OwnerID, req.SinkID)
		if err != nil {
			return sink, err
		}
		tagData, err := json.Marshal(sink.Tags)
		if err != nil {
			return sinkRes{}, err
		}
		configData, err := json.Marshal(sink.Config)
		if err != nil {
			return sinkRes{}, err
		}

		res := sinkRes{
			id:          sink.ID,
			name:        sink.Name.String(),
			description: sink.Description,
			tags:        tagData,
			state:       sink.State.String(),
			error:       sink.Error,
			backend:     sink.Backend,
			config:      configData,
		}
		return res, err
	}
}
