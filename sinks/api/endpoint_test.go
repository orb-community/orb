// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"encoding/json"
	"fmt"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	thmocks "github.com/mainflux/mainflux/things/mocks"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	skmocks "github.com/ns1labs/orb/sinks/mocks"
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
)

const (
	contentType = "application/json"
	token       = "token"
	invalidToken = "invalid_token"
	wrongID		= 0
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

func TestListSinks(t *testing.T) {
	svc := newService(map[string]string{token: email})
	server := newServer(svc)
	defer server.Close()

	var data []sinks.Sink
	for i := 0; i < 20; i++ {
		var skName, _ = types.NewIdentifier(fmt.Sprintf("name%d", i))
		snk := sinks.Sink{
			Name:        skName,
			Description: "An example prometheus sink",
			Backend:     "prometheus",
			Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
			Tags:        map[string]string{"cloud": "aws"},
		}

		sk, err := svc.CreateSink(context.Background(), token, snk)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		data = append(data, sk)
	}

	sinkURL := fmt.Sprintf("%s/sinks", server.URL)

	cases := []struct{
		desc        string
		auth        string
		status      int
		url    		string
		res			[]sinks.Sink
	}{
		{
			desc: "get a list of sinks",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 0, 5),
			res: data[0:5],
		},
		{
			desc: "get a list of sinks with empty token",
			auth: "",
			status: http.StatusUnauthorized,
			url: fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 0, 1),
			res: nil,
		},
		{
			desc: "get a list of sinks with invalid token",
			auth: invalidToken,
			status: http.StatusUnauthorized,
			url: fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 0, 1),
			res: nil,
		},
		{
			desc: "get a list of sinks ordered by name descendent",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=desc", sinkURL, 0, 5),
			res: data[0:5],
		},
		{
			desc: "get a list of sinks ordered by name ascendent",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=asc", sinkURL, 0, 5),
			res: data[0:5],
		},
		{
			desc: "get a list of sinks with invalid order",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&order=wrong", sinkURL, 0, 5),
			res: nil,
		},
		{
			desc: "get a list of sinks with invalid dir",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=wrong", sinkURL, 0, 5),
			res: nil,
		},
		{
			desc: "get a list of sinks with negative offset",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, -1, 5),
			res: nil,
		},
		{
			desc: "get a list of sinks with negative limit",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 1, -5),
			res: nil,
		},
		{
			desc: "get a list of sinks with offset 1 and zero limit",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 1, 0),
			res: data[1:8],
		},
		{
			desc: "get a list of sinks without offset",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?limit=%d", sinkURL, 5),
			res: data[0:5],
		},
		{
			desc: "get a list of sinks without limit",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?offset=%d", sinkURL, 1),
			res: data[1:11],
		},
		{
			desc: "get a list of sinks with redundant query params",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&value=something", sinkURL, 0, 5),
			res: data[0:5],
		},
		{
			desc: "get a list of sinks with limit greater than max",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 0, 110),
			res: nil,
		},
		{
			desc: "get a list of sinks with default URL",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s%s", sinkURL, ""),
			res: data[0:10],
		},
		{
			desc: "get a list of sinks with invalid number of params",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s%s", sinkURL, "?offset=4&limit=4&limit=5&offset=5"),
			res: nil,
		},
		{
			desc: "get a list of sinks with invalid offset",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s%s", sinkURL, "?offset=e&limit=5"),
			res: nil,
		},
		{
			desc: "get a list of sinks with invalid limit",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s%s", sinkURL, "?offset=5&limit=e"),
			res: nil,
		},
		{
			desc: "get a list of sinks filtering with invalid name",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&name=%s", sinkURL, 0, 5, invalidName),
			res: nil,
		},
		{
			desc: "get a list of sinks sorted by name ascendent",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=asc", sinkURL, 0, 5),
			res: data[0:5],
		},
		{
			desc: "get a list of sinks sorted by name descendent",
			auth: token,
			status: http.StatusOK,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=desc", sinkURL, 0, 5),
			res: data[0:5],
		},
		{
			desc: "get a list of sinks sorted with invalid order",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&order=wrong&dir=desc", sinkURL, 0, 5),
			res: nil,
		},
		{
			desc: "get a list of sinks sorted with invalid order",
			auth: token,
			status: http.StatusBadRequest,
			url: fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=wrong", sinkURL, 0, 5),
			res: nil,
		},
	}

	for _, sinkCase := range cases {
		req := testRequest{
			client: server.Client(),
			method: http.MethodGet,
			url: sinkCase.url,
			token: sinkCase.auth,
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		assert.Equal(t, sinkCase.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d, got %d", sinkCase.desc, sinkCase.status, res.StatusCode))
	}
}

func TestDeleteSink(t *testing.T) {
	svc := newService(map[string]string{token: email})
	server := newServer(svc)
	defer server.Close()

	sk, err := svc.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := []struct{
		desc string
		id string
		auth string
		status int
	}{
		{
			desc: "delete existing sink",
			id: sk.ID,
			auth: token,
			status: http.StatusOK,
		},
		{
			desc: "delete non-existent sink",
			id: strconv.FormatUint(wrongID, 10),
			auth: token,
			status: http.StatusNotFound,
		},
		{
			desc: "delete sink with invalid token",
			id: sk.ID,
			auth: invalidToken,
			status: http.StatusUnauthorized,
		},
		{
			desc: "delete sink with empty token",
			id: sk.ID,
			auth: "",
			status: http.StatusUnauthorized,
		},
	}

	for _, sinkCase := range cases {
		req := testRequest{
			client: server.Client(),
			method: http.MethodDelete,
			url: fmt.Sprintf("%s/sinks/%s", server.URL, sinkCase.id),
			token: sinkCase.auth,
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", sinkCase.desc, err))
		assert.Equal(t, sinkCase.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", sinkCase.desc, sinkCase.status, res.StatusCode))
	}
}