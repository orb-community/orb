// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/backend"
	"go.uber.org/zap"
)

func omitSecretInformation(configSvc *sinks.Configuration, inputSink sinks.Sink) (returnSink sinks.Sink, err error) {
	a, err := configSvc.Authentication.OmitInformation("object", inputSink.Config)
	if err != nil {
		return sinks.Sink{}, err
	}
	returnSink = inputSink
	authMeta := a.(types.Metadata)
	returnSink.Config = authMeta
	if inputSink.Format == "yaml" {
		configData, newErr := configSvc.Authentication.ConfigToFormat(inputSink.Format, authMeta)
		if newErr != nil {
			err = newErr
			return
		}
		returnSink.ConfigData = configData.(string)
	}
	return
}

func addEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addReq)
		if err := req.validate(); err != nil {
			svc.GetLogger().Error("got error in validating request", zap.Error(err))
			return nil, err
		}

		nID, err := types.NewIdentifier(req.Name)
		if err != nil {
			svc.GetLogger().Error("got error in creating new identifier", zap.Error(err))
			return nil, err
		}
		var exporterConfig types.Metadata
		var authConfig types.Metadata
		var configSvc *sinks.Configuration
		if len(req.Format) > 0 && req.Format == "yaml" {
			if len(req.ConfigData) > 0 {
				configSvc, exporterConfig, authConfig, err = GetConfigurationAndMetadataFromYaml(req.Backend, req.ConfigData)
				if err != nil {
					svc.GetLogger().Error("got error in parse and validate configuration", zap.Error(err))
					return nil, errors.Wrap(errors.ErrMalformedEntity, err)
				}
			} else {
				svc.GetLogger().Error("got error in parse and validate configuration", zap.Error(err))
				return nil, errors.Wrap(errors.ErrMalformedEntity, errors.New("missing required field when format is sent, config_data must be sent also"))
			}
		} else {
			configSvc, exporterConfig, authConfig, err = GetConfigurationAndMetadataFromMeta(req.Backend, req.Config)
			if err != nil {
				svc.GetLogger().Error("got error in parse and validate configuration", zap.Error(err))
				return nil, errors.Wrap(errors.ErrMalformedEntity, err)
			}
		}
		config := types.Metadata{
			"exporter":                            exporterConfig,
			authentication_type.AuthenticationKey: authConfig,
		}
		sink := sinks.Sink{
			Name:        nID,
			Backend:     req.Backend,
			Config:      config,
			Description: &req.Description,
			Tags:        req.Tags,
			ConfigData:  req.ConfigData,
			Format:      req.Format,
			Created:     time.Now(),
		}
		saved, err := svc.CreateSink(ctx, req.token, sink)
		if err != nil {
			svc.GetLogger().Error("received error on creating sink", zap.Error(err))
			return nil, err
		}

		omittedSink, err := omitSecretInformation(configSvc, saved)
		if err != nil {
			svc.GetLogger().Error("sink was created, but got error in the response build", zap.Error(err))
			return nil, err
		}
		res := sinkRes{
			ID:          saved.ID,
			Name:        saved.Name.String(),
			Description: *saved.Description,
			Tags:        saved.Tags,
			State:       saved.State.String(),
			Error:       saved.Error,
			Backend:     saved.Backend,
			Config:      omittedSink.Config,
			ConfigData:  omittedSink.ConfigData,
			Format:      saved.Format,
			TsCreated:   saved.Created,
			created:     true,
		}

		return res, nil
	}
}

func updateSinkEndpoint(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateSinkReq)
		if req.token == "" {
			return nil, errors.ErrUnauthorizedAccess
		}
		currentSink, err := svc.ViewSink(ctx, req.token, req.id)
		if err != nil {
			svc.GetLogger().Error("could not find sink with id", zap.String("sinkID", req.id), zap.Error(err))
			return nil, err
		}
		if err := req.validate(); err != nil {
			svc.GetLogger().Error("error validating request", zap.Error(err))
			return nil, err
		}
		// Update only the fields in currentSink that are populated in req
		if req.Tags != nil {
			currentSink.Tags = req.Tags
		}
		if req.ConfigData != "" {
			currentSink.ConfigData = req.ConfigData
		}
		if req.Format != "" {
			currentSink.Format = req.Format
		}
		if req.Description != nil {
			currentSink.Description = req.Description
		}
		if req.Name != "" {
			nameID, err := types.NewIdentifier(req.Name)
			if err != nil {
				svc.GetLogger().Error("error on getting new identifier", zap.Error(err))
				return nil, errors.ErrConflict
			}
			currentSink.Name = nameID
		}
		var configSvc *sinks.Configuration
		var exporterConfig types.Metadata
		var authConfig types.Metadata

		// Update the config if either req.Config or req.ConfigData is populated
		if req.Config != nil || req.ConfigData != "" {
			if req.Format == "yaml" {
				configSvc, exporterConfig, authConfig, err = GetConfigurationAndMetadataFromYaml(currentSink.Backend, req.ConfigData)
				if err != nil {
					svc.GetLogger().Error("got error in parse and validate configuration", zap.Error(err))
					return nil, errors.Wrap(errors.ErrMalformedEntity, err)
				}
			} else if req.Config != nil {
				configSvc, exporterConfig, authConfig, err = GetConfigurationAndMetadataFromMeta(currentSink.Backend, req.Config)
				if err != nil {
					svc.GetLogger().Error("got error in parse and validate configuration", zap.Error(err))
					return nil, errors.Wrap(errors.ErrMalformedEntity, err)
				}
			}

			currentSink.Config = types.Metadata{
				"exporter":                            exporterConfig,
				authentication_type.AuthenticationKey: authConfig,
			}
		}

		if err := req.validate(); err != nil {
			svc.GetLogger().Error("error validating request", zap.Error(err))
			return nil, err
		}

		sinkEdited, err := svc.UpdateSink(ctx, req.token, currentSink)
		if err != nil {
			svc.GetLogger().Error("error on updating sink", zap.Error(err))
			return nil, err
		}

		var omittedSink sinks.Sink
		omittedSink, err = omitSecretInformation(configSvc, sinkEdited)
		if err != nil {
			svc.GetLogger().Error("sink was updated, but got error in the response build", zap.Error(err))
			return nil, err
		}

		res := sinkRes{
			ID:          sinkEdited.ID,
			Name:        sinkEdited.Name.String(),
			Description: *sinkEdited.Description,
			Tags:        sinkEdited.Tags,
			State:       sinkEdited.State.String(),
			Error:       sinkEdited.Error,
			Backend:     sinkEdited.Backend,
			Config:      omittedSink.Config,
			ConfigData:  omittedSink.ConfigData,
			Format:      sinkEdited.Format,
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
			reqBackend := backend.GetBackend(sink.Backend)
			reqAuthType, _ := authentication_type.GetAuthType(sink.GetAuthenticationTypeName())
			cfg := sinks.Configuration{
				Exporter:       reqBackend,
				Authentication: reqAuthType,
			}
			responseSink, err := omitSecretInformation(&cfg, sink)
			if err != nil {
				return nil, err
			}
			view := sinkRes{
				ID:         sink.ID,
				Name:       sink.Name.String(),
				Tags:       sink.Tags,
				State:      sink.State.String(),
				Error:      sink.Error,
				Backend:    sink.Backend,
				Config:     responseSink.Config,
				ConfigData: responseSink.ConfigData,
				Format:     sink.Format,
				TsCreated:  sink.Created,
			}
			if sink.Description != nil {
				view.Description = *sink.Description
			}
			res.Sinks = append(res.Sinks, view)
		}
		return res, nil
	}
}

func listAuthenticationTypes(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listBackendsReq)
		if err = req.validate(); err != nil {
			return nil, err
		}

		authtypes, err := svc.ListAuthenticationTypes(ctx, req.token)
		if err != nil {
			return nil, err
		}

		response = sinkAuthTypesRes{
			AuthenticationTypes: authtypes,
		}

		return
	}
}

func viewAuthenticationType(svc sinks.SinkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err = req.validate(); err != nil {
			return nil, err
		}
		authType, err := svc.ViewAuthenticationType(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		response = sinkAuthTypeRes{
			AuthenticationTypes: authType,
		}
		return
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
		reqBackend := backend.GetBackend(sink.Backend)
		reqAuthType, _ := authentication_type.GetAuthType(sink.GetAuthenticationTypeName())
		cfg := sinks.Configuration{
			Exporter:       reqBackend,
			Authentication: reqAuthType,
		}
		responseSink, err := omitSecretInformation(&cfg, sink)
		res := sinkRes{
			ID:          sink.ID,
			Name:        sink.Name.String(),
			Description: *sink.Description,
			Tags:        sink.Tags,
			State:       sink.State.String(),
			Error:       sink.Error,
			Backend:     sink.Backend,
			Config:      responseSink.Config,
			ConfigData:  responseSink.ConfigData,
			Format:      sink.Format,
			TsCreated:   sink.Created,
		}
		if sink.Description != nil {
			res.Description = *sink.Description
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
			Description: &req.Description,
			Tags:        req.Tags,
		}

		validated, err := svc.ValidateSink(ctx, req.token, sink)
		if err != nil {
			return nil, err
		}

		res := validateSinkRes{
			ID:          validated.ID,
			Name:        validated.Name.String(),
			Description: *validated.Description,
			Tags:        validated.Tags,
			State:       validated.State.String(),
			Error:       validated.Error,
			Backend:     validated.Backend,
			Config:      validated.Config,
		}

		return res, err
	}
}
