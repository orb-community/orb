// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/mainflux/mainflux/things"
	thingsapi "github.com/mainflux/mainflux/things/api/things/http"
	"github.com/ns1labs/orb/fleet"
	api "github.com/ns1labs/orb/fleet/api"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	contentType  = "application/json"
	token        = "token"
	invalidToken = "invalid"
	email        = "user@example.com"
	validJson    = "{\n	\"name\": \"eu-agents\", \n	\"tags\": {\n		\"region\": \"eu\", \n		\"node_type\": \"dns\"\n	}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"
	invalidJson  = "{"
	channelsNum  = 3
)

var (
	agentGroup = fleet.AgentGroup{
		ID:             "",
		MFOwnerID:      "",
		Name:           types.Identifier{},
		Description:    "",
		MFChannelID:    "",
		MatchingAgents: nil,
		Tags:           nil,
		Created:        time.Time{},
	}
	metadata = map[string]interface{}{"meta": "data"}
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
	body        io.Reader
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}
	if tr.token != "" {
		req.Header.Set("Authorization", tr.token)
	}
	if tr.contentType != "" {
		req.Header.Set("Content-Type", tr.contentType)
	}
	return tr.client.Do(req)
}

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

func newServer(svc fleet.Service) *httptest.Server {
	mux := api.MakeHandler(mocktracer.New(), "fleet", svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestCreateAgentGroup(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)
	fleetServer := newServer(fleetService)
	defer fleetServer.Close()

	cases := []struct {
		desc        string
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		{
			desc:        "add a valid agent group",
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusCreated,
			location:    "/agent_groups",
		},
		{
			desc:        "add a valid agent group with invalid token",
			req:         validJson,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/agent_groups",
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      fleetServer.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/agent_groups", fleetServer.URL),
			contentType: tc.contentType,
			token:       tc.auth,
			body:        strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}

}

func TestViewAgentGroup(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)
	fleetServer := newServer(fleetService)
	defer fleetServer.Close()

	validName, _ := types.NewIdentifier("")
	agentGroup.Name = validName
	ag, err := fleetService.CreateAgentGroup(context.Background(), token, agentGroup)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		id          string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"view a existing agent group": {
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			location:    "/agent_groups/{id}",
		},
		"view a non-existing agent group": {
			id:          "9bb1b244-a199-93c2-aa03-28067b431e2c",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
			location:    "/agent_groups/{id}",
		},
		"view a agent group with invalid content type": {
			id:          ag.ID,
			contentType: "application",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/agent_groups/{id}",
		},
		"view a agent group with empty content type": {
			id:          ag.ID,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/agent_groups/{id}",
		},
		"view a agent group with a invalid token": {
			id:          ag.ID,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/agent_groups/{id}",
		},
		"view a agent group with a empty token": {
			id:          ag.ID,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			location:    "/agent_groups/{id}",
		},
	}

	for desc, tc := range cases {
		req := testRequest{
			client:      fleetServer.Client(),
			method:      http.MethodGet,
			url:         fmt.Sprintf("%s/agent_groups/%s", fleetServer.URL, tc.id),
			contentType: tc.contentType,
			token:       tc.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected erro %s", desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
	}

}
