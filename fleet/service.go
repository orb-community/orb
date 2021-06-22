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
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
	"time"
)

var (
	ErrCreateSelector = errors.New("failed to create selector")

	ErrCreateAgent = errors.New("failed to create agent")

	// ErrThings indicates failure to communicate with Mainflux Things service.
	// It can be due to networking error or invalid/unauthorized request.
	ErrThings = errors.New("failed to receive response from Things service")

	errCreateThing   = errors.New("failed to create thing")
	errThingNotFound = errors.New("thing not found")
)

type Service interface {
	AgentService
	SelectorService
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total    uint64
	Offset   uint64         `json:"offset,omitempty"`
	Limit    uint64         `json:"limit,omitempty"`
	Name     string         `json:"name,omitempty"`
	Order    string         `json:"order,omitempty"`
	Dir      string         `json:"dir,omitempty"`
	Metadata types.Metadata `json:"metadata,omitempty"`
}

var _ Service = (*fleetService)(nil)

type fleetService struct {
	logger *zap.Logger
	// for AuthN/AuthZ
	auth mainflux.AuthServiceClient
	// for Thing manipulation
	mfsdk mfsdk.SDK
	// Agents and Selectors
	agentRepo    AgentRepository
	selectorRepo SelectorRepository
}

func (svc fleetService) ListAgents(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Page{}, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return svc.agentRepo.RetrieveAll(ctx, res.GetId(), pm)
}

func (svc fleetService) identify(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
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
			return mfsdk.Thing{}, errors.Wrap(errThingNotFound, errors.ErrNotFound)
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

	md := map[string]interface{}{"type": "orb_agent"}

	// create new Thing
	mfThing, err := svc.thing(token, "", a.Name.String(), md)
	if err != nil {
		return Agent{}, errors.Wrap(ErrCreateAgent, err)
	}

	a.MFThingID = mfThing.ID
	a.MFKeyID = mfThing.Key

	// create main Agent RPC Channel
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

	// RPC Channel to Agent
	err = svc.mfsdk.Connect(mfsdk.ConnectionIDs{
		ChannelIDs: []string{mfChannelID},
		ThingIDs:   []string{mfThing.ID},
	}, token)
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

func NewFleetService(logger *zap.Logger, auth mainflux.AuthServiceClient, agentRepo AgentRepository, selectorRepo SelectorRepository, mfsdk mfsdk.SDK) Service {
	return &fleetService{
		logger:       logger,
		auth:         auth,
		agentRepo:    agentRepo,
		selectorRepo: selectorRepo,
		mfsdk:        mfsdk,
	}
}
