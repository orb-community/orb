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
	"strings"
	"testing"
	"time"
)

const (
	token        = "token"
	invalidToken = "invalid"
	email        = "user@example.com"
	channelsNum  = 3
	maxNameSize  = 1024
	limit        = 10
	wrongID      = "9bb1b244-a199-93c2-aa03-28067b431e2c"
)

var (
	agentGroup = fleet.AgentGroup{
		ID:          "",
		MFOwnerID:   "",
		Name:        types.Identifier{},
		Description: "",
		MFChannelID: "",
		Tags:        nil,
		Created:     time.Time{},
	}
	invalidName = strings.Repeat("m", maxNameSize+1)
	metadata    = map[string]interface{}{"meta": "data"}
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
	agentComms := flmocks.NewFleetCommService()
	logger, _ := zap.NewDevelopment()
	config := mfsdk.Config{
		BaseURL: url,
	}

	mfsdk := mfsdk.NewSDK(config)
	return fleet.NewFleetService(logger, auth, agentRepo, agentGroupRepo, agentComms, mfsdk)
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
	}

	for desc, tc := range cases {
		_, err := fleetService.CreateAgentGroup(context.Background(), tc.token, tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}

}

func TestViewAgentGroup(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ag, err := createAgentGroup(t, "ue-agent-group", fleetService)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"view a existing agent group": {
			id:    ag.ID,
			token: token,
			err:   nil,
		},
		"view agent group with wrong credentials": {
			id:    ag.ID,
			token: "wrong",
			err:   things.ErrUnauthorizedAccess,
		},
		"view non-existing agent group": {
			id:    "9bb1b244-a199-93c2-aa03-28067b431e2c",
			token: token,
			err:   things.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		_, err := fleetService.ViewAgentGroupByID(context.Background(), tc.token, tc.id)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestListAgentGroup(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	var agents []fleet.AgentGroup
	for i := 0; i < limit; i++ {
		ag, err := createAgentGroup(t, fmt.Sprintf("ue-agent-group-%d", i), fleetService)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		agents = append(agents, ag)
	}

	cases := map[string]struct {
		token string
		pm    fleet.PageMetadata
		size  uint64
		err   error
	}{
		"retrieve a list of agent groups": {
			token: token,
			pm: fleet.PageMetadata{
				Limit:  limit,
				Offset: 0,
			},
			size: limit,
			err:  nil,
		},
		"list half": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: limit / 2,
				Limit:  limit,
			},
			size: limit / 2,
			err:  nil,
		},
		"list last agent group": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: limit - 1,
				Limit:  limit,
			},
			size: 1,
			err:  nil,
		},
		"list empty set": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: limit + 1,
				Limit:  limit,
			},
			size: 0,
			err:  nil,
		},
		"list with zero limit": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: 1,
				Limit:  0,
			},
			size: 0,
			err:  nil,
		},
		"list with wrong credentials": {
			token: "wrong",
			pm: fleet.PageMetadata{
				Offset: 0,
				Limit:  0,
			},
			size: 0,
			err:  fleet.ErrUnauthorizedAccess,
		},
		"list all agent groups sorted by name ascendent": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: 0,
				Limit:  limit,
				Order:  "name",
				Dir:    "asc",
			},
			size: limit,
			err:  nil,
		},
		"list all agent groups sorted by name descendent": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: 0,
				Limit:  limit,
				Order:  "name",
				Dir:    "desc",
			},
			size: limit,
			err:  nil,
		},
	}

	for desc, tc := range cases {
		page, err := fleetService.ListAgentGroups(context.Background(), tc.token, tc.pm)
		size := uint64(len(page.AgentGroups))
		assert.Equal(t, size, tc.size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, size))
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		testSortAgentGroups(t, tc.pm, page.AgentGroups)

	}
}

func TestUpdateAgentGroup(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ag, err := createAgentGroup(t, "ue-agent-group", fleetService)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	matching := types.Metadata{"total": 0, "online": 0}
	wrongAgentGroup := fleet.AgentGroup{ID: wrongID}
	readyOnlyAgentGroup := fleet.AgentGroup{ID: wrongID, MatchingAgents: matching}
	cases := map[string]struct {
		group fleet.AgentGroup
		token string
		err   error
	}{
		"update existing sink": {
			group: ag,
			token: token,
			err:   nil,
		},
		"update group with wrong credentials": {
			group: ag,
			token: invalidToken,
			err:   fleet.ErrUnauthorizedAccess,
		},
		"update a non-existing group": {
			group: wrongAgentGroup,
			token: token,
			err:   fleet.ErrNotFound,
		},
		"update group read only fields": {
			group: readyOnlyAgentGroup,
			token: token,
			err:   errors.ErrUpdateEntity,
		},
	}

	for desc, tc := range cases {
		_, err := fleetService.EditAgentGroup(context.Background(), tc.token, tc.group)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %d got %d", desc, tc.err, err))
	}
}

func TestRemoveAgentGroup(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ag, err := createAgentGroup(t, "ue-agent-group", fleetService)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"remove existing agent group": {
			id:    ag.ID,
			token: token,
			err:   nil,
		},
		"remove agent group with wrong credentials": {
			id:    ag.ID,
			token: "wrong",
			err:   things.ErrUnauthorizedAccess,
		},
		"remove removed agent group": {
			id:    ag.ID,
			token: token,
			err:   nil,
		},
		"remove non-existing thing": {
			id:    wrongID,
			token: token,
			err:   nil,
		},
	}

	for desc, tc := range cases {
		err := fleetService.RemoveAgentGroup(context.Background(), tc.token, tc.id)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func createAgentGroup(t *testing.T, name string, svc fleet.AgentGroupService) (fleet.AgentGroup, error) {
	t.Helper()
	agCopy := agentGroup
	validName, err := types.NewIdentifier(name)
	if err != nil {
		return fleet.AgentGroup{}, err
	}
	agCopy.Name = validName
	ag, err := svc.CreateAgentGroup(context.Background(), token, agentGroup)
	if err != nil {
		return fleet.AgentGroup{}, err
	}
	return ag, nil
}

func testSortAgentGroups(t *testing.T, pm fleet.PageMetadata, ags []fleet.AgentGroup) {
	t.Helper()
	switch pm.Order {
	case "name":
		current := ags[0]
		for _, res := range ags {
			if pm.Dir == "asc" {
				assert.GreaterOrEqual(t, res.Name.String(), current.Name.String())
			}
			if pm.Dir == "desc" {
				assert.GreaterOrEqual(t, current.Name.String(), res.Name.String())
			}
			current = res
		}
	default:
		break
	}
}
