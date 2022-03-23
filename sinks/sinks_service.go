// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks/backend"
	"github.com/ns1labs/orb/sinks/prometheus"
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

	if err := validateConfig(sink.Config); err != nil {
		sink.State = Error
		sink.Error = err.Error()
	} else {
		sink.State = Active
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

	err = validateConfig(sink.Config)
	if err != nil {
		sink.State = Error
		sink.Error = err.Error()
	} else {
		sink.State = Active
	}

	err = svc.sinkRepo.Update(ctx, sink)
	if err != nil {
		return Sink{}, err
	}

	sinkEdited, err := svc.sinkRepo.RetrieveById(ctx, sink.ID)
	if err != nil {
		return Sink{}, err
	}

	return sinkEdited, nil
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
	if !backend.HaveBackend(sink.Backend) {
		return ErrInvalidBackend
	}
	return nil
}

func validateConfig(config types.Metadata) error {
	var cfgRepo SinkConfig
	data, _ := json.Marshal(config)
	err := json.Unmarshal(data, &cfgRepo)

	cfg := prometheus.NewConfig(
		prometheus.WriteURLOption(cfgRepo.Url),
	)

	promClient, err := prometheus.NewClient(cfg)
	if err != nil {
		return err
	}

	var dpFlag prometheus.Datapoint
	dpFlag.Value = 1
	dpFlag.Timestamp = time.Now()

	var labelsListFlag []prometheus.Label
	label := prometheus.Label{
		Name:  "__name__",
		Value: "login",
	}

	labelsListFlag = append(labelsListFlag, label)

	tsList := prometheus.TSList{
		struct {
			Labels    []prometheus.Label
			Datapoint prometheus.Datapoint
		}{Labels: labelsListFlag, Datapoint: dpFlag},
	}

	var headers = make(map[string]string)
	headers["Authorization"] = encodeBase64(cfgRepo.User, cfgRepo.Password)
	_, writeErr := promClient.WriteTimeSeries(context.Background(), tsList,
		prometheus.WriteOptions{Headers: headers})
	if err := error(writeErr); err != nil {
		return errors.Wrap(errors.New(err.Error()), err)
	}

	return nil
}

func encodeBase64(user string, password string) string {
	sEnc := b64.URLEncoding.EncodeToString([]byte(user + ":" + password))
	return fmt.Sprintf("Basic %s", sEnc)
}