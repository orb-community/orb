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
	"strings"
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
	limit        = 10
	invalidToken = "invalid"
	maxNameSize  = 1024
	contentType  = "application/json"
	wrongID      = "28ea82e7-0224-4798-a848-899a75cdc650"
)

var (
	validJson = `{
    "name": "mypktvisorpolicyyaml-3",
    "description": "my pktvisor policy yaml",
    "tags": {
        "region": "eu",
        "node_type": "dns"
    },
    "format": "yaml",
    "policy_data": "version: \"1.0\"\nvisor:\n    foo: \"bar\""
}`
	metadata    = map[string]interface{}{"type": "orb_agent"}
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
	return policies.New(nil, auth, policyRepo)
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
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/policies/agent/%s", cli.server.URL, tc.ID),
				token:  tc.token,
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: Unexpected error: %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestListPolicies(t *testing.T) {
	cli := newClientServer(t)

	var data []policyRes
	for i := 0; i < limit; i++ {
		p := createPolicy(t, &cli, fmt.Sprintf("policy-%d", i))
		data = append(data, policyRes{
			ID:      p.ID,
			Name:    p.Name.String(),
			Backend: p.Backend,
			created: true,
		})
	}

	cases := map[string]struct {
		auth   string
		status int
		url    string
		res    []policyRes
		total  uint64
	}{
		"retrieve a list of policies": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, limit),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of policies with empty token": {
			auth:   "",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 1),
			res:    nil,
			total:  0,
		},
		"get a list of policies with invalid token": {
			auth:   invalidToken,
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 1),
			res:    nil,
			total:  0,
		},
		"get a list of policies with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of policies with invalid dir": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=name&dir=wrong", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of policies with negative offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", -1, 5),
			res:    nil,
			total:  0,
		},
		"get a list of policies with negative limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, -5),
			res:    nil,
			total:  0,
		},
		"get a list of policies with zero limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 0),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of policies without offset": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?limit=%d", limit),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of policies without limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("?offset=%d", 1),
			res:    data[1:limit],
			total:  limit - 1,
		},
		"get a list of policies with limit greater than max": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d", 0, 110),
			res:    nil,
			total:  0,
		},
		"get a list of policies with default URL": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s", ""),
			res:    data[0:limit],
			total:  limit,
		},
		"get a list of policies with invalid number of params": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=4&limit=4&limit=5&offset=5"),
			res:    nil,
			total:  0,
		},
		"get a list of policies with invalid offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=e&limit=5"),
			res:    nil,
			total:  0,
		},
		"get a list of policies with invalid limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=5&limit=e"),
			res:    nil,
			total:  0,
		},
		"get a list of policies filtering with invalid name": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&name=%s", 0, 5, invalidName),
			res:    nil,
			total:  0,
		},
		"get a list of policies sorted with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("?offset=%d&limit=%d&order=wrong&dir=desc", 0, 5),
			res:    nil,
			total:  0,
		},
		"get a list of policies sorted with invalid direction": {
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
				client: cli.server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf(fmt.Sprintf("%s/policies/agent%s", cli.server.URL, tc.url)),
				token:  tc.auth,
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var body policiesPageRes
			json.NewDecoder(res.Body).Decode(&body)
			total := uint64(len(body.Policies))
			assert.Equal(t, res.StatusCode, tc.status, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
			assert.Equal(t, total, tc.total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.total, total))
		})
	}
}

func TestPolicyEdition(t *testing.T) {
	cli := newClientServer(t)
	policy := createPolicy(t, &cli, "policy")

	cases := map[string]struct {
		id          string
		contentType string
		auth        string
		status      int
		data        string
	}{
		"update a existing policy": {
			id:          policy.ID,
			contentType: "application/json",
			auth:        token,
			status:      http.StatusOK,
			data:        validJson,
		},
		"update policy with a empty json request": {
			id:          policy.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			data:        "{}",
		},
		"update policy with a invalid id": {
			data:        validJson,
			id:          "invalid",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"update non-existing policy": {
			data:        validJson,
			id:          wrongID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"update policy with invalid user token": {
			data:        validJson,
			id:          policy.ID,
			contentType: contentType,
			auth:        "invalid",
			status:      http.StatusUnauthorized,
		},
		"update policy with empty user token": {
			data:        validJson,
			id:          policy.ID,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
		},
		"update policy with invalid content type": {
			data:        validJson,
			id:          policy.ID,
			contentType: "invalid",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"update policy without content type": {
			data:        validJson,
			id:          policy.ID,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"update policy with a empty request": {
			data:        "",
			id:          policy.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update policy with a invalid data format": {
			data:        "{",
			id:          policy.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPut,
				url:         fmt.Sprintf("%s/policies/agent/%s", cli.server.URL, tc.id),
				contentType: tc.contentType,
				token:       tc.auth,
				body:        strings.NewReader(tc.data),
			}
			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: Unexpected error: %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestPolicyRemoval(t *testing.T) {
	cli := newClientServer(t)

	plcy := createPolicy(t, &cli, "policy")

	cases := map[string]struct {
		id     string
		auth   string
		status int
	}{
		"delete a existing policy": {
			id:     plcy.ID,
			auth:   token,
			status: http.StatusNoContent,
		},
		"delete non-existent policy": {
			id:     wrongID,
			auth:   token,
			status: http.StatusNoContent,
		},
		"delete policy with invalid token": {
			id:     plcy.ID,
			auth:   invalidToken,
			status: http.StatusUnauthorized,
		},
		"delete policy with empty token": {
			id:     plcy.ID,
			auth:   "",
			status: http.StatusUnauthorized,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: cli.server.Client(),
				method: http.MethodDelete,
				url:    fmt.Sprintf("%s/policies/agent/%s", cli.server.URL, tc.id),
				token:  tc.auth,
			}

			res, err := req.make()
			require.Nil(t, err, fmt.Sprintf("%s: Unexpected error: %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestDatasetValidation(t *testing.T){
	var (
		validYaml = `{"name": "mydatasetyaml-3", "agent_group_id": "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", "agent_policy_id": "86b7b412-1b7f-f5bc-c78b-f79087d6e49b", "sink_id": "urn:uuid:f5b2d342-211d-a9ab-1233-63199a3fc16f", "tags": {"region": "eu", "node_type": "dns"}}`
		invalidJson = `{`
		//testvalidJson         = "{\n    \"name\": \"mydatasetyaml-3\",\n    \"agent_group_id\": \"8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db\",\n    \"agent_policy_id\": \"86b7b412-1b7f-f5bc-c78b-f79087d6e49b\",\n    \"sink_id\": urn:uuid:f5b2d342-211d-a9ab-1233-63199a3fc16f\n,\n    \"tags\": {\n        \"region\": \"eu\",\n        \"node_type\": \"dns\"\n    }}"
		invalidNameYaml = `{"name": "9...DATASET", "agent_group_id": "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", "agent_policy_id": "86b7b412-1b7f-f5bc-c78b-f79087d6e49b", "sink_id": "urn:uuid:f5b2d342-211d-a9ab-1233-63199a3fc16f", "tags": {"region": "eu", "node_type": "dns"}}`
		invalidTagYaml = `{"name": "mydatasetyaml-3", "agent_group_id": "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", "agent_policy_id": "86b7b412-1b7f-f5bc-c78b-f79087d6e49b", "sink_id": "urn:uuid:f5b2d342-211d-a9ab-1233-63199a3fc16f", "tags": "invalidTag"}`
		invalidFieldYaml = `{"naamme": "mydatasetyaml-3", "agent_group_id": "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", "agent_policy_id": "86b7b412-1b7f-f5bc-c78b-f79087d6e49b", "sink_id": "urn:uuid:f5b2d342-211d-a9ab-1233-63199a3fc16f", "tags": {"region": "eu", "node_type": "dns"}}`
	)
	cli := newClientServer(t)

	cases := map[string]struct {
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"Validate a valid dataset": {
			req:         validYaml,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			location:    "/policies/dataset/validate",
		},
		"Validate a invalid yaml": {
			req:         invalidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/policies/dataset/validate",
		},
		"Validate a dataset with a empty token": {
			req:         validYaml,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			location:    "/policies/dataset/validate",
		},
		"Validate a dataset with a invalid token": {
			req:         validYaml,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/policies/dataset/validate",
		},
		"Validate a dataset without a content type": {
			req:         validYaml,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/policies/dataset/validate",
		},
		"Validate a dataset with a invalid name value": {
			req:         invalidNameYaml,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusBadRequest,
			location:    "/policies/dataset/validate",
		},
		"Validate a dataset with a invalid tag value": {
			req:         invalidTagYaml,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/policies/dataset/validate",
		},
		"Validate a dataset with a invalid field": {
			req:         invalidFieldYaml,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/policies/dataset/validate",
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      cli.server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/policies/dataset/validate", cli.server.URL),
				contentType: tc.contentType,
				token:       tc.auth,
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			if err != nil {
				require.Nil(t, err, "%s: Unexpected error: %s", desc, err)
			}
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func createPolicy(t *testing.T, cli *clientServer, name string) policies.Policy {
	t.Helper()
	ID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	validName, err := types.NewIdentifier(name)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := policies.Policy{
		ID:      ID.String(),
		Name:    validName,
		Backend: "pktvisor",
	}

	res, err := cli.service.AddPolicy(context.Background(), token, policy, format, policy_data)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
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

type policyRes struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Backend string `json:"backend"`
	created bool
}

type policiesPageRes struct {
	Total    uint64      `json:"total"`
	Offset   uint64      `json:"offset"`
	Limit    uint64      `json:"limit"`
	Policies []policyRes `json:"data"`
}
