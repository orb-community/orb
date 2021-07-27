// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	thmocks "github.com/mainflux/mainflux/things/mocks"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	skmocks "github.com/ns1labs/orb/sinks/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

const (
	contentType  = "application/json"
	token        = "token"
	invalidToken = "invalid"
	email        = "user@example.com"
	n            = uint64(10)
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
	wrongID, _ = uuid.NewV4()
)

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

func TestCreateSink(t *testing.T) {
	service := newService(map[string]string{token: email})

	cases := map[string]struct {
		sink  sinks.Sink
		token string
		err   error
	}{
		"create a new sink": {
			sink:  sink,
			token: token,
			err:   nil,
		},
		"add a sink with a invalid token": {
			sink:  sink,
			token: "invalid",
			err:   sinks.ErrUnauthorizedAccess,
		},
	}

	for desc, sinkCase := range cases {
		_, err := service.CreateSink(context.Background(), sinkCase.token, sink)
		assert.True(t, errors.Contains(err, sinkCase.err), fmt.Sprintf("%s: expected %s got %s", desc, err, sinkCase.err))
	}

}

func TestUpdateSink(t *testing.T) {
	service := newService(map[string]string{token: email})

	sk, err := service.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	wrongSink := sinks.Sink{ID: wrongID.String()}
	sink.ID = sk.ID

	cases := map[string]struct {
		sink  sinks.Sink
		token string
		err   error
	}{
		"update existing sink": {
			sink:  sink,
			token: token,
			err:   nil,
		},
		"update sink with wrong credentials": {
			sink:  sink,
			token: invalidToken,
			err:   sinks.ErrUnauthorizedAccess,
		},
		"update a non-existing thing": {
			sink:  wrongSink,
			token: token,
			err:   sinks.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		err := service.UpdateSink(context.Background(), tc.token, tc.sink)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %d got %d", desc, tc.err, err))
	}
}

func TestViewSink(t *testing.T) {
	service := newService(map[string]string{token: email})

	sk, err := service.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		key   string
		token string
		err   error
	}{
		"view a existing sink": {
			key:   sk.ID,
			token: token,
			err:   nil,
		},
		"view a existing sink with wrong credentials": {
			key:   sk.ID,
			token: invalidToken,
			err:   sinks.ErrUnauthorizedAccess,
		},
		"view a non-existing sink": {
			key:   wrongID.String(),
			token: token,
			err:   sinks.ErrNotFound,
		},
	}

	for desc, sinkCase := range cases {
		_, err := service.ViewSink(context.Background(), sinkCase.token, sinkCase.key)
		assert.True(t, errors.Contains(err, sinkCase.err), fmt.Sprintf("%s: expected %s got %s\n", desc, sinkCase.err, err))
	}
}

func TestListThings(t *testing.T) {
	service := newService(map[string]string{token: email})
	metadata := make(map[string]interface{})
	metadata["serial"] = "12345"
	var sks []sinks.Sink
	for i := uint64(0); i < n; i++ {
		sink.Name, _ = types.NewIdentifier(fmt.Sprintf("my-sink-%d", i))
		sk, err := service.CreateSink(context.Background(), token, sink)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		sks = append(sks, sk)
	}

	cases := map[string]struct {
		token        string
		pageMetadata sinks.PageMetadata
		size         uint64
		err          error
	}{
		"list all sinks": {
			token: token,
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
			},
			size: n,
			err:  nil,
		},
		"list half of sinks": {
			token: token,
			pageMetadata: sinks.PageMetadata{
				Offset: n / 2,
				Limit:  n,
			},
			size: n / 2,
			err:  nil,
		},
		"list last sinks": {
			token: token,
			pageMetadata: sinks.PageMetadata{
				Offset: n - 1,
				Limit:  n,
			},
			size: 1,
			err:  nil,
		},
		"list empty set": {
			token: token,
			pageMetadata: sinks.PageMetadata{
				Offset: n + 1,
				Limit:  n,
			},
			size: 0,
			err:  nil,
		},
		"list with zero limit": {
			token: token,
			pageMetadata: sinks.PageMetadata{
				Offset: 1,
				Limit:  0,
			},
			size: 0,
			err:  nil,
		},
		"list sinks with wrong credentials": {
			token: invalidToken,
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
			},
			size: 0,
			err:  sinks.ErrUnauthorizedAccess,
		},
		"list sinks with metadata": {
			token: token,
			pageMetadata: sinks.PageMetadata{
				Offset:   0,
				Limit:    n,
				Metadata: metadata,
			},
			size: n,
			err:  nil,
		},
		"list all sinks sorted by name asc": {
			token: token,
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Order:  "name",
				Dir:    "asc",
			},
			size: n,
			err:  nil,
		},
		"list all sinks sorted by name desc": {
			token: token,
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Order:  "name",
				Dir:    "desc",
			},
			size: n,
			err:  nil,
		},
	}

	for desc, sinkCase := range cases {
		page, err := service.ListSinks(context.Background(), sinkCase.token, sinkCase.pageMetadata)
		size := uint64(len(page.Sinks))
		assert.Equal(t, sinkCase.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, sinkCase.size, size))
		assert.True(t, errors.Contains(err, sinkCase.err), fmt.Sprintf("%s: expected %s got %s", desc, sinkCase.err, err))

		testSortSinks(t, sinkCase.pageMetadata, page.Sinks)
	}

}

func TestViewBackends(t *testing.T) {
	service := newService(map[string]string{token: email})

	cases := map[string]struct {
		token   string
		backend string
		err     error
	}{
		"view a existing backend": {
			token:   token,
			backend: "prometheus",
			err:     nil,
		},
		"view a non-existing backend": {
			token:   token,
			backend: "grafana",
			err:     sinks.ErrNotFound,
		},
		"view sinks with wrong credentials": {
			token:   invalidToken,
			backend: "prometheus",
			err:     sinks.ErrUnauthorizedAccess,
		},
	}

	for desc, sinkCase := range cases {
		_, err := service.ViewBackend(context.Background(), sinkCase.token, sinkCase.backend)
		assert.True(t, errors.Contains(err, sinkCase.err), fmt.Sprintf("%s: expected %s got %s", desc, sinkCase.err, err))
	}

}

func TestListBackends(t *testing.T) {
	service := newService(map[string]string{token: email})

	cases := map[string]struct {
		token string
		err   error
	}{
		"list all backends": {
			token: token,
			err:   nil,
		},
		"list backends with wrong credentials": {
			token: invalidToken,
			err:   sinks.ErrUnauthorizedAccess,
		},
	}

	for desc, sinkCase := range cases {
		_, err := service.ListBackends(context.Background(), sinkCase.token)
		assert.True(t, errors.Contains(err, sinkCase.err), fmt.Sprintf("%s: expected %s got %s", desc, sinkCase.err, err))
	}

}

func testSortSinks(t *testing.T, pm sinks.PageMetadata, sks []sinks.Sink) {
	switch pm.Order {
	case "name":
		current := sks[0]
		for _, res := range sks {
			if pm.Dir == "asc" {
				assert.GreaterOrEqual(t, res.Name.String(), current.Name.String())
			}
			if pm.Dir == "desc" {
				assert.GreaterOrEqual(t, current.Name.String(), res.Name.String())
			}
			current = res
		}
	default:
		break
	}
}
