// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
	"time"
)

type Service interface {
	AgentService
	AgentGroupService
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
	Tags     types.Tags     `json:"tags,omitempty"`
}

var _ Service = (*fleetService)(nil)

type fleetService struct {
	logger *zap.Logger
	// for AuthN/AuthZ
	auth mainflux.AuthServiceClient
	// for Thing manipulation
	mfsdk mfsdk.SDK
	// Agents and Agent Groups
	agentRepo            AgentRepository
	agentGroupRepository AgentGroupRepository
	// Agent Comms
	agentComms AgentCommsService
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

func NewFleetService(logger *zap.Logger, auth mainflux.AuthServiceClient, agentRepo AgentRepository, agentGroupRepository AgentGroupRepository, agentComms AgentCommsService, mfsdk mfsdk.SDK, db *sqlx.DB) Service {

	return &fleetService{
		logger:               logger,
		auth:                 auth,
		agentRepo:            agentRepo,
		agentGroupRepository: agentGroupRepository,
		agentComms:           agentComms,
		mfsdk:                mfsdk,
	}

}
