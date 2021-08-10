// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/mainflux/mainflux/things"
	thingsapi "github.com/mainflux/mainflux/things/api/things/http"
	"github.com/ns1labs/orb/fleet"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

const (
	token       = "token"
	invalidToken       = ""
	email       = "user@example.com"
	channelsNum = 3
)

var (
	metadata = map[string]interface{}{"meta": "data"}
)

func generateChannels() map[string]things.Channel {
	channels := make(map[string]things.Channel, channelsNum)
	for i := 0; i < channelsNum; i++ {
		id := strconv.Itoa(i + 1)
		channels[id] = things.Channel{
			ID:       id,
			Owner:    email,
			Metadata: metadata,
		}
	}
	return channels
}

func newThingsService(auth mainflux.AuthServiceClient) things.Service {
	return flmocks.NewThingsService(map[string]things.Thing{}, generateChannels(), auth)
}

func newThingsServer(svc things.Service) *httptest.Server {
	mux := thingsapi.MakeHandler(mocktracer.New(), svc)
	return httptest.NewServer(mux)
}

func newService(auth mainflux.AuthServiceClient, url string) fleet.Service {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()
	var logger *zap.Logger
	config := mfsdk.Config{
		BaseURL: url,
	}

	mfsdk := mfsdk.NewSDK(config)
	return fleet.NewFleetService(logger, auth, agentRepo, agentGroupRepo, nil, mfsdk)
}

func TestCreateAgentGroup(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	nameID, err := types.NewIdentifier("eu-agents")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	validAgent := fleet.AgentGroup{
		MFOwnerID:   ownerID.String(),
		Name:        nameID,
		Description: "An example agent group representing european dns nodes",
		Tags:        make(map[string]string),
		Created:     time.Time{},
	}

	validAgent.Tags = map[string]string{
		"region":    "eu",
		"node_type": "dns",
	}

	cases := map[string]struct {
		agent fleet.AgentGroup
		token string
		err   error
	}{
		"add a valid agent group": {
			agent: validAgent,
			token: token,
			err:   nil,
		},
		"add a agent group with invalid token": {
			agent: validAgent,
			token: invalidToken,
			err:   fleet.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		_, err := fleetService.CreateAgentGroup(context.Background(), tc.token, tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}

}

func TestValidateAgentGroup(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	nameID, err := types.NewIdentifier("eu-agents")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	validAgent := fleet.AgentGroup{
		MFOwnerID:   ownerID.String(),
		Name:        nameID,
		Description: "An example agent group representing european dns nodes",
		Tags:        make(map[string]string),
		Created:     time.Time{},
	}

	validAgent.Tags = map[string]string{
		"region":    "eu",
		"node_type": "dns",
	}

	cases := map[string]struct {
		agent fleet.AgentGroup
		token string
		err   error
	}{
		"validate a valid agent group": {
			agent: validAgent,
			token: token,
			err:   nil,
		},
		"validate a agent group with a invalid token": {
			agent: validAgent,
			token: invalidToken,
			err:   fleet.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		_, err := fleetService.ValidateAgentGroup(context.Background(), tc.token, tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}

}