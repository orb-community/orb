// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"encoding/json"
	"fmt"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	thmocks "github.com/mainflux/mainflux/things/mocks"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	skmocks "github.com/ns1labs/orb/sinks/mocks"
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
	contentType = "application/json"
	token       = "token"
	email       = "user@example.com"
	validJson   = "{\n    \"name\": \"my-prom-sink\",\n    \"backend\": \"prometheus\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": {\n        \"cloud\": \"aws\"\n    },\n    \"validate_only\": false\n}"
	invalidJson = "{"
)

var (
	sink = sinks.Sink{
		Name:        types.Identifier{},
		Description: "An example prometheus sink",
		Backend:     "prometheus",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}
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

func newService(tokens map[string]string) sinks.Service {
	auth := thmocks.NewAuthService(tokens)
	sinkRepo := skmocks.NewSinkRepository()
	var logger *zap.Logger

	config := mfsdk.Config{
		BaseURL:      "localhost",
		ThingsPrefix: "",
	}

	mfsdk := mfsdk.NewSDK(config)
	return sinks.NewSinkService(logger, auth, sinkRepo, mfsdk)
}

func newServer(svc sinks.Service) *httptest.Server {
	mux := MakeHandler(mocktracer.New(), "sinks", svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestCreateSinks(t *testing.T) {
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
			desc:        "add a valid sink",
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			location:    "/sinks",
		},
		{
			desc:        "add a duplicate sink",
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusConflict,
			location:    "/sinks",
		},
		{
			desc:        "add sink with invalid json",
			req:         invalidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks",
		},
		{
			desc:        "add a sink with a invalid token",
			req:         validJson,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			location:    "/sinks",
		},
		{
			desc:        "add a valid without content type",
			req:         validJson,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/sinks",
		},
	}

	for _, sinkCase := range cases {
		req := testRequest{
			client:      server.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/sinks", server.URL),
			contentType: sinkCase.contentType,
			token:       sinkCase.auth,
			body:        strings.NewReader(sinkCase.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("unexpect erro %s", err))

		assert.Equal(t, sinkCase.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", sinkCase.desc, sinkCase.status, res.StatusCode))
	}

}
