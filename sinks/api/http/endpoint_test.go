// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
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
	contentType       = "application/json"
	token             = "token"
	invalidToken      = "invalid"
	email             = "user@example.com"
	validJson         = "{\n    \"name\": \"my-prom-sink\",\n    \"backend\": \"prometheus\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": {\n        \"cloud\": \"aws\"\n    },\n    \"validate_only\": false\n}"
	conflictValidJson = "{\n    \"name\": \"conflict\",\n    \"backend\": \"prometheus\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": {\n        \"cloud\": \"aws\"\n    },\n    \"validate_only\": false\n}"
	invalidJson       = "{"
)

var (
	nameID, _ = types.NewIdentifier("my-sink")
	sink      = sinks.Sink{
		Name:        nameID,
		Description: "An example prometheus sink",
		Backend:     "prometheus",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}
	invalidName        = strings.Repeat("m", maxNameSize+1)
	notFoundRes        = toJSON(errorRes{sinks.ErrNotFound.Error()})
	unauthRes          = toJSON(errorRes{sinks.ErrUnauthorizedAccess.Error()})
	notSupported       = toJSON(errorRes{sinks.ErrUnsupportedContentTypeSink.Error()})
	malformedEntityRes = toJSON(errorRes{sinks.ErrMalformedEntity.Error()})
	wrongID, _         = uuid.NewV4()
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

func newService(tokens map[string]string) sinks.SinkService {
	auth := skmocks.NewAuthService(tokens)
	sinkRepo := skmocks.NewSinkRepository()
	var logger *zap.Logger

	config := mfsdk.Config{
		BaseURL:      "localhost",
		ThingsPrefix: "",
	}

	mfsdk := mfsdk.NewSDK(config)
	return sinks.NewSinkService(logger, auth, sinkRepo, mfsdk)
}

func newServer(svc sinks.SinkService) *httptest.Server {
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

	// Conflict creation scenario
	sinkConflict := sink
	conflictNameID, err := types.NewIdentifier("conflict")
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	sinkConflict.Name = conflictNameID
	_, err = service.CreateSink(context.Background(), token, sinkConflict)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var invalidNameJson = "{\n    \"name\": \"s\",\n    \"backend\": \"prometheus\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": {\n        \"cloud\": \"aws\"\n    },\n    \"validate_only\": false\n}"
	var emptyNameJson = "{\n    \"name\": \"\",\n    \"backend\": \"prometheus\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": {\n        \"cloud\": \"aws\"\n    },\n    \"validate_only\": false\n}"

	cases := map[string]struct {
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"add a valid sink": {
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusCreated,
			location:    "/sinks",
		},
		"add a duplicate sink": {
			req:         conflictValidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusConflict,
			location:    "/sinks",
		},
		"add sink with invalid json": {
			req:         invalidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks",
		},
		"add a sink with a invalid token": {
			req:         validJson,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			location:    "/sinks",
		},
		"add a valid without content type": {
			req:         validJson,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/sinks",
		},
		"add a sink with invalid name": {
			req:         invalidNameJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks",
		},
		"add a sink with empty name": {
			req:         emptyNameJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks",
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/sinks", server.URL),
				contentType: tc.contentType,
				token:       tc.auth,
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}

}

func TestUpdateSink(t *testing.T) {
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	sk, err := service.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	data := toJSON(updateSinkReq{
		Name:        sk.Name.String(),
		Description: sk.Description,
		Config:      sink.Config,
		Tags:        sk.Tags,
	})

	dataInvalidName := toJSON(updateSinkReq{
		Name:        invalidName,
		Description: sk.Description,
		Config:      sink.Config,
		Tags:        sk.Tags,
	})

	dataInvalidRgxName := toJSON(updateSinkReq{
		Name:        "&*sink*&",
		Description: sk.Description,
		Config:      sink.Config,
		Tags:        sk.Tags,
	})

	cases := map[string]struct {
		req         string
		id          string
		contentType string
		auth        string
		status      int
	}{
		"update existing sink": {
			req:         data,
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
		},
		"update sink with a empty json request": {
			req:         "{}",
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update sink with a invalid id": {
			req:         data,
			id:          "invalid",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"update non-existing sink": {
			req:         data,
			id:          wrongID.String(),
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
		},
		"update sink with invalid user token": {
			req:         data,
			id:          sk.ID,
			contentType: contentType,
			auth:        "invalid",
			status:      http.StatusUnauthorized,
		},
		"update sink with empty user token": {
			req:         data,
			id:          sk.ID,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
		},
		"update sink with invalid content type": {
			req:         data,
			id:          sk.ID,
			contentType: "invalid",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"update sink without content type": {
			req:         data,
			id:          sk.ID,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
		},
		"update sink with a empty request": {
			req:         "",
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update sink with a invalid data format": {
			req:         invalidJson,
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update sink with different owner": {
			req:         invalidJson,
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update existing sink with a invalid name": {
			req:         dataInvalidName,
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
		"update existing sink with a invalid regex name": {
			req:         dataInvalidRgxName,
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      server.Client(),
				method:      http.MethodPut,
				url:         fmt.Sprintf("%s/sinks/%s", server.URL, tc.id),
				contentType: tc.contentType,
				token:       tc.auth,
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			require.Nil(t, err, "%s: unexpected error: %s", desc, err)
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
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

	cases := map[string]struct {
		auth   string
		status int
		url    string
		res    []sinks.Sink
	}{
		"get a list of sinks": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 0, 5),
			res:    data[0:5],
		},
		"get a list of sinks with empty token": {
			auth:   "",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 0, 1),
			res:    nil,
		},
		"get a list of sinks with invalid token": {
			auth:   invalidToken,
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 0, 1),
			res:    nil,
		},
		"get a list of sinks ordered by name descendent": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=desc", sinkURL, 0, 5),
			res:    data[0:5],
		},
		"get a list of sinks ordered by name ascendent": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=asc", sinkURL, 0, 5),
			res:    data[0:5],
		},
		"get a list of sinks with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&order=wrong", sinkURL, 0, 5),
			res:    nil,
		},
		"get a list of sinks with invalid dir": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=wrong", sinkURL, 0, 5),
			res:    nil,
		},
		"get a list of sinks with negative offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, -1, 5),
			res:    nil,
		},
		"get a list of sinks with negative limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 1, -5),
			res:    nil,
		},
		"get a list of sinks with offset 1 and zero limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 1, 0),
			res:    data[1:8],
		},
		"get a list of sinks without offset": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?limit=%d", sinkURL, 5),
			res:    data[0:5],
		},
		"get a list of sinks without limit": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d", sinkURL, 1),
			res:    data[1:11],
		},
		"get a list of sinks with redundant query params": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&value=something", sinkURL, 0, 5),
			res:    data[0:5],
		},
		"get a list of sinks with limit greater than max": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", sinkURL, 0, 110),
			res:    nil,
		},
		"get a list of sinks with default URL": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s%s", sinkURL, ""),
			res:    data[0:10],
		},
		"get a list of sinks with invalid number of params": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s%s", sinkURL, "?offset=4&limit=4&limit=5&offset=5"),
			res:    nil,
		},
		"get a list of sinks with invalid offset": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s%s", sinkURL, "?offset=e&limit=5"),
			res:    nil,
		},
		"get a list of sinks with invalid limit": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s%s", sinkURL, "?offset=5&limit=e"),
			res:    nil,
		},
		"get a list of sinks filtering with invalid name": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&name=%s", sinkURL, 0, 5, invalidName),
			res:    nil,
		},
		"get a list of sinks sorted by name ascendent": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=asc", sinkURL, 0, 5),
			res:    data[0:5],
		},
		"get a list of sinks sorted by name descendent": {
			auth:   token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=desc", sinkURL, 0, 5),
			res:    data[0:5],
		},
		"get a list of sinks sorted with invalid order": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&order=wrong&dir=desc", sinkURL, 0, 5),
			res:    nil,
		},
		"get a list of sinks sorted with invalid direction": {
			auth:   token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d&order=name&dir=wrong", sinkURL, 0, 5),
			res:    nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: server.Client(),
				method: http.MethodGet,
				url:    tc.url,
				token:  tc.auth,
			}

			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d, got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestViewBackend(t *testing.T) {
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	bes, err := service.ListBackends(context.Background(), token)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	id := strings.Trim(string(bes[0]), "\n")
	be, err := service.ViewBackend(context.Background(), token, id)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	data := toJSON(sinksBackendRes{
		Backend: be.Metadata(),
	})

	cases := map[string]struct {
		id     string
		auth   string
		status int
		res    string
	}{
		"view existing backend": {
			id:     id,
			auth:   token,
			status: http.StatusOK,
			res:    data,
		},
		"view non-existing backend": {
			id:     "logstash",
			auth:   token,
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
		"view backend by passing invalid token": {
			id:     id,
			auth:   "blah",
			status: http.StatusUnauthorized,
			res:    unauthRes,
		},
		"view backend by passing empty token": {
			id:     id,
			auth:   "",
			status: http.StatusUnauthorized,
			res:    unauthRes,
		},
		"view backend by passing invalid id": {
			id:     "invalid",
			auth:   token,
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
		"view backend with empty id": {
			id:     "",
			auth:   token,
			status: http.StatusBadRequest,
			res:    malformedEntityRes,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/features/sinks/%s", server.URL, tc.id),
				token:  tc.auth,
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
			body, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
			data := strings.Trim(string(body), "\n")
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
			assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", desc, tc.res, data))
		})
	}

}

func TestViewBackends(t *testing.T) {
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	bes, err := service.ListBackends(context.Background(), token)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var backends []interface{}
	for _, bk := range bes {
		b, err := service.ViewBackend(context.Background(), token, bk)
		if err != nil {
			require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		}
		backends = append(backends, b.Metadata())
	}

	data := toJSON(sinksBackendsRes{
		Backends: backends,
	})

	cases := map[string]struct {
		auth   string
		status int
		res    string
	}{
		"view existing backends": {
			auth:   token,
			status: http.StatusOK,
			res:    data,
		},
		"view backends by passing invalid token": {
			auth:   "blah",
			status: http.StatusUnauthorized,
			res:    unauthRes,
		},
		"view backends by passing empty token": {
			auth:   "",
			status: http.StatusUnauthorized,
			res:    unauthRes,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: server.Client(),
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/features/sinks", server.URL),
				token:  tc.auth,
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
			body, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
			data := strings.Trim(string(body), "\n")
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
			assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", desc, tc.res, data))
		})
	}

}

func TestViewSink(t *testing.T) {
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	sk, err := service.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	data := toJSON(sinkRes{
		ID:          sk.ID,
		Name:        sk.Name.String(),
		Description: sk.Description,
		Backend:     sk.Backend,
		Config:      sk.Config,
		Tags:        sk.Tags,
		State:       sk.State.String(),
		Error:       sk.Error,
		TsCreated:   sk.Created,
	})

	cases := map[string]struct {
		id          string
		contentType string
		auth        string
		status      int
		res         string
	}{
		"view existing sink": {
			id:          sk.ID,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			res:         data,
		},
		"view non-existing sink": {
			id:          "logstash",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
			res:         notFoundRes,
		},
		"view backend by passing invalid token": {
			id:          sink.ID,
			contentType: contentType,
			auth:        "blah",
			status:      http.StatusUnauthorized,
			res:         unauthRes,
		},
		"view backend by passing empty token": {
			id:          sink.ID,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			res:         unauthRes,
		},
		"view backend by passing invalid id": {
			id:          "invalid",
			contentType: contentType,
			auth:        token,
			status:      http.StatusNotFound,
			res:         notFoundRes,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      server.Client(),
				method:      http.MethodGet,
				contentType: tc.contentType,
				url:         fmt.Sprintf("%s/sinks/%s", server.URL, tc.id),
				token:       tc.auth,
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
			body, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
			data := strings.Trim(string(body), "\n")
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
			assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", desc, tc.res, data))
		})
	}

}

func TestDeleteSink(t *testing.T) {
	svc := newService(map[string]string{token: email})
	server := newServer(svc)
	defer server.Close()
	sk, err := svc.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	cases := map[string]struct {
		id     string
		auth   string
		status int
	}{
		"delete existing sink": {
			id:     sk.ID,
			auth:   token,
			status: http.StatusNoContent,
		},
		"delete non-existent sink": {
			id:     wrongID.String(),
			auth:   token,
			status: http.StatusNoContent,
		},
		"delete sink with invalid token": {
			id:     sk.ID,
			auth:   invalidToken,
			status: http.StatusUnauthorized,
		},
		"delete sink with empty token": {
			id:     sk.ID,
			auth:   "",
			status: http.StatusUnauthorized,
		},
		"delete sink with empty id": {
			id:     "",
			auth:   token,
			status: http.StatusBadRequest,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client: server.Client(),
				method: http.MethodDelete,
				url:    fmt.Sprintf("%s/sinks/%s", server.URL, tc.id),
				token:  tc.auth,
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", desc, err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestValidateSink(t *testing.T) {
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	var invalidSinkField = "{\n    \"namee\": \"my-prom-sink\",\n    \"backend\": \"prometheus\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": {\n        \"cloud\": \"aws\"\n    }}"
	var invalidSinkValueName = "{\n    \"name\": \"my...SINK1\",\n    \"backend\": \"prometheus\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": {\n        \"cloud\": \"aws\"\n    }}"
	var invalidSinkValueBackend = "{\n    \"name\": \"my-prom-sink\",\n    \"backend\": \"invalidBackend\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": {\n        \"cloud\": \"aws\"\n    }}"
	var invalidSinkValueTag = "{\n    \"name\": \"my-prom-sink\",\n    \"backend\": \"prometheus\",\n    \"config\": {\n        \"remote_host\": \"my.prometheus-host.com\",\n        \"username\": \"dbuser\"\n    },\n    \"description\": \"An example prometheus sink\",\n    \"tags\": \"invalidTag\"}"

	cases := map[string]struct {
		req         string
		contentType string
		auth        string
		status      int
		location    string
	}{
		"validate a valid sink": {
			req:         validJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusOK,
			location:    "/sinks/validate",
		},
		"validate an invalid json": {
			req:         invalidJson,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks/validate",
		},
		"validate a sink with a empty token": {
			req:         validJson,
			contentType: contentType,
			auth:        "",
			status:      http.StatusUnauthorized,
			location:    "/sinks/validate",
		},
		"validate a sink with an invalid token": {
			req:         validJson,
			contentType: contentType,
			auth:        invalidToken,
			status:      http.StatusUnauthorized,
			location:    "/sinks/validate",
		},
		"validate a valid sink without content type": {
			req:         validJson,
			contentType: "",
			auth:        token,
			status:      http.StatusUnsupportedMediaType,
			location:    "/sinks/validate",
		},
		"validate an invalid sink field": {
			req:         invalidSinkField,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks/validate",
		},
		"validate a sink with invalid name value": {
			req:         invalidSinkValueName,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks/validate",
		},
		"validate a sink with invalid backend value": {
			req:         invalidSinkValueBackend,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks/validate",
		},
		"validate a sink with invalid tag value": {
			req:         invalidSinkValueTag,
			contentType: contentType,
			auth:        token,
			status:      http.StatusBadRequest,
			location:    "/sinks/validate",
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			req := testRequest{
				client:      server.Client(),
				method:      http.MethodPost,
				url:         fmt.Sprintf("%s/sinks/validate", server.URL),
				contentType: tc.contentType,
				token:       tc.auth,
				body:        strings.NewReader(tc.req),
			}
			res, err := req.make()
			assert.Nil(t, err, fmt.Sprintf("unexpected erro %s", err))
			assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", desc, tc.status, res.StatusCode))
		})
	}
}

func TestOmitPasswords(t *testing.T) {
	cases := map[string]struct {
		inputMetadata    types.Metadata
		expectedMetadata types.Metadata
	}{
		"omit configuration with password": {
			inputMetadata:    types.Metadata{"user": 387157, "password": "s3cr3tp@ssw0rd", "url": "someUrl"},
			expectedMetadata: types.Metadata{"user": 387157, "password": "", "url": "someUrl"},
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			metadata := omitSecretInformation(tc.inputMetadata)
			assert.Equal(t, tc.expectedMetadata, metadata)
		})
	}
}
