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
	"github.com/gofrs/uuid"
	"github.com/mainflux/mainflux"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	api "github.com/ns1labs/orb/policies/api/http"
	plmocks "github.com/ns1labs/orb/policies/mocks"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	token       = "token"
	email       = "user@example.com"
	format      = "yaml"
	policy_data = `version: "1.0"
visor:
  taps:
    anycast:
      type: pcap
      config:
        iface: eth0`
)

var (
	metadata = map[string]interface{}{"type": "orb_agent"}
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
	service policies.Service
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

func newService(auth mainflux.AuthServiceClient) policies.Service {
	policyRepo := plmocks.NewPoliciesRepository()
	return policies.New(auth, policyRepo)
}

func newServer(svc policies.Service) *httptest.Server {
	mux := api.MakeHandler(mocktracer.New(), "policies", svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestViewPolicy(t *testing.T) {
	cli := newClientServer(t)
	policy := createPolicy(t, &cli, "policy")

	cases := map[string]struct {
		ID     string
		token  string
		status int
	}{
		"view a existing policy": {
			ID:     policy.ID,
			token:  token,
			status: http.StatusOK,
		},
		"view a non-existing policy": {
			ID:     "d0967904-8824-4ed1-b11c-9a92f9e4e43c",
			token:  token,
			status: http.StatusNotFound,
		},
		"view a policy with a invalid token": {
			ID:     policy.ID,
			token:  "invalid",
			status: http.StatusUnauthorized,
		},
		"view a policy with a empty token": {
			ID:     policy.ID,
			token:  "",
			status: http.StatusUnauthorized,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			//t.Parallel()
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/policies/%s", cli.server.URL, tc.ID),
				token:  tc.token,
			}
			res, err := req.make()
			if err != nil {
				require.Nil(t, err, "%s: Unexpected error: %s", desc, err)
			}
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestListPolicies(t *testing.T) {
	cli := newClientServer(t)

	var policies []policies.Policy
	for i := 0; i < 10; i++ {
		p := createPolicy(t, &cli, fmt.Sprintf("policy-%d", i))
		policies = append(policies, p)
	}
}

func createPolicy(t *testing.T, cli *clientServer, name string) policies.Policy {
	t.Helper()
	ID, err := uuid.NewV4()
	if err != nil {
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	}

	validName, err := types.NewIdentifier(name)
	if err != nil {
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	}

	policy := policies.Policy{
		ID:      ID.String(),
		Name:    validName,
		Backend: "pktvisor",
	}

	res, err := cli.service.CreatePolicy(context.Background(), token, policy, format, policy_data)
	if err != nil {
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	}
	return res
}

func newClientServer(t *testing.T) clientServer {
	t.Helper()
	users := flmocks.NewAuthService(map[string]string{token: email})

	policiesService := newService(users)
	policiesServer := newServer(policiesService)

	return clientServer{
		service: policiesService,
		server:  policiesServer,
	}
}
