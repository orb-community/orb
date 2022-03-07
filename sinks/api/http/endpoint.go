// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
)

func addEndpoint(svc sinks.SinkService) endpoint.Endpoint {
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
			State:       saved.State.String(),
			Error:       saved.Error,
			Backend:     saved.Backend,
			Config:      saved.Config,
			TsCreated:   saved.Created,
			created:     true,
		}

		return res, nil
	}
}

func updateSinkEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateSinkReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		nameID, _ := types.NewIdentifier(req.Name)
		sink := sinks.Sink{
			Name:        nameID,
			ID:          req.id,
			Tags:        req.Tags,
			Config:      req.Config,
			Description: req.Description,
		}

		if _, err := svc.UpdateSink(ctx, req.token, sink); err != nil {
			return nil, err
		}
		res := sinkRes{
			ID:          sink.ID,
			Name:        sink.Name.String(),
			Description: sink.Description,
			Tags:        sink.Tags,
			State:       sink.State.String(),
			Error:       sink.Error,
			Backend:     sink.Backend,
			Config:      sink.Config,
			created:     false,
		}
		return res, nil
	}
}

func listSinksEndpoint(svc sinks.SinkService) endpoint.Endpoint {
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
				State:       sink.State.String(),
				Error:       sink.Error,
				Backend:     sink.Backend,
				Config:      sink.Config,
				TsCreated:   sink.Created,
			}
			res.Sinks = append(res.Sinks, view)
		}
		return res, nil
	}
}

func listBackendsEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listBackendsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		backends, err := svc.ListBackends(ctx, req.token)
		if err != nil {
			return nil, err
		}

		var completeBackends []interface{}
		for _, bk := range backends {
			b, err := svc.ViewBackend(ctx, req.token, bk)
			if err != nil {
				return nil, err
			}
			completeBackends = append(completeBackends, b.Metadata())
		}

		res := sinksBackendsRes{
			Backends: completeBackends,
		}

		return res, nil
	}
}

func viewBackendEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		backend, err := svc.ViewBackend(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := sinksBackendRes{
			backend.Metadata(),
		}

		return res, nil
	}
}

func viewSinkEndpoint(svc sinks.SinkService) endpoint.Endpoint {
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
			State:       sink.State.String(),
			Error:       sink.Error,
			Backend:     sink.Backend,
			Config:      sink.Config,
			TsCreated:   sink.Created,
		}
		return res, err
	}
}

func deleteSinkEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteSinkReq)

		err = req.validate()
		if err != nil {
			return nil, err
		}

		if err := svc.DeleteSink(ctx, req.token, req.id); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}

func validateSinkEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(validateReq)
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

		validated, err := svc.ValidateSink(ctx, req.token, sink)
		if err != nil {
			return nil, err
		}

		res := validateSinkRes{
			ID:          validated.ID,
			Name:        validated.Name.String(),
			Description: validated.Description,
			Tags:        validated.Tags,
			State:       validated.State.String(),
			Error:       validated.Error,
			Backend:     validated.Backend,
			Config:      validated.Config,
		}

		return res, err
	}
}
