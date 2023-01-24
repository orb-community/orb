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
	nameID, _   = types.NewIdentifier("my-sink")
	description = "An example prometheus sink"
	sink        = sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}
	wrongID, _ = uuid.NewV4()
)

func newService(tokens map[string]string) sinks.SinkService {
	auth := thmocks.NewAuthService(tokens, make(map[string][]thmocks.MockSubjectSet))
	sinkRepo := skmocks.NewSinkRepository()
	var logger *zap.Logger

	config := mfsdk.Config{
		ThingsURL: "localhost",
	}

	mfsdk := mfsdk.NewSDK(config)
	pwdSvc := sinks.NewPasswordService(logger, "_testing_string_")
	return sinks.NewSinkService(logger, auth, sinkRepo, mfsdk, pwdSvc)
}

func TestCreateSink(t *testing.T) {
	service := newService(map[string]string{token: email})

	description := "An example prometheus sink"
	var invalidBackendSink = sinks.Sink{
		Name:        nameID,
		Description: &description,
		Backend:     "invalid",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

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
		"create a sink with a invalid backend": {
			sink:  invalidBackendSink,
			token: token,
			err:   sinks.ErrInvalidBackend,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := service.CreateSink(context.Background(), tc.token, tc.sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
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

	noConfig := sk
	noConfig.Config = make(map[string]interface{})

	description := "An example prometheus sink"
	newDescription := "new description"

	nameTestConfigAttribute, _ := types.NewIdentifier("configSink")
	sinkTestConfigAttribute, err := service.CreateSink(context.Background(), token, sinks.Sink{
		Name:        nameTestConfigAttribute,
		Description: &description,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	nameTestDescriptionAttribute, _ := types.NewIdentifier("emptyDescSink")
	sinkTestDescriptionAttribute, err := service.CreateSink(context.Background(), token, sinks.Sink{
		Name:        nameTestDescriptionAttribute,
		Description: &description,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	emptyTagsSinkName, _ := types.NewIdentifier("emptyTagsSink")
	emptyTagsSink, err := service.CreateSink(context.Background(), token, sinks.Sink{
		Name:        emptyTagsSinkName,
		Description: &description,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	addTagsToSinkName, _ := types.NewIdentifier("addTagsToSinkName")
	addTagsToSink, err := service.CreateSink(context.Background(), token, sinks.Sink{
		Name:        addTagsToSinkName,
		Description: &description,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		incomingSink sinks.Sink
		expectedSink sinks.Sink
		token        string
		err          error
	}{
		"update existing sink": {
			incomingSink: sk,
			expectedSink: sk,
			token:        token,
			err:          nil,
		},
		"update sink with wrong credentials": {
			incomingSink: sink,
			token:        invalidToken,
			err:          sinks.ErrUnauthorizedAccess,
		},
		"update a non-existing sink": {
			incomingSink: wrongSink,
			token:        token,
			err:          sinks.ErrNotFound,
		},
		"update sink read only fields": {
			incomingSink: sink,
			token:        token,
			err:          errors.ErrUpdateEntity,
		},
		"update existing sink - only updating config": {
			incomingSink: sinks.Sink{
				ID: sinkTestConfigAttribute.ID,
				Config: types.Metadata{
					"test": "config",
				},
				Error: "",
			},
			expectedSink: sinks.Sink{
				Name: sinkTestConfigAttribute.Name,
				Config: types.Metadata{
					"test": "config",
				},
				Description: sinkTestConfigAttribute.Description,
				Tags:        sinkTestConfigAttribute.Tags,
			},
			token: token,
			err:   nil,
		},
		"update existing sink - omitted config": {
			incomingSink: sinks.Sink{
				ID:    sink.ID,
				Error: "",
			},
			expectedSink: sinks.Sink{
				Name:        sink.Name,
				Config:      sink.Config,
				Description: sink.Description,
				Tags:        sink.Tags,
			},
			token: token,
			err:   nil,
		},
		"update existing sink using empty tags": {
			incomingSink: sinks.Sink{
				ID:    emptyTagsSink.ID,
				Error: "",
				Tags:  make(map[string]string),
			},
			expectedSink: sinks.Sink{
				Name:        emptyTagsSink.Name,
				Config:      emptyTagsSink.Config,
				Description: emptyTagsSink.Description,
				Tags:        make(map[string]string),
			},
			token: token,
			err:   nil,
		},
		"update existing sink - only updating description": {
			incomingSink: sinks.Sink{
				ID:          sinkTestDescriptionAttribute.ID,
				Description: &newDescription,
			},
			expectedSink: sinks.Sink{
				Name:        sinkTestDescriptionAttribute.Name,
				Description: &newDescription,
				Tags:        sinkTestDescriptionAttribute.Tags,
				Config:      sinkTestDescriptionAttribute.Config,
			},
			token: token,
			err:   nil,
		},
		"update existing sink - omitted description": {
			incomingSink: sinks.Sink{
				ID: sink.ID,
			},
			expectedSink: sinks.Sink{
				Name:        sink.Name,
				Description: sink.Description,
				Tags:        sink.Tags,
				Config:      sink.Config,
			},
			token: token,
			err:   nil,
		},
		"update sink tags with new tags": {
			incomingSink: sinks.Sink{
				ID:   addTagsToSink.ID,
				Tags: map[string]string{"cloud": "aws", "test": "true"},
			},
			expectedSink: sinks.Sink{
				Name:        addTagsToSink.Name,
				Config:      addTagsToSink.Config,
				Description: addTagsToSink.Description,
				Tags:        types.Tags{"cloud": "aws", "test": "true"},
			},
			token: token,
			err:   nil,
		},
		"update sink tags with omitted tags": {
			incomingSink: sinks.Sink{
				ID: sink.ID,
			},
			expectedSink: sinks.Sink{
				Name:        sink.Name,
				Config:      sink.Config,
				Description: sink.Description,
				Tags:        sink.Tags,
			},
			token: token,
			err:   nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			res, err := service.UpdateSink(context.Background(), tc.token, tc.incomingSink)
			if err == nil {
				assert.Equal(t, tc.expectedSink.Config, res.Config, fmt.Sprintf("%s: expected %s got %s", desc, tc.expectedSink.Config, res.Config))
				assert.Equal(t, tc.expectedSink.Name.String(), res.Name.String(), fmt.Sprintf("%s: expected name %s got %s", desc, tc.expectedSink.Name.String(), res.Name.String()))
				assert.Equal(t, *tc.expectedSink.Description, *res.Description, fmt.Sprintf("%s: expected description %s got %s", desc, *tc.expectedSink.Description, *res.Description))
				assert.Equal(t, tc.expectedSink.Tags, res.Tags, fmt.Sprintf("%s: expected tags %s got %s", desc, tc.expectedSink.Tags, res.Tags))
			}
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

	description := "An example prometheus sink"

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
		"validate a sink with a invalid backend": {
			sink: sinks.Sink{
				Name:        nameID,
				Description: &description,
				Backend:     "invalid",
				Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
				Tags:        map[string]string{"cloud": "aws"},
			},
			token: token,
			err:   sinks.ErrValidateSink,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := service.ValidateSink(context.Background(), tc.token, tc.sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}

}

func TestViewSinkInternal(t *testing.T) {
	service := newService(map[string]string{token: email})

	sk, err := service.CreateSink(context.Background(), token, sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		key     string
		ownerID string
		err     error
	}{
		"view a existing sink": {
			key:     sk.ID,
			ownerID: sk.MFOwnerID,
			err:     nil,
		},
		"view a existing sink with wrong credentials": {
			key:     sk.ID,
			ownerID: "invalid",
			err:     sinks.ErrNotFound,
		},
		"view a non-existing sink": {
			key:     wrongID.String(),
			ownerID: sk.MFOwnerID,
			err:     sinks.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := service.ViewSinkInternal(context.Background(), tc.ownerID, tc.key)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
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
