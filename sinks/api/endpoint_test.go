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
	"github.com/mainflux/mainflux/things"
	thmocks "github.com/mainflux/mainflux/things/mocks"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	skmocks "github.com/ns1labs/orb/sinks/mocks"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
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
	notFoundRes  = toJSON(errorRes{things.ErrNotFound.Error()})
	unauthRes    = toJSON(errorRes{things.ErrUnauthorizedAccess.Error()})
	notSupported = toJSON(errorRes{sinks.ErrUnsupportedContentTypeSink.Error()})
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
	//return skmocks.NewSinkServiceMock()
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

func TestViewBackend(t *testing.T) {
	f := skmocks.NewSinkServiceMock()
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	bes, err := f.ListBackends(context.Background(), token)
	//bes, err := service.ListBackends(context.Background(), token)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	id := strings.Trim(string(bes[0]), "\n")
	be, err := f.GetBackend(context.Background(), token, id)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	data := toJSON(sinksBackendRes{
		Backend:     be.GetName(),
		Description: be.GetDescription(),
		Config:      be.GetConfig(),
	})

	cases := []struct {
		desc        string
		id          string
		contentType string
		auth        string
		status      int
		res         string
	}{
		{
			desc:        "view existing backend",
			id:          id,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			res:         data,
		},
		{
			desc:        "view non-existing backend",
			id:          "logstash",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
			res:         notFoundRes,
		},
		{
			desc:        "view backend by passing invalid token",
			id:          id,
			contentType: contentType,
			auth:        "blah",
			status:      http.StatusUnauthorized,
			res:         unauthRes,
		},
		{
			desc:        "view backend by passing empty token",
			id:          id,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			res:         unauthRes,
		},
		{
			desc:        "view backend by passing invalid id",
			id:          "invalid",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
			res:         notFoundRes,
		},
	}

	for _, sinkCase := range cases {
		req := testRequest{
			client:      server.Client(),
			method:      http.MethodGet,
			contentType: sinkCase.contentType,
			url:         fmt.Sprintf("%s/features/sinks/%s", server.URL, sinkCase.id),
			token:       sinkCase.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("unexpect error %s", err))
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("unexpect error %s", err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, sinkCase.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", sinkCase.desc, sinkCase.status, res.StatusCode))
		assert.Equal(t, sinkCase.res, data, fmt.Sprintf("%s: expected body %s got %s", sinkCase.desc, sinkCase.res, data))
	}

}

func TestViewBackends(t *testing.T) {
	f := skmocks.NewSinkServiceMock()
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	bes, err := f.ListBackends(context.Background(), token)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var backends []interface{}
	for _, bk := range bes {
		b, err := f.GetBackend(context.Background(), token, bk)
		if err != nil {
			require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		}
		backends = append(backends, b)
	}

	data := toJSON(sinksBackendsRes{
		Backends: backends,
	})

	cases := []struct {
		desc        string
		auth        string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "view existing backends",
			auth:        token,
			contentType: contentType,
			status:      http.StatusOK,
			res:         data,
		},
		{
			desc:        "view backends by passing invalid content type",
			auth:        token,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         notSupported,
		},
		{
			desc:        "view backends by passing invalid token",
			auth:        "blah",
			contentType: contentType,
			status:      http.StatusUnauthorized,
			res:         unauthRes,
		},
		{
			desc:        "view backends by passing empty token",
			auth:        "",
			contentType: contentType,
			status:      http.StatusUnauthorized,
			res:         unauthRes,
		},
	}

	for _, sinkCase := range cases {
		req := testRequest{
			client:      server.Client(),
			method:      http.MethodGet,
			contentType: sinkCase.contentType,
			url:         fmt.Sprintf("%s/features/sinks", server.URL),
			token:       sinkCase.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("unexpect error %s", err))
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("unexpect error %s", err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, sinkCase.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", sinkCase.desc, sinkCase.status, res.StatusCode))
		assert.Equal(t, sinkCase.res, data, fmt.Sprintf("%s: expected body %s got %s", sinkCase.desc, sinkCase.res, data))
	}

}


func TestViewSink(t *testing.T) {
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	sk, err := service.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	data := toJSON(sinkRes{
		ID: sk.ID,
		Name: sk.Name.String() ,
		Description: sk.Description,
		Backend: sk.Backend,
		Config: sk.Config,
		Tags: sk.Tags,
		TsCreated: sk.Created,
	})

	cases := []struct {
		desc        string
		id          string
		contentType string
		auth        string
		status      int
		res         string
	}{
		{
			desc:        "view existing sink",
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			res:         data,
		},
		{
			desc:        "view non-existing sink",
			id:          "logstash",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
			res:         notFoundRes,
		},
		{
			desc:        "view backend by passing invalid token",
			id:          sink.ID,
			contentType: contentType,
			auth:        "blah",
			status:      http.StatusUnauthorized,
			res:         unauthRes,
		},
		{
			desc:        "view backend by passing empty token",
			id:          sink.ID,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			res:         unauthRes,
		},
		{
			desc:        "view backend by passing invalid id",
			id:          "invalid",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
			res:         notFoundRes,
		},
	}

	for _, sinkCase := range cases {
		req := testRequest{
			client:      server.Client(),
			method:      http.MethodGet,
			contentType: sinkCase.contentType,
			url:         fmt.Sprintf("%s/sinks/%s", server.URL, sinkCase.id),
			token:       sinkCase.auth,
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("unexpect error %s", err))
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("unexpect error %s", err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, sinkCase.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", sinkCase.desc, sinkCase.status, res.StatusCode))
		assert.Equal(t, sinkCase.res, data, fmt.Sprintf("%s: expected body %s got %s", sinkCase.desc, sinkCase.res, data))
	}

}