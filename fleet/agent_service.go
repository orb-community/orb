// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"fmt"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/fleet/backend"
	"github.com/ns1labs/orb/fleet/backend/pktvisor"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
)

var (
	ErrCreateAgent = errors.New("failed to create agent")

	// ErrThings indicates failure to communicate with Mainflux Things service.
	// It can be due to networking error or invalid/unauthorized request.
	ErrThings = errors.New("failed to receive response from Things service")

	errCreateThing   = errors.New("failed to create thing")
	errThingNotFound = errors.New("thing not found")
)

func (svc fleetService) ViewAgentByID(ctx context.Context, token string, thingID string) (Agent, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}
	return svc.agentRepo.RetrieveByID(ctx, ownerID, thingID)
}

func (svc fleetService) ViewAgentByIDInternal(ctx context.Context, ownerID string, id string) (Agent, error) {
	return svc.agentRepo.RetrieveByID(ctx, ownerID, id)
}

func (svc fleetService) ListAgents(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Page{}, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return svc.agentRepo.RetrieveAll(ctx, res.GetId(), pm)
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

func (svc fleetService) EditAgent(ctx context.Context, token string, agent Agent) (Agent, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}
	agent.MFOwnerID = ownerID

	err = svc.agentRepo.UpdateAgentByID(ctx, ownerID, agent)
	if err != nil {
		return Agent{}, err
	}

	res, err := svc.agentRepo.RetrieveByID(ctx, ownerID, agent.MFThingID)
	if err != nil {
		return Agent{}, err
	}

	err = svc.agentComms.NotifyAgentGroupMembership(res)
	if err != nil {
		svc.logger.Error("failure during agent group membership comms", zap.Error(err))
	}

	return res, nil
}

func (svc fleetService) ValidateAgent(ctx context.Context, token string, a Agent) (Agent, error) {
	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}

	a.MFOwnerID = mfOwnerID

	return a, nil
}

func (svc fleetService) RemoveAgent(ctx context.Context, token, thingID string) error {
	ownerID, err := svc.identify(token)
	if err != nil {
		return err
	}

	res, err := svc.agentRepo.RetrieveByID(ctx, ownerID, thingID)
	if err != nil {
		return nil
	}

	if errT := svc.mfsdk.DeleteThing(res.MFThingID, token); errT != nil {
		svc.logger.Error("failed to delete thing", zap.Error(errT), zap.String("thing_id", res.MFThingID))
	}

	if errT := svc.mfsdk.DeleteChannel(res.MFChannelID, token); errT != nil {
		svc.logger.Error("failed to delete channel", zap.Error(errT), zap.String("channel_id", res.MFChannelID))
	}

	err = svc.agentRepo.Delete(ctx, ownerID, thingID)
	if err != nil {
		return err
	}

	return nil
}

func (svc fleetService) ListAgentBackends(ctx context.Context, token string) ([]string, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	return backend.GetList(), nil
}

func (svc fleetService) ViewAgentBackend(ctx context.Context, token string, name string) (interface{}, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	if backend.HaveBackend(name) {
		return backend.GetBackend(name).Metadata(), nil
	}
	return nil, errors.ErrNotFound
}

func (svc fleetService) ViewAgentBackendHandler(ctx context.Context, token string, name string) (types.Metadata, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	if backend.HaveBackend(name) {
		bk, err := pktvisor.Handlers()
		if err != nil {
			return nil, err
		}
		return bk, nil
	}
	return nil, errors.ErrNotFound
}

func (svc fleetService) ViewAgentBackendInput(ctx context.Context, token string, name string) (types.Metadata, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	if backend.HaveBackend(name) {
		bk, err := pktvisor.Inputs()
		if err != nil {
			return nil, err
		}
		return bk, nil
	}
	return nil, errors.ErrNotFound
}

func (svc fleetService) ViewAgentBackendTaps(ctx context.Context, token string, name string) ([]BackendTaps, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	if !backend.HaveBackend(name) {
		return nil, errors.ErrNotFound
	}
	metadataList, err := svc.agentRepo.RetrieveAgentMetadataByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	var list []types.Metadata
	for _, mt := range metadataList {
		extractTaps(mt, &list)
	}

	res, err := toBackendTaps(list)
	if err != nil {
		return nil, err
	}
	tapsGroup := groupTaps(res)

	return tapsGroup, nil
}

// Used to get the taps from policy json
func extractTaps(mt map[string]interface{}, list *[]types.Metadata) {
	for k, v := range mt {
		if k == "taps" {
			m, _ := v.(map[string]interface{})
			*list = append(*list, m)
		} else {
			m, _ := v.(map[string]interface{})
			extractTaps(m, list)
		}
	}
}

// Used to cast the map[string]interface for a concrete struct
func toBackendTaps(list []types.Metadata) ([]BackendTaps, error) {
	var bkTaps []BackendTaps
	for _, tc := range list {
		bkTap := BackendTaps{}
		var idx int
		for k, v := range tc {
			bkTap.Name = k
			m, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("Error to group taps")
			}
			for k, v := range m {
				if k == "config" {
					m, ok := v.(map[string]interface{})
					if !ok {
						return nil, errors.New("Error to group taps")
					}
					for k, _ := range m {
						bkTap.ConfigPredefined = append(bkTap.ConfigPredefined, []string{k}...)
					}
				} else {
					bkTap.InputType = k
				}
			}
			idx++
			bkTap.TotalAgents += uint64(idx)
			bkTaps = append(bkTaps, bkTap)
		}
	}
	return bkTaps, nil
}

// Used to aggregate and sumarize the taps and return a slice of BackendTaps
func groupTaps(taps []BackendTaps) []BackendTaps {
	//TODO sort taps before group
	tapsMap := make(map[string]BackendTaps)
	for _, tap := range taps {
		key := key(tap.Name, tap.InputType)
		if v, ok := tapsMap[key]; ok {
			v.ConfigPredefined = append(v.ConfigPredefined, tap.ConfigPredefined...)
			v.TotalAgents += 1
			tapsMap[key] = v
		} else {
			tapsMap[key] = BackendTaps{
				Name:             tap.Name,
				InputType:        tap.InputType,
				ConfigPredefined: tap.ConfigPredefined,
				TotalAgents:      tap.TotalAgents,
			}
		}
	}
	var bkTaps []BackendTaps
	for _, v := range tapsMap {
		bkTaps = append(bkTaps, v)
	}
	return bkTaps
}

func key(name string, inputType string) string {
	return fmt.Sprintf("%s-%s", name, inputType)
}
