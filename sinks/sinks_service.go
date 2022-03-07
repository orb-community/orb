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
	"github.com/ns1labs/orb/sinks/backend"
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

	id, err := svc.sinkRepo.Save(ctx, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}
	sink.ID = id
	return sink, nil
}

func (svc sinkService) UpdateSink(ctx context.Context, token string, sink Sink) (Sink, error) {
	skOwnerID, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}

	if sink.Backend != "" || sink.Error != "" {
		return Sink{}, errors.ErrUpdateEntity
	}

	sink.MFOwnerID = skOwnerID
	err = svc.sinkRepo.Update(ctx, sink)
	if err != nil {
		return Sink{}, err
	}
	return sink, nil
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
	return res, nil
}

func (svc sinkService) ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	res, err := svc.identify(token)
	if err != nil {
		return Page{}, err
	}

	return svc.sinkRepo.RetrieveAll(ctx, res, pm)
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
