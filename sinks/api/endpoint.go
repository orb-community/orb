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
	"github.com/ns1labs/orb/sinks"
)

func addEndpoint(svc sinks.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			return nil, err
		}

		sink := sinks.Sink{
			Name:        nID,
			Backend:     req.Backend,
			Config:      req.Config,
			Description: req.Description,
			Tags:        req.Tags,
		}
		saved, err := svc.CreateSink(ctx, req.token, sink)
		if err != nil {
			return nil, err
		}

		res := sinkRes{
			ID:          saved.ID,
			Name:        saved.Name.String(),
			Description: saved.Description,
			Tags:        saved.Tags,
			Backend:     saved.Backend,
			Config:      saved.Config,
			TsCreated:   saved.Created,
		}

		return res, nil
	}
}

func listSinksEndpoint(svc sinks.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listResourcesReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListSinks(ctx, req.token, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		res := sinksPagesRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
				Order:  page.Order,
				Dir:    page.Dir,
			},
			Sinks: []sinkRes{},
		}

		for _, sink := range page.Sinks {
			view := sinkRes{
				ID:          sink.ID,
				Name:        sink.Name.String(),
				Description: sink.Description,
				Tags:        sink.Tags,
				Backend:     sink.Backend,
				Config:      sink.Config,
				TsCreated:   sink.Created,
			}
			res.Sinks = append(res.Sinks, view)
		}
		return res, nil
	}
}

func listBackendsEndpoint(svc sinks.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listBackendsReq)

		backends, err := svc.ListBackends(ctx, req.token)
		if err != nil {
			return nil, err
		}

		var completeBackends []interface{}
		for _, bk := range backends {
			b, err := svc.GetBackend(ctx, req.token, bk)
			if err != nil {
				return nil, err
			}
			completeBackends = append(completeBackends, b)
		}

		res := sinksBackendsRes{
			Backends: completeBackends,
		}

		return res, nil
	}
}

func viewBackendEndpoint(svc sinks.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)

		backend, err := svc.GetBackend(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := sinksBackendRes{
			Backend:     backend.GetName(),
			Description: backend.GetDescription(),
			Config:      backend.GetConfig(),
		}

		return res, nil
	}
}

func viewSinkEndpoint(svc sinks.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		sink, err := svc.ViewSink(ctx, req.token, req.id)
		if err != nil {
			return sink, err
		}

		res := sinkRes{
			ID:          sink.ID,
			Name:        sink.Name.String(),
			Description: sink.Description,
			Tags:        sink.Tags,
			Backend:     sink.Backend,
			Config:      sink.Config,
			TsCreated:   sink.Created,
		}
		return res, err
	}
}
