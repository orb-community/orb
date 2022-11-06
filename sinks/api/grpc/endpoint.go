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
	"github.com/etaques/orb/sinks"
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

func retrieveSinksEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(sinksFilter)
		filter := sinks.Filter{
			OpenTelemetry: req.isOtel,
		}
		sinksInternal, err := svc.ListSinksInternal(ctx, filter)
		if err != nil {
			return sinksInternal, err
		}
		responseStr := sinksRes{}
		for _, sink := range sinksInternal {
			sinkResponse, err := buildSinkResponse(sink)
			if err != nil {

			} else {
				responseStr.sinks = append(responseStr.sinks, sinkResponse)
			}
		}

		return responseStr, err
	}
}

func buildSinkResponse(sink sinks.Sink) (sinkRes, error) {
	tagData, err := json.Marshal(sink.Tags)
	if err != nil {
		return sinkRes{}, err
	}
	configData, err := json.Marshal(sink.Config)
	if err != nil {
		return sinkRes{}, err
	}

	return sinkRes{
		id:          sink.ID,
		name:        sink.Name.String(),
		description: sink.Description,
		tags:        tagData,
		state:       sink.State.String(),
		error:       sink.Error,
		backend:     sink.Backend,
		config:      configData,
	}, nil
}
