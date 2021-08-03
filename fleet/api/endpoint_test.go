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
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/fleet"
	api "github.com/ns1labs/orb/fleet/api"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
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

func newService(tokens map[string]string) fleet.Service {
	auth := flmocks.NewAuthService(tokens)
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	var logger *zap.Logger
	config := mfsdk.Config{
		BaseURL: "http://localhost",
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
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

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
			client:      server.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/agent_groups", server.URL),
			contentType: tc.contentType,
			token:       tc.auth,
			body:        strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}

}
