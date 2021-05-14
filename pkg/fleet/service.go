// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/pkg/errors"
	"time"
)

var (
	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrMalformedEntity indicates malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = errors.New("failed to scan metadata")

	ErrCreateSelector = errors.New("failed to create selector")

	ErrCreateAgent = errors.New("failed to create agent")

	// ErrThings indicates failure to communicate with Mainflux Things service.
	// It can be due to networking error or invalid/unauthorized request.
	ErrThings = errors.New("failed to receive response from Things service")

	errCreateThing   = errors.New("failed to create thing")
	errThingNotFound = errors.New("thing not found")
)

// A flat kv pair object
type Tags map[string]interface{}

// Maybe a full object hierarchy
type Metadata map[string]interface{}

type Service interface {
	AgentService
	SelectorService
}

var _ Service = (*fleetService)(nil)

type fleetService struct {
	// for AuthN/AuthZ
	auth mainflux.AuthServiceClient
	// for Thing manipulation
	mfsdk mfsdk.SDK
	// Agents and Selectors
	agentRepo    AgentRepository
	selectorRepo SelectorRepository
}

func (svc fleetService) identify(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(ErrUnauthorizedAccess, err)
	}

	return res.GetId(), nil
}

func (svc fleetService) CreateSelector(ctx context.Context, token string, s Selector) (Selector, error) {
	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Selector{}, err
	}

	s.MFOwnerID = mfOwnerID

	err = svc.selectorRepo.Save(ctx, s)
	if err != nil {
		return Selector{}, errors.Wrap(ErrCreateSelector, err)
	}

	return s, nil
}

// Method thing retrieves Mainflux Thing creating one if an empty ID is passed.
func (svc fleetService) thing(token, id string, name string, md map[string]interface{}) (mfsdk.Thing, error) {
	thingID := id
	var err error

	if id == "" {
		thingID, err = svc.mfsdk.CreateThing(mfsdk.Thing{Name: name, Metadata: md}, token)
		if err != nil {
			return mfsdk.Thing{}, errors.Wrap(errCreateThing, err)
		}
	}

	thing, err := svc.mfsdk.Thing(thingID, token)
	if err != nil {
		if errors.Contains(err, mfsdk.ErrFailedFetch) {
			return mfsdk.Thing{}, errors.Wrap(errThingNotFound, ErrNotFound)
		}

		if id != "" {
			if errT := svc.mfsdk.DeleteThing(thingID, token); errT != nil {
				err = errors.Wrap(err, errT)
			}
		}

		return mfsdk.Thing{}, errors.Wrap(ErrThings, err)
	}

	return thing, nil
}

func (svc fleetService) CreateAgent(ctx context.Context, token string, a Agent) (Agent, error) {
	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}

	a.MFOwnerID = mfOwnerID

	md := map[string]interface{}{"type": "orb-agent"}

	// create new Thing
	mfThing, err := svc.thing(token, "", a.Name.String(), md)
	if err != nil {
		return Agent{}, errors.Wrap(ErrCreateAgent, err)
	}

	a.MFThingID = mfThing.ID
	a.MFKeyID = mfThing.Key

	// create main Channel
	mfChannelID, err := svc.mfsdk.CreateChannel(mfsdk.Channel{
		Name:     a.Name.String(),
		Metadata: md,
	}, token)
	if err != nil {
		if errT := svc.mfsdk.DeleteThing(mfThing.ID, token); errT != nil {
			err = errors.Wrap(err, errT)
		}
		return Agent{}, errors.Wrap(ErrCreateAgent, err)
	}

	a.MFChannelID = mfChannelID

	err = svc.agentRepo.Save(ctx, a)
	if err != nil {
		if errT := svc.mfsdk.DeleteThing(mfThing.ID, token); errT != nil {
			err = errors.Wrap(err, errT)
			// fall through
		}
		if errT := svc.mfsdk.DeleteChannel(mfChannelID, token); errT != nil {
			err = errors.Wrap(err, errT)
		}
		return Agent{}, errors.Wrap(ErrCreateAgent, err)
	}

	return a, nil
}

func NewFleetService(auth mainflux.AuthServiceClient, agentRepo AgentRepository, selectorRepo SelectorRepository, mfsdk mfsdk.SDK) Service {
	return &fleetService{
		auth:         auth,
		agentRepo:    agentRepo,
		selectorRepo: selectorRepo,
		mfsdk:        mfsdk,
	}
}
