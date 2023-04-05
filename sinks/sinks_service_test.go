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
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"github.com/orb-community/orb/sinks/authentication_type"
	skmocks "github.com/orb-community/orb/sinks/mocks"
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
		Config: map[string]interface{}{
			"exporter":       map[string]interface{}{"remote_host": "https://orb.community/"},
			"authentication": map[string]interface{}{"type": "basicauth", "username": "dbuser", "password": "dbpass"},
		},
		Tags: map[string]string{"cloud": "aws"},
	}
	wrongID, _ = uuid.NewV4()
)

func newService(tokens map[string]string) sinks.SinkService {
	logger := zap.NewNop()
	auth := thmocks.NewAuthService(tokens, make(map[string][]thmocks.MockSubjectSet))
	pwdSvc := authentication_type.NewPasswordService(logger, "_testing_string_")
	sinkRepo := skmocks.NewSinkRepository(pwdSvc)

	config := mfsdk.Config{
		ThingsURL: "localhost",
	}

	newSDK := mfsdk.NewSDK(config)
	return sinks.NewSinkService(logger, auth, sinkRepo, newSDK, pwdSvc)
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
		Config: map[string]interface{}{
			"exporter":       map[string]interface{}{"remote_host": "https://orb.community/"},
			"authentication": map[string]interface{}{"type": "basicauth", "username": "dbuser", "password": "dbpass"},
		},
		Tags: map[string]string{"cloud": "aws"},
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

func TestIdempotencyUpdateSink(t *testing.T) {
	ctx := context.Background()
	service := newService(map[string]string{token: email})
	jsonSinkName, err := types.NewIdentifier("initial-json-Sink")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	yamlSinkName, err := types.NewIdentifier("initial-yaml-Sink")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	aInitialDescription := "A initial description worthy reading"
	initialJsonSink := sinks.Sink{
		Name:        jsonSinkName,
		Description: &aInitialDescription,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config: map[string]interface{}{
			"exporter":       map[string]interface{}{"remote_host": "https://orb.community/"},
			"authentication": map[string]interface{}{"type": "basicauth", "username": "dbuser", "password": "dbpass"},
		},
		Tags: map[string]string{"cloud": "aws"},
	}
	initialYamlSink := sinks.Sink{
		Name:        yamlSinkName,
		Description: &aInitialDescription,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		ConfigData:  "exporter: \n    remote_host: https://orb.community/\nauthentication:\n    type: basicauth\n    username: dbuser\n    password: dbpass\n",
		Format:      "yaml",
		Config: map[string]interface{}{
			"exporter":       map[string]interface{}{"remote_host": "https://orb.community/"},
			"authentication": map[string]interface{}{"type": "basicauth", "username": "dbuser", "password": "dbpass"},
		},
		MFOwnerID: "OrbCommunity",
		Tags:      map[string]string{"cloud": "aws"},
	}
	jsonCreatedSink, err := service.CreateSink(ctx, token, initialJsonSink)
	require.NoError(t, err, "failed to create entity")
	require.NotEmptyf(t, jsonCreatedSink.ID, "id must not be empty")
	yamlCreatedSink, err := service.CreateSink(ctx, token, initialYamlSink)
	require.NoError(t, err, "failed to create entity")
	initialJsonSink.ID = jsonCreatedSink.ID
	initialYamlSink.ID = yamlCreatedSink.ID
	var cases = map[string]struct {
		name        string
		requestSink sinks.Sink
		expected    func(t *testing.T, value sinks.Sink, err error)
		token       string
	}{
		"idempotency json update": {
			requestSink: initialJsonSink,
			expected: func(t *testing.T, value sinks.Sink, err error) {
				require.NoError(t, err, "no error expected")
				require.NotNilf(t, value.Description, "description is nil")
				desc := *value.Description
				require.Equal(t, desc, aInitialDescription, "description is not equal")
				require.Equal(t, value.Name, jsonSinkName, "sink name is not equal")
				tagVal, tagOk := value.Tags["cloud"]
				require.True(t, tagOk)
				require.Equal(t, "aws", tagVal)
				require.Equalf(t, "https://orb.community/", value.Config["remote_host"], "remote host is not equal")
				require.Equalf(t, "netops", value.Config["username"], "username is not equal")
			},
			token: token,
		},
		"idempotency yaml update": {
			requestSink: initialYamlSink,
			expected: func(t *testing.T, value sinks.Sink, err error) {
				require.NoError(t, err, "no error expected")
				require.NotNilf(t, value.Description, "description is nil")
				desc := *value.Description
				require.Equal(t, desc, aInitialDescription, "description is not equal")
				require.Equal(t, value.Name, yamlSinkName, "sink name is not equal")
				tagVal, tagOk := value.Tags["cloud"]
				require.True(t, tagOk)
				require.Equal(t, "aws", tagVal)
				require.Equalf(t, "https://orb.community/", value.Config["remote_host"], "remote host is not equal")
			},
			token: token,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			res, err := service.UpdateSink(ctx, tc.token, tc.requestSink)
			tc.expected(t, res, err)
		})
	}
}

func TestPartialUpdateSink(t *testing.T) {
	ctx := context.Background()
	service := newService(map[string]string{token: email})
	jsonSinkName, err := types.NewIdentifier("initial-json-Sink")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	yamlSinkName, err := types.NewIdentifier("initial-yaml-Sink")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	//newSinkName, err := types.NewIdentifier("updated-Sink")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	//aNewDescription := "A new description worthy reading"
	aInitialDescription := "A initial description worthy reading"
	initialJsonSink := sinks.Sink{
		Name:        jsonSinkName,
		Description: &aInitialDescription,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config: map[string]interface{}{
			"exporter":       map[string]interface{}{"remote_host": "https://orb.community/"},
			"authentication": map[string]interface{}{"type": "basicauth", "username": "dbuser", "password": "dbpass"},
		},
		Tags: map[string]string{"cloud": "aws"},
	}
	initialUsername := "netops"
	initialPassword := "w0w-orb-Rocks!"
	initialYamlSink := sinks.Sink{
		Name:        yamlSinkName,
		Description: &aInitialDescription,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		ConfigData:  "remote_host:https://orb.community/\nusername: netops\npassword: w0w-orb-Rocks!",
		Format:      "yaml",
		MFOwnerID:   "OrbCommunity",
		Config:      map[string]interface{}{"remote_host": "https://orb.community/", "username": &initialUsername, "password": &initialPassword},
		Tags:        map[string]string{"cloud": "aws"},
	}
	jsonCreatedSink, err := service.CreateSink(ctx, token, initialJsonSink)
	require.NoError(t, err, "failed to create entity")
	require.NotEmptyf(t, jsonCreatedSink.ID, "id must not be empty")
	yamlCreatedSink, err := service.CreateSink(ctx, token, initialYamlSink)
	require.NoError(t, err, "failed to create entity")
	initialJsonSink.ID = jsonCreatedSink.ID
	initialYamlSink.ID = yamlCreatedSink.ID
	var cases = map[string]struct {
		name        string
		requestSink sinks.Sink
		expected    func(t *testing.T, value sinks.Sink, err error)
		token       string
	}{
		// TODO this will fail locally because of password encryption,
		// TODO we will revisit this whenever there is an update on password encryption
		//"update only name": {
		//	requestSink: sinks.Sink{
		//		ID:   jsonCreatedSink.ID,
		//		Name: newSinkName,
		//	},
		//	expected: func(t *testing.T, value sinks.Sink, err error) {
		//		require.NoError(t, err, "no error expected")
		//		require.Equal(t, value.Name, newSinkName, "sink name is not equal")
		//	},
		//	token: token,
		//},
		//"update only description": {
		//	requestSink: sinks.Sink{
		//		ID:          jsonCreatedSink.ID,
		//		Description: &aNewDescription,
		//	},
		//	expected: func(t *testing.T, value sinks.Sink, err error) {
		//		require.NoError(t, err, "no error expected")
		//		require.NotNilf(t, value.Description, "description is nil")
		//		desc := *value.Description
		//		require.Equal(t, desc, aNewDescription, "description is not equal")
		//	},
		//	token: token,
		//}, "update only tags": {
		//	requestSink: sinks.Sink{
		//		ID:   jsonCreatedSink.ID,
		//		Tags: map[string]string{"cloud": "gcp", "from_aws": "true"},
		//	},
		//	expected: func(t *testing.T, value sinks.Sink, err error) {
		//		require.NoError(t, err, "no error expected")
		//		tagVal, tagOk := value.Tags["cloud"]
		//		tag2Val, tag2Ok := value.Tags["from_aws"]
		//		require.True(t, tagOk)
		//		require.Equal(t, "gcp", tagVal)
		//		require.True(t, tag2Ok)
		//		require.Equal(t, "true", tag2Val)
		//	},
		//	token: token,
		//},
		"update config json": {
			requestSink: sinks.Sink{
				ID:     jsonCreatedSink.ID,
				Config: map[string]interface{}{"remote_host": "https://orb.community/prom/push", "username": "netops_admin", "password": "w0w-orb-Rocks!"},
			},
			expected: func(t *testing.T, value sinks.Sink, err error) {
				require.NoError(t, err, "no error expected")
				require.Equalf(t, "https://orb.community/prom/push", value.Config["remote_host"], "want %s, got %s", "https://orb.community/prom/push", value.Config["remote_host"])
				require.Equalf(t, "netops_admin", value.Config["username"], "want %s, got %s", "netops_admin", value.Config["username"])
			},
			token: token,
		}, "update config yaml": {
			requestSink: sinks.Sink{
				ID:         yamlCreatedSink.ID,
				Format:     "yaml",
				ConfigData: "remote_host: https://orb.community/prom/push2\nusername: netops_admin\npassword: \"w0w-orb-Rocks!\"",
			},
			expected: func(t *testing.T, value sinks.Sink, err error) {
				require.NoError(t, err, "no error expected")
				require.Equalf(t, "https://orb.community/prom/push2", value.Config["remote_host"], "want %s, got %s", "https://orb.community/prom/push2", value.Config["remote_host"])
			},
			token: token,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			res, err := service.UpdateSink(ctx, tc.token, tc.requestSink)
			tc.expected(t, res, err)
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
		Config: map[string]interface{}{
			"exporter":       map[string]interface{}{"remote_host": "https://orb.community/"},
			"authentication": map[string]interface{}{"type": "basicauth", "username": "dbuser", "password": "dbpass"},
		},
		Tags: map[string]string{"cloud": "aws"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	nameTestDescriptionAttribute, _ := types.NewIdentifier("emptyDescSink")
	sinkTestDescriptionAttribute, err := service.CreateSink(context.Background(), token, sinks.Sink{
		Name:        nameTestDescriptionAttribute,
		Description: &description,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config: map[string]interface{}{
			"exporter":       map[string]interface{}{"remote_host": "https://orb.community/"},
			"authentication": map[string]interface{}{"type": "basicauth", "username": "dbuser", "password": "dbpass"},
		},
		Tags: map[string]string{"cloud": "aws"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	emptyTagsSinkName, _ := types.NewIdentifier("emptyTagsSink")
	emptyTagsSink, err := service.CreateSink(context.Background(), token, sinks.Sink{
		Name:        emptyTagsSinkName,
		Description: &description,
		Backend:     "prometheus",
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "https://orb.community/", "username": "dbuser"},
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
		Config:      map[string]interface{}{"remote_host": "https://orb.community/", "username": "dbuser"},
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
		"update existing sink - only updating config": {
			incomingSink: sinks.Sink{
				ID: sinkTestConfigAttribute.ID,
				Config: types.Metadata{
					"remote_host": "https://orb.community/",
				},
				Error: "",
			},
			expectedSink: sinks.Sink{
				Name: sinkTestConfigAttribute.Name,
				Config: types.Metadata{
					"opentelemetry": "enabled", "remote_host": "https://orb.community/",
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
				assert.Equal(t, tc.expectedSink.Config, res.Config, "config not as expected")
				assert.Equal(t, tc.expectedSink.Name.String(), res.Name.String(), "sink name not as expected")
				assert.Equal(t, *tc.expectedSink.Description, *res.Description, "sink description not as expected")
				assert.Equal(t, tc.expectedSink.Tags, res.Tags, "sink tags not as expected")
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
				Config: map[string]interface{}{
					"exporter":       map[string]interface{}{"remote_host": "https://orb.community/"},
					"authentication": map[string]interface{}{"type": "basicauth", "username": "dbuser", "password": "dbpass"},
				},
				Tags: map[string]string{"cloud": "aws"},
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
