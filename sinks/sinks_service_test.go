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
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}
	wrongID, _ = uuid.NewV4()
)

func newService(tokens map[string]string) sinks.SinkService {
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

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := service.CreateSink(context.Background(), tc.token, sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, err, tc.err))
			t.Log(tc.token)
		})
	}

}

func TestUpdateSink(t *testing.T) {
	service := newService(map[string]string{token: email})
	sk, err := service.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	sk.Backend = ""
	sk.State = sinks.Unknown
	sk.Error = ""
	wrongSink := sinks.Sink{ID: wrongID.String()}
	sink.ID = sk.ID

	cases := map[string]struct {
		sink  sinks.Sink
		token string
		err   error
	}{
		"update existing sink": {
			sink:  sk,
			token: token,
			err:   nil,
		},
		"update sink with wrong credentials": {
			sink:  sink,
			token: invalidToken,
			err:   sinks.ErrUnauthorizedAccess,
		},
		"update a non-existing sink": {
			sink:  wrongSink,
			token: token,
			err:   sinks.ErrNotFound,
		},
		"update sink read only fields": {
			sink:  sink,
			token: token,
			err:   errors.ErrUpdateEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := service.UpdateSink(context.Background(), tc.token, tc.sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %d got %d", desc, tc.err, err))
		})
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

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := service.ViewSink(context.Background(), tc.token, tc.key)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestListSinks(t *testing.T) {
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

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			page, err := service.ListSinks(context.Background(), tc.token, tc.pageMetadata)
			size := uint64(len(page.Sinks))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))

			testSortSinks(t, tc.pageMetadata, page.Sinks)
		})
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

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := service.ViewBackend(context.Background(), tc.token, tc.backend)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
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

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := service.ListBackends(context.Background(), tc.token)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}

}

func TestDeleteSink(t *testing.T) {
	svc := newService(map[string]string{token: email})

	sk, err := svc.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"delete existing sink": {
			id:    sk.ID,
			token: token,
			err:   nil,
		},
		"delete non-existent sink": {
			id:    wrongID.String(),
			token: token,
			err:   nil,
		},
		"delete sink with wrong credentials": {
			id:    sk.ID,
			token: invalidToken,
			err:   sinks.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := svc.DeleteSink(context.Background(), tc.token, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestValidateSink(t *testing.T) {
	service := newService(map[string]string{token: email})

	cases := map[string]struct {
		sink  sinks.Sink
		token string
		err   error
	}{
		"validate a new sink": {
			sink:  sink,
			token: token,
			err:   nil,
		},
		"validate a sink with a invalid token": {
			sink:  sink,
			token: invalidToken,
			err:   sinks.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := service.ValidateSink(context.Background(), tc.token, sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, err, tc.err))
		})
	}

}

func TestSinksStatisticsSummary(t *testing.T) {
	service := newService(map[string]string{token: email})

	var (
		sks []sinks.Sink
		err error
	)
	for i := 0; i < 10; i++ {
		sink.Name, err = types.NewIdentifier(fmt.Sprintf("my-sink-%d", i))
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		sink.State = sinks.Active
		sk, err := service.CreateSink(context.Background(), token, sink)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		sks = append(sks, sk)
	}

	var skSummary []sinks.SinkStates
	skSummary = append(skSummary, sinks.SinkStates{
		State: sinks.Active,
		Count: 10,
	})

	cases := map[string]struct{
		token string
		statistics sinks.SinksStatistics
		err error
	}{
		"retrieve all sinks statistics summary": {
			token: token,
			statistics: sinks.SinksStatistics{
				StatesSummary: skSummary,
				TotalSinks: 10,
			},
			err: nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			statistics, err := service.SinksStatistics(context.Background(), tc.token)
			assert.Equal(t, statistics, tc.statistics, fmt.Sprintf("%s: expected %+v got %+v", desc, tc.statistics, statistics))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
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
