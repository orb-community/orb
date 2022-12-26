// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/logger"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/mainflux/mainflux/things"
	thingsapi "github.com/mainflux/mainflux/things/api/things/http"
	"github.com/ns1labs/orb/fleet"
	http2 "github.com/ns1labs/orb/fleet/api/http"
	"github.com/ns1labs/orb/fleet/backend/pktvisor"
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
	contentType       = "application/json"
	token             = "token"
	invalidToken      = "invalid"
	email             = "user@example.com"
	validJson         = "{\n	\"name\": \"eu-agents\", \n	\"tags\": {\n		\"region\": \"eu\", \n		\"node_type\": \"dns\"\n	}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"
	conflictValidJson = "{\n	\"name\": \"eu-agents-conflict\", \n	\"tags\": {\n		\"region\": \"eu\", \n		\"node_type\": \"dns\"\n	}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"
	invalidJson       = "{"
	wrongID           = "9bb1b244-a199-93c2-aa03-28067b431e2c"
	maxNameSize       = 1024
	channelsNum       = 3
	limit             = 10
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
	agent = fleet.Agent{
		Name:          types.Identifier{},
		MFOwnerID:     "",
		MFThingID:     "",
		MFKeyID:       "",
		MFChannelID:   "",
		Created:       time.Time{},
		OrbTags:       nil,
		AgentTags:     nil,
		AgentMetadata: nil,
		State:         0,
		LastHBData:    nil,
		LastHB:        time.Time{},
	}
	metadata    = map[string]interface{}{"type": "orb_agent"}
	tags        = types.Tags{"region": "us", "node_type": "dns"}
	invalidName = strings.Repeat("m", maxNameSize+1)
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
	body        io.Reader
}

type clientServer struct {
	service fleet.Service
	server  *httptest.Server
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
	log := logger.NewMock()
	mux := thingsapi.MakeHandler(mocktracer.New(), svc, log)
	return httptest.NewServer(mux)
}

func newService(auth mainflux.AuthServiceClient, url string) fleet.Service {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()
	agentComms := flmocks.NewFleetCommService(agentRepo, agentGroupRepo)
	logger, _ := zap.NewDevelopment()
	config := mfsdk.Config{
		ThingsURL: url,
	}

	mfsdk := mfsdk.NewSDK(config)
	pktvisor.Register(auth, agentRepo)
	aDone := make(chan bool)
	return fleet.NewFleetService(logger, auth, agentRepo, agentGroupRepo, agentComms, mfsdk, aDone)
}

func newServer(svc fleet.Service) *httptest.Server {
	mux := http2.MakeHandler(mocktracer.New(), "fleet", svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestCreateAgentGroup(t *testing.T) {
	cli := newClientServer(t)
	defer cli.server.Close()

	var missingTagsJson = "{\n	\"name\": \"group\", \n	\"tags\": {}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"
	var invalidNameJson = "{\n	\"name\": \"g\", \n	\"tags\": {\n		\"region\": \"eu\", \n		\"node_type\": \"dns\"\n	}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"

	// Conflict scenario
	createAgentGroup(t, "eu-agents-conflict", &cli)

	cases := map[string]struct {
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"add a valid agent group": {
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusCreated,
			location:    "/agent_groups",
		},
		"add a duplicated agent group": {
			req:         conflictValidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusConflict,
			location:    "/agent_groups",
		},
		"add a valid agent group with invalid token": {
			req:         validJson,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/agent_groups",
		},
		"add a agent group with a invalid json": {
			req:         invalidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agent_groups",
		},
		"add a agent group without a content type": {
			req:         validJson,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/agent_groups",
		},
		"add a agent group without tags": {
			req:         missingTagsJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agent_groups",
		},
		"add a agent group with invalid name": {
			req:         invalidNameJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agent_groups",
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/agent_groups", cli.server.URL),
				contentType: tc.contentType,
				token:       fmt.Sprintf("Bearer %s", tc.auth),
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestViewAgentGroup(t *testing.T) {
	cli := newClientServer(t)

	ag, err := createAgentGroup(t, "ue-agent-group", &cli)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		id       string
		auth     string
		status   int
		location string
	}{
		"view a existing agent group": {
			id:     ag.ID,
			auth:   token,
			status: http.StatusOK,
		},
		"view a non-existing agent group": {
			id:     wrongID,
			auth:   token,
			status: http.StatusNotFound,
		},
		"view a agent group with a invalid token": {
			id:     ag.ID,
			auth:   invalidToken,
			status: http.StatusUnauthorized,
		},
		"view a agent group with a empty token": {
			id:     ag.ID,
			auth:   "",
			status: http.StatusUnauthorized,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/agent_groups/%s", cli.server.URL, tc.id),
				token:  fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected erro %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestListAgentGroup(t *testing.T) {
	cli := newClientServer(t)

	var data []agentGroupRes
	for i := 0; i < limit; i++ {
		ag, err := createAgentGroup(t, fmt.Sprintf("ue-agent-group-%d", i), &cli)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		data = append(data, agentGroupRes{
			ID:             ag.ID,
			Name:           ag.Name.String(),
			Description:    ag.Description,
			Tags:           ag.Tags,
			TsCreated:      ag.Created,
			MatchingAgents: nil,
		})
	}

	cases := map[string]struct {
		auth   string
		status int
		url    string
		res    []agentGroupRes
		total  uint64
	}{
		"retrieve a list of agent groups": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, limit),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agent group with empty token": {
			auth:   "",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 1),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid token": {
			auth:   invalidToken,
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 1),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid dir": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=name&dir=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with negative offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", -1, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with negative limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, -5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with zero limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 0),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agent group without offset": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?limit=%d", limit),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agent group without limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d", 1),
			res:    data[1:limit],
			total:  limit - 1,
		},
		"get a list of agent group with limit greater than max": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 110),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with default URL": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s", ""),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agent group with invalid number of params": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=4&limit=4&limit=5&offset=5"),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=e&limit=5"),
			res:    nil,
			total:  0,
		},
		"get a list of agent group with invalid limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=5&limit=e"),
			res:    nil,
			total:  0,
		},
		"get a list of agent group filtering with invalid name": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&name=%s", 0, 5, invalidName),
			res:    nil,
			total:  0,
		},
		"get a list of agent group sorted with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=wrong&dir=desc", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agent group sorted with invalid direction": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=name&dir=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodGet,
				url:         fmt.Sprintf(fmt.Sprintf("%s/agent_groups%s", cli.server.URL, tc.url)),
				contentType: contentType,
				token:       fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var body agentGroupsPageRes
			json.NewDecoder(res.Body).Decode(&body)
			total := uint64(len(body.AgentGroups))
			assert.Equal(t, res.StatusCode, tc.status, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
			assert.Equal(t, total, tc.total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.total, total))
		})
	}

}

func TestUpdateAgentGroup(t *testing.T) {
	cli := newClientServer(t)

	ag, err := createAgentGroup(t, "ue-agent-group", &cli)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		req         string
		id          string
		contentType string
		auth        string
		status      int
	}{
		"update existing agent group": {
			req: toJSON(updateAgentGroupReq{
				Name:        ag.Name.String(),
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
		},
		"update agent group with a empty json request": {
			req:         "{}",
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update agent group with a invalid id": {
			req: toJSON(updateAgentGroupReq{
				Name:        ag.Name.String(),
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          "invalid",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"update non-existing agent group": {
			req: toJSON(updateAgentGroupReq{
				Name:        ag.Name.String(),
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          wrongID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"update agent group with invalid user token": {
			req: toJSON(updateAgentGroupReq{
				Name:        ag.Name.String(),
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          ag.ID,
			contentType: contentType,
			auth:        "invalid",
			status:      http.StatusUnauthorized,
		},
		"update agent group with empty user token": {
			req: toJSON(updateAgentGroupReq{
				Name:        ag.Name.String(),
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          ag.ID,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
		},
		"update agent group with invalid content type": {
			req: toJSON(updateAgentGroupReq{
				Name:        ag.Name.String(),
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          ag.ID,
			contentType: "invalid",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"update agent group without content type": {
			req: toJSON(updateAgentGroupReq{
				Name:        ag.Name.String(),
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          ag.ID,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"update agent group with a empty request": {
			req:         "",
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update agent group with a invalid data format": {
			req:         invalidJson,
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update agent group with different owner": {
			req:         invalidJson,
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"add a agent group with invalid name": {
			req: toJSON(updateAgentGroupReq{
				Name:        "g",
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update existing agent group without tags": {
			req: toJSON(updateAgentGroupReq{
				Name:        ag.Name.String(),
				Description: ag.Description,
				Tags:        map[string]string{},
			}),
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update existing agent group without name": {
			req: toJSON(updateAgentGroupReq{
				Description: ag.Description,
				Tags:        ag.Tags,
			}),
			id:          ag.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPut,
				url:         fmt.Sprintf("%s/agent_groups/%s", cli.server.URL, tc.id),
				contentType: tc.contentType,
				token:       fmt.Sprintf("Bearer %s", tc.auth),
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			require.Nil(t, err, "%s: unexpected error: %s", desc, err)
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestDeleteAgentGroup(t *testing.T) {

	cli := newClientServer(t)

	_, err := createAgent(t, "agent", &cli)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	ag, err := createAgentGroup(t, "ue-agent-group", &cli)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		id     string
		auth   string
		status int
	}{
		"delete existing agent group": {
			id:     ag.ID,
			auth:   token,
			status: http.StatusNoContent,
		},
		"delete non-existent agent group": {
			id:     wrongID,
			auth:   token,
			status: http.StatusNotFound,
		},
		"delete agent group with invalid token": {
			id:     ag.ID,
			auth:   invalidToken,
			status: http.StatusUnauthorized,
		},
		"delete agent group with empty token": {
			id:     ag.ID,
			auth:   "",
			status: http.StatusUnauthorized,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodDelete,
				contentType: contentType,
				url:         fmt.Sprintf("%s/agent_groups/%s", cli.server.URL, tc.id),
				token:       fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestValidateAgentGroup(t *testing.T) {
	cli := newClientServer(t)
	defer cli.server.Close()

	var (
		invalidValueTag  = "{\n \"name\": \"eu-agents\", \n    \"tags\": {\n       \"invalidTag\", \n      \"node_type\": \"dns\"\n    }, \n   \"description\": \"An example agent group representing european dns nodes\", \n \"validate_only\": false \n}"
		invalidValueName = "{\n \"name\": \",,AGENT 6,\", \n	\"tags\": {\n		\"region\": \"eu\", \n		\"node_type\": \"dns\"\n	}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"
		invalidField     = "{\n \"nname\": \",,AGENT 6,\", \n	\"tags\": {\n		\"region\": \"eu\", \n		\"node_type\": \"dns\"\n	}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"
		invalidNameJson  = "{\n	\"name\": \"g\", \n	\"tags\": {\n		\"region\": \"eu\", \n		\"node_type\": \"dns\"\n	}, \n	\"description\": \"An example agent group representing european dns nodes\", \n	\"validate_only\": false \n}"
	)
	cases := map[string]struct {
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"validate a valid agent group": {
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			location:    "/agent_groups/validate",
		},

		"validate a agent group invalid json": {
			req:         invalidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agent_groups/validate",
		},

		"validate a empty token": {
			req:         validJson,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			location:    "/agent_groups/validate",
		},
		"validate a agent group without content type": {
			req:         validJson,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/agent_groups/validate",
		},
		"validate a agent group with a invalid tag": {
			req:         invalidValueTag,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agent_groups/validate",
		},
		"validate a agent group with a name not respecting RegEx": {
			req:         invalidValueName,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agent_groups/validate",
		},
		"validate a agent group with a invalid token": {
			req:         validJson,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/agent_groups/validate",
		},
		"validate a agent group with a invalid agent group field": {
			req:         invalidField,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusBadRequest,
			location:    "/agent_groups/validate",
		},
		"validate a agent group with a invalid name": {
			req:         invalidNameJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agent_groups/validate",
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/agent_groups/validate", cli.server.URL),
				contentType: tc.contentType,
				token:       fmt.Sprintf("Bearer %s", tc.auth),
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestViewAgent(t *testing.T) {
	cli := newClientServer(t)

	ag, err := createAgent(t, "my-agent1", &cli)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		id     string
		auth   string
		status int
	}{
		"view a existing agent": {
			id:     ag.MFThingID,
			auth:   token,
			status: http.StatusOK,
		},
		"view a non-existing agent": {
			id:     wrongID,
			auth:   token,
			status: http.StatusNotFound,
		},
		"view a agent with a invalid token": {
			id:     ag.MFThingID,
			auth:   invalidToken,
			status: http.StatusUnauthorized,
		},
		"view a agent with a empty token": {
			id:     ag.MFThingID,
			auth:   "",
			status: http.StatusUnauthorized,
		},
		"view a existing agent with empty id": {
			id:     "",
			auth:   token,
			status: http.StatusBadRequest,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/agents/%s", cli.server.URL, tc.id),
				token:  fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected erro %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestViewAgentMatchingGroups(t *testing.T) {
	cli := newClientServer(t)

	ag, err := createAgent(t, "my-agent1", &cli)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		id     string
		auth   string
		status int
	}{
		"view a existing agent matching groups": {
			id:     ag.MFThingID,
			auth:   token,
			status: http.StatusOK,
		},
		"view a agent matching groups with a invalid token": {
			id:     ag.MFThingID,
			auth:   invalidToken,
			status: http.StatusUnauthorized,
		},
		"view a agent matching groups with a empty token": {
			id:     ag.MFThingID,
			auth:   "",
			status: http.StatusUnauthorized,
		},
		"view agent matching groups with empty id": {
			id:     "",
			auth:   token,
			status: http.StatusBadRequest,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/agents/%s/matching_groups", cli.server.URL, tc.id),
				token:  fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected erro %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestListAgent(t *testing.T) {
	cli := newClientServer(t)

	var data []agentRes
	for i := 0; i < limit; i++ {
		ag, err := createAgent(t, fmt.Sprintf("my-agent-%d", i), &cli)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		data = append(data, agentRes{
			ID:            ag.MFThingID,
			Name:          ag.Name.String(),
			ChannelID:     ag.MFChannelID,
			AgentTags:     ag.AgentTags,
			OrbTags:       ag.OrbTags,
			TsCreated:     ag.Created,
			AgentMetadata: ag.AgentMetadata,
			State:         ag.State.String(),
			LastHBData:    ag.LastHBData,
			TsLastHB:      ag.LastHB,
		})
	}

	cases := map[string]struct {
		auth   string
		status int
		url    string
		res    []agentRes
		total  uint64
	}{
		"retrieve a list of agents": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, limit),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agents with empty token": {
			auth:   "",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 1),
			res:    nil,
			total:  0,
		},
		"get a list of agents with invalid token": {
			auth:   invalidToken,
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 1),
			res:    nil,
			total:  0,
		},
		"get a list of agents with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agents with invalid dir": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=name&dir=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agents with negative offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", -1, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agents with negative limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, -5),
			res:    nil,
			total:  0,
		},
		"get a list of agents with zero limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 0),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agents without offset": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?limit=%d", limit),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agents without limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d", 1),
			res:    data[1:limit],
			total:  limit - 1,
		},
		"get a list of agents with limit greater than max": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 110),
			res:    nil,
			total:  0,
		},
		"get a list of agents with default URL": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s", ""),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of agents with invalid number of params": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=4&limit=4&limit=5&offset=5"),
			res:    nil,
			total:  0,
		},
		"get a list of agents with invalid offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=e&limit=5"),
			res:    nil,
			total:  0,
		},
		"get a list of agents with invalid limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=5&limit=e"),
			res:    nil,
			total:  0,
		},
		"get a list of agents filtering with invalid name": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&name=%s", 0, 5, invalidName),
			res:    nil,
			total:  0,
		},
		"get a list of agents sorted with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=wrong&dir=desc", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of agents sorted with invalid direction": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=name&dir=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodGet,
				url:         fmt.Sprintf(fmt.Sprintf("%s/agents%s", cli.server.URL, tc.url)),
				contentType: contentType,
				token:       fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var body agentsPageRes
			json.NewDecoder(res.Body).Decode(&body)
			total := uint64(len(body.Agents))
			assert.Equal(t, res.StatusCode, tc.status, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
			assert.Equal(t, total, tc.total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.total, total))
		})
	}

}

func TestUpdateAgent(t *testing.T) {
	cli := newClientServer(t)

	ag, err := createAgent(t, "my-agent1", &cli)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		req         string
		id          string
		contentType string
		auth        string
		status      int
	}{
		"update existing agent": {
			req: toJSON(updateAgentReq{
				Name: ag.Name.String(),
				Tags: ag.OrbTags,
			}),
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
		},
		"update agent with a empty json request": {
			req:         "{}",
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update agent with a invalid id": {
			req: toJSON(updateAgentReq{
				Name: ag.Name.String(),
				Tags: ag.OrbTags,
			}),
			id:          "invalid",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"update non-existing agent": {
			req: toJSON(updateAgentReq{
				Name: ag.Name.String(),
				Tags: ag.OrbTags,
			}),
			id:          wrongID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"update agent with invalid user token": {
			req: toJSON(updateAgentReq{
				Name: ag.Name.String(),
				Tags: ag.OrbTags,
			}),
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        "invalid",
			status:      http.StatusUnauthorized,
		},
		"update agent with empty user token": {
			req: toJSON(updateAgentReq{
				Name: ag.Name.String(),
				Tags: ag.OrbTags,
			}),
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
		},
		"update agent with invalid content type": {
			req: toJSON(updateAgentReq{
				Name: ag.Name.String(),
				Tags: ag.OrbTags,
			}),
			id:          ag.MFThingID,
			contentType: "invalid",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"update agent without content type": {
			req: toJSON(updateAgentReq{
				Name: ag.Name.String(),
				Tags: ag.OrbTags,
			}),
			id:          ag.MFThingID,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"update agent with a empty request": {
			req:         "",
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update agent with a invalid data format": {
			req:         invalidJson,
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update agent with different owner": {
			req:         invalidJson,
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update existing agent with invalid name": {
			req: toJSON(updateAgentReq{
				Name: "a",
				Tags: ag.OrbTags,
			}),
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update existing agent without name": {
			req: toJSON(updateAgentReq{
				Tags: ag.OrbTags,
			}),
			id:          ag.MFThingID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPut,
				url:         fmt.Sprintf("%s/agents/%s", cli.server.URL, tc.id),
				contentType: tc.contentType,
				token:       fmt.Sprintf("Bearer %s", tc.auth),
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestValidateAgent(t *testing.T) {
	var (
		validJson       = "{\"name\":\"eu-agents\",\"orb_tags\": {\"region\":\"eu\",   \"node_type\":\"dns\"}}"
		invalidTag      = "{\"name\":\"eu-agents\",\"orb_tags\": {\n\"invalidTag\", \n \"node_type\":\"dns\"}}"
		invalidName     = "{\"name\":\",,AGENT 6,\",\"orb_tags\": {\"region\":\"eu\",   \"node_type\":\"dns\"}}"
		invalidField    = "{\"nname\":\"eu-agents\",\"orb_tags\": {\"region\":\"eu\",   \"node_type\":\"dns\"}}"
		invalidNameJson = "{\"name\":\"a\",\"orb_tags\": {\"region\":\"eu\",   \"node_type\":\"dns\"}}"
	)

	cli := newClientServer(t)
	defer cli.server.Close()

	cases := map[string]struct {
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"validate a valid agent": {
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			location:    "/agents/validate",
		},
		"validate a agent with invalid json": {
			req:         invalidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agents/validate",
		},
		"validate a agent with a empty token": {
			req:         validJson,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			location:    "/agents/validate",
		},
		"validate a agent without a content type": {
			req:         validJson,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/agents/validate",
		},
		"validate a agent with a invalid tag": {
			req:         invalidTag,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agents/validate",
		},
		"validate a agent with name not following regex": {
			req:         invalidName,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agents/validate",
		},
		"validate a agent with a invalid token": {
			req:         validJson,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/agents/validate",
		},
		"validate a agent with a invalid agent field": {
			req:         invalidField,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agents/validate",
		},
		"validate a agent with invalid name": {
			req:         invalidNameJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agents/validate",
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/agents/validate", cli.server.URL),
				contentType: tc.contentType,
				token:       fmt.Sprintf("Bearer %s", tc.auth),
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestCreateAgent(t *testing.T) {
	var validJson = "{\"name\":\"eu-agents\",\"orb_tags\": {\"region\":\"eu\",   \"node_type\":\"dns\"}}"
	var conflictValidJson = "{\"name\":\"conflict\",\"orb_tags\": {\"region\":\"eu\",   \"node_type\":\"dns\"}}"
	var invalidNameJson = "{\"name\":\"a\",\"orb_tags\": {\"region\":\"eu\",   \"node_type\":\"dns\"}}"

	cli := newClientServer(t)
	defer cli.server.Close()

	_, err := createAgent(t, "conflict", &cli)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	cases := map[string]struct {
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"add a valid agent": {
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusCreated,
			location:    "/agents",
		},
		"add a duplicated agent": {
			req:         conflictValidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusConflict,
			location:    "/agents",
		},
		"add a valid agent with invalid token": {
			req:         validJson,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/agents",
		},
		"add a agent with an invalid json": {
			req:         invalidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agents",
		},
		"add a agent without a content type": {
			req:         validJson,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/agents",
		},
		"add a agent with an invalid content type": {
			req:         validJson,
			contentType: "invalid",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/agents",
		},
		"add a agent with an empty request": {
			req:         "{}",
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agents",
		},
		"add a agent with invalid name": {
			req:         invalidNameJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/agents",
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/agents", cli.server.URL),
				contentType: tc.contentType,
				token:       fmt.Sprintf("Bearer %s", tc.auth),
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestDeleteAgent(t *testing.T) {

	cli := newClientServer(t)

	ag, err := createAgent(t, "my-agent", &cli)
	require.Nil(t, err, "unexpected error: %s", err)

	cases := map[string]struct {
		id     string
		auth   string
		status int
	}{
		"delete existing agent": {
			id:     ag.MFThingID,
			auth:   token,
			status: http.StatusNoContent,
		},
		"delete non-existent agent": {
			id:     wrongID,
			auth:   token,
			status: http.StatusNoContent,
		},
		"delete agent with invalid token": {
			id:     ag.MFThingID,
			auth:   invalidToken,
			status: http.StatusUnauthorized,
		},
		"delete agent with empty token": {
			id:     ag.MFThingID,
			auth:   "",
			status: http.StatusUnauthorized,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodDelete,
				contentType: contentType,
				url:         fmt.Sprintf("%s/agents/%s", cli.server.URL, tc.id),
				token:       fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestAgentBackends(t *testing.T) {
	cli := newClientServer(t)

	cases := map[string]struct {
		auth   string
		status int
	}{
		"Return a list of available backends": {
			auth:   token,
			status: http.StatusOK,
		},
		"Return a list of available backends with invalid token": {
			auth:   invalidToken,
			status: http.StatusUnauthorized,
		},
		"Return a list of available backends with empty token": {
			auth:   "",
			status: http.StatusUnauthorized,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/agents/backends", cli.server.URL),
				token:  fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: Unexpected error: %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status %d got %d", desc, tc.status, res.StatusCode))

		})
	}
}

func TestAgentBackendHandler(t *testing.T) {
	cli := newClientServer(t)

	cases := map[string]struct {
		backend string
		auth    string
		status  int
	}{
		"Retrieve backend handler": {
			backend: "pktvisor",
			auth:    token,
			status:  http.StatusOK,
		},
		"Retrieve a handler with a non-existing backend": {
			backend: "orb",
			auth:    token,
			status:  http.StatusNotFound,
		},
		"Retrieve a handler with a invalid token": {
			backend: "pktvisor",
			auth:    invalidToken,
			status:  http.StatusUnauthorized,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/agents/backends/%s/handlers", cli.server.URL, tc.backend),
				token:  fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: Unexpected error: %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: Expected status %d got %d", desc, tc.status, res.StatusCode))

		})
	}
}

func TestAgentBackendInput(t *testing.T) {
	cli := newClientServer(t)

	cases := map[string]struct {
		backend string
		auth    string
		status  int
	}{
		"Retrieve backend input": {
			backend: "pktvisor",
			auth:    token,
			status:  http.StatusOK,
		},
		"Retrieve a backend input with a non-existing backend": {
			backend: "orb",
			auth:    token,
			status:  http.StatusNotFound,
		},
		"Retrieve a backend input with a invalid token": {
			backend: "pktvisor",
			auth:    invalidToken,
			status:  http.StatusUnauthorized,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/agents/backends/%s/inputs", cli.server.URL, tc.backend),
				token:  fmt.Sprintf("Bearer %s", tc.auth),
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: Unexpected error: %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: Expected status %d got %d", desc, tc.status, res.StatusCode))

		})
	}
}

func TestAgentBackendTaps(t *testing.T) {
	cli := newClientServer(t)

	cases := map[string]struct {
		token   string
		backend string
		status  int
	}{
		//"Retrieve taps by a provided backend": {
		//	token:   token,
		//	backend: "pktvisor",
		//	status:  http.StatusOK,
		//},
		"Retrieve taps by a non-existing backend": {
			token:   token,
			backend: "orb",
			status:  http.StatusNotFound,
		},
		"Retrieve taps by a provided backend with a invalid token": {
			token:   invalidToken,
			backend: "pktvisor",
			status:  http.StatusUnauthorized,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/agents/backends/%s/taps", cli.server.URL, tc.backend),
				token:  fmt.Sprintf("Bearer %s", tc.token),
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: Unexpected error: %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: Expected status %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func createAgentGroup(t *testing.T, name string, cli *clientServer) (fleet.AgentGroup, error) {
	t.Helper()
	agCopy := agentGroup
	validName, err := types.NewIdentifier(name)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	agCopy.Name = validName
	agCopy.Tags = tags
	ag, err := cli.service.CreateAgentGroup(context.Background(), token, agCopy)
	if err != nil {
		return fleet.AgentGroup{}, err
	}
	return ag, nil
}

func createAgent(t *testing.T, name string, cli *clientServer) (fleet.Agent, error) {
	t.Helper()
	aCopy := agent
	validName, err := types.NewIdentifier(name)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	aCopy.Name = validName
	aCopy.OrbTags = tags
	a, err := cli.service.CreateAgent(context.Background(), token, aCopy)
	if err != nil {
		return fleet.Agent{}, err
	}
	return a, nil
}

func newClientServer(t *testing.T) clientServer {
	t.Helper()
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)
	fleetServer := newServer(fleetService)

	return clientServer{
		service: fleetService,
		server:  fleetServer,
	}
}

type agentGroupRes struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description,omitempty"`
	Tags           types.Tags     `json:"tags"`
	TsCreated      time.Time      `json:"ts_created,omitempty"`
	MatchingAgents types.Metadata `json:"matching_agents,omitempty"`
	created        bool
}

type agentGroupsPageRes struct {
	Total       uint64          `json:"total"`
	Offset      uint64          `json:"offset"`
	Limit       uint64          `json:"limit"`
	AgentGroups []agentGroupRes `json:"agentGroups"`
}
type agentRes struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	State         string         `json:"state"`
	Key           string         `json:"key,omitempty"`
	ChannelID     string         `json:"channel_id,omitempty"`
	AgentTags     types.Tags     `json:"agent_tags"`
	OrbTags       types.Tags     `json:"orb_tags"`
	AgentMetadata types.Metadata `json:"agent_metadata"`
	LastHBData    types.Metadata `json:"last_hb_data"`
	TsCreated     time.Time      `json:"ts_created"`
	TsLastHB      time.Time      `json:"ts_last_hb"`
	created       bool
}

type agentsPageRes struct {
	Total  uint64     `json:"total"`
	Offset uint64     `json:"offset"`
	Limit  uint64     `json:"limit"`
	Agents []agentRes `json:"agents"`
}

type updateAgentGroupReq struct {
	token       string
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Tags        types.Tags `json:"tags"`
}

type updateAgentReq struct {
	token string
	Name  string     `json:"name,omitempty"`
	Tags  types.Tags `json:"orb_tags"`
}
