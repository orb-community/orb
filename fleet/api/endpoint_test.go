// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/mainflux/mainflux/things"
	thingsapi "github.com/mainflux/mainflux/things/api/things/http"
	"github.com/ns1labs/orb/fleet"
	api "github.com/ns1labs/orb/fleet/api"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
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
	var logger *zap.Logger
	config := mfsdk.Config{
		BaseURL: url,
	}

	mfsdk := mfsdk.NewSDK(config)
	return fleet.NewFleetService(logger, auth, nil, agentGroupRepo, nil, mfsdk)
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
			status:      http.StatusOK,
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
