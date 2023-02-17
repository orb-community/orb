// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"context"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks/backend"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrCreateSink                 = errors.New("failed to create Sink")
	ErrConflictSink               = errors.New("entity already exists")
	ErrUnsupportedContentTypeSink = errors.New("unsupported content type")
	ErrValidateSink               = errors.New("failed to validate Sink")
)

func (svc sinkService) CreateSink(ctx context.Context, token string, sink Sink) (Sink, error) {

	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}

	sink.MFOwnerID = mfOwnerID

	err = validateBackend(&sink)
	if err != nil {
		return Sink{}, err
	}

	// Validate remote_host
	_, err = url.ParseRequestURI(sink.Config["remote_host"].(string))
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}
	err, _ = svc.validateConfig(ctx, sink.Config)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}

	// encrypt data for the password
	sink, err = svc.encryptMetadata(sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}

	//// add default values
	defaultMetadata := make(types.Metadata, 1)
	defaultMetadata["opentelemetry"] = "enabled"
	sink.Config.Merge(defaultMetadata)

	id, err := svc.sinkRepo.Save(ctx, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}
	sink.ID = id

	// After creating, decrypt Metadata to send correct information to Redis
	sink, err = svc.decryptMetadata(sink)

	return sink, nil
}

func (svc sinkService) encryptMetadata(sink Sink) (Sink, error) {
	var err error
	sink.Config.FilterMap(func(key string) bool {
		return key == backend.ConfigFeatureTypePassword
	}, func(key string, value interface{}) (string, interface{}) {
		newValue, err2 := svc.passwordService.EncodePassword(value.(string))
		if err2 != nil {
			err = err2
			return key, value
		}
		return key, newValue
	})
	return sink, err
}

func (svc sinkService) decryptMetadata(sink Sink) (Sink, error) {
	var err error
	sink.Config.FilterMap(func(key string) bool {
		return key == backend.ConfigFeatureTypePassword
	}, func(key string, value interface{}) (string, interface{}) {
		newValue, err2 := svc.passwordService.DecodePassword(value.(string))
		if err2 != nil {
			err = err2
			return key, value
		}
		return key, newValue
	})
	return sink, err
}

func (svc sinkService) UpdateSink(ctx context.Context, token string, sink Sink) (Sink, error) {
	skOwnerID, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}

	var currentSink Sink
	currentSink, err = svc.sinkRepo.RetrieveById(ctx, sink.ID)
	if err != nil {
		return Sink{}, err
	}

	if sink.Config == nil {
		sink.Config = currentSink.Config
	} else {
		// Validate remote_host
		_, err := url.ParseRequestURI(sink.Config["remote_host"].(string))
		if err != nil {
			return Sink{}, errors.Wrap(ErrUpdateEntity, err)
		}
		err, ok := svc.validateConfig(ctx, sink.Config)
		if err != nil {
			return Sink{}, errors.Wrap(ErrValidateSink, err)
		}
		if ok {
			// This will keep the previous tags
			currentSink.Config.Merge(sink.Config)
			sink.Config = currentSink.Config
		}
	}

	if sink.Tags == nil {
		sink.Tags = currentSink.Tags
	}

	if sink.Description == nil {
		sink.Description = currentSink.Description
	}

	if newName := sink.Name.String(); newName == "" {
		sink.Name = currentSink.Name
	}

	if sink.Backend != "" || sink.Error != "" {
		return Sink{}, errors.ErrUpdateEntity
	}
	sink.MFOwnerID = skOwnerID
	sink, err = svc.encryptMetadata(sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}
	err = svc.sinkRepo.Update(ctx, sink)
	if err != nil {
		return Sink{}, err
	}
	sinkEdited, err := svc.sinkRepo.RetrieveById(ctx, sink.ID)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}
	sinkEdited, err = svc.decryptMetadata(sinkEdited)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}

	return sinkEdited, nil
}

// validateConfig Validate on BackEnd to check whether the combination for host/username/password is valid
func (svc sinkService) validateConfig(ctx context.Context, config types.Metadata) (err error, ok bool) {
	requestCtx, requestCancel := context.WithTimeout(ctx, 2*time.Second)
	defer requestCancel()
	configUrl := config["remote_host"].(string)
	configUsername := config["username"].(string)
	configPassword := config["password"].(string)
	err, ok = svc.requestAuth(requestCtx, configUrl, configUsername, configPassword)
	if err != nil {
		svc.logger.Error("got error during validation username configuration")
		return
	}
	return
}

func (svc sinkService) requestAuth(ctx context.Context, url, username, password string) (err error, ok bool) {
	client := &http.Client{
		Timeout: time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err, false
	}
	req.SetBasicAuth(username, password)
	response, err := client.Do(req)
	if err != nil {
		return err, false
	}
	svc.logger.Info("response code", zap.String("HTTP Status Code", response.Status))
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	return nil, true
}

func (svc sinkService) ListBackends(ctx context.Context, token string) ([]string, error) {
	_, err := svc.identify(token)
	if err != nil {
		return []string{}, err
	}
	return backend.GetList(), nil
}

func (svc sinkService) ViewBackend(ctx context.Context, token string, key string) (backend.Backend, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	res := backend.GetBackend(key)
	if res == nil {
		return nil, errors.Wrap(errors.ErrNotFound, err)
	}
	return res, nil
}

func (svc sinkService) ViewSink(ctx context.Context, token string, key string) (Sink, error) {
	_, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}
	res, err := svc.sinkRepo.RetrieveById(ctx, key)
	if err != nil {
		return Sink{}, errors.Wrap(errors.ErrNotFound, err)
	}
	return res, nil
}

func (svc sinkService) ViewSinkInternal(ctx context.Context, ownerID string, key string) (Sink, error) {
	res, err := svc.sinkRepo.RetrieveByOwnerAndId(ctx, ownerID, key)
	if err != nil {
		return Sink{}, errors.Wrap(errors.ErrNotFound, err)
	}
	res, err = svc.decryptMetadata(res)
	if err != nil {
		return Sink{}, errors.Wrap(errors.ErrViewEntity, err)
	}
	return res, nil
}

func (svc sinkService) ListSinksInternal(ctx context.Context, filter Filter) (sinks []Sink, err error) {
	sinks, err = svc.sinkRepo.SearchAllSinks(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNotFound, err)
	}
	for _, sink := range sinks {
		sink, err = svc.decryptMetadata(sink)
		if err != nil {
			return nil, errors.Wrap(errors.ErrViewEntity, err)
		}
	}

	return
}

func (svc sinkService) ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	res, err := svc.identify(token)
	if err != nil {
		return Page{}, err
	}

	return svc.sinkRepo.RetrieveAllByOwnerID(ctx, res, pm)
}

func (svc sinkService) DeleteSink(ctx context.Context, token string, id string) error {
	res, err := svc.identify(token)
	if err != nil {
		return err
	}

	return svc.sinkRepo.Remove(ctx, res, id)
}

func (svc sinkService) ValidateSink(ctx context.Context, token string, sink Sink) (Sink, error) {

	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}

	sink.MFOwnerID = mfOwnerID

	err = validateBackend(&sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrValidateSink, err)
	}

	return sink, nil
}

func (svc sinkService) ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state State) error {
	return svc.sinkRepo.UpdateSinkState(ctx, sinkID, msg, ownerID, state)
}

func validateBackend(sink *Sink) error {
	if backend.HaveBackend(sink.Backend) {
		sink.State = Unknown
	} else {
		return ErrInvalidBackend
	}
	return nil
}
