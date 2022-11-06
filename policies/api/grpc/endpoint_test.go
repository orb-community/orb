// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc_test

import (
	"context"
	"fmt"
	"github.com/etaques/orb/policies/pb"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go/mocktracer"

	policiesgrpc "github.com/etaques/orb/policies/api/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRetrievePolicy(t *testing.T) {

	usersAddr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(usersAddr, grpc.WithInsecure())
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	cli := policiesgrpc.NewClient(mocktracer.New(), conn, time.Second*5)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cases := map[string]struct {
		id   string
		code codes.Code
	}{
		"retrieve existing policy": {
			id:   policy.ID,
			code: codes.OK,
		},
		//"retrieve non-existent policy": {
		//	id:   "nonexist",
		//	code: codes.NotFound,
		//},
	}

	for desc, tc := range cases {
		id, err := cli.RetrievePolicy(ctx, &pb.PolicyByIDReq{
			PolicyID: policy.ID,
			OwnerID:  policy.MFOwnerID,
		})
		e, ok := status.FromError(err)
		assert.True(t, ok, "OK expected to be true")
		assert.Equal(t, tc.id, id.GetId(), fmt.Sprintf("%s: expected %s got %s", desc, tc.id, id.GetId()))
		assert.Equal(t, tc.code, e.Code(), fmt.Sprintf("%s: expected %s got %s", desc, tc.code, e.Code()))
	}
}

func TestRetrievePoliciesByGroups(t *testing.T) {

	usersAddr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(usersAddr, grpc.WithInsecure())
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	cli := policiesgrpc.NewClient(mocktracer.New(), conn, time.Second*5)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cases := map[string]struct {
		ids     []string
		results int
		code    codes.Code
	}{
		"retrieve existing policy by group": {
			ids:     []string{dataset.AgentGroupID},
			results: 1,
			code:    codes.OK,
		},
		//"retrieve non-existent policy": {
		//	id:   "nonexist",
		//	code: codes.NotFound,
		//},
	}

	for desc, tc := range cases {
		plist, err := cli.RetrievePoliciesByGroups(ctx, &pb.PoliciesByGroupsReq{
			GroupIDs: tc.ids,
			OwnerID:  policy.MFOwnerID,
		})
		e, ok := status.FromError(err)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
		assert.True(t, ok, "OK expected to be true")
		assert.Equal(t, tc.results, len(plist.Policies), fmt.Sprintf("%s: expected %d got %d", desc, tc.results, len(plist.Policies)))
		assert.Equal(t, tc.code, e.Code(), fmt.Sprintf("%s: expected %s got %s", desc, tc.code, e.Code()))
	}
}
