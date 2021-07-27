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
			token: "",
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
