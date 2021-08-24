// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http_test

import (
	"encoding/json"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/things"
	thingsapi "github.com/mainflux/mainflux/things/api/things/http"
	"github.com/ns1labs/orb/fleet"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/ns1labs/orb/policies"
	api "github.com/ns1labs/orb/policies/api/http"
	"github.com/opentracing/opentracing-go/mocktracer"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
)

const (
	email       = "user@example.com"
	channelsNum = 3
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
	mux := thingsapi.MakeHandler(mocktracer.New(), svc)
	return httptest.NewServer(mux)
}

func newService(auth mainflux.AuthServiceClient, url string) policies.Service {
	return policies.New(auth, nil)
}

func newServer(svc policies.Service) *httptest.Server {
	mux := api.MakeHandler(mocktracer.New(), "policies", svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
