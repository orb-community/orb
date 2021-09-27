// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet/pb"
	"google.golang.org/grpc"
)

var _ pb.FleetServiceClient = (*fleetGrpcClientMock)(nil)

type fleetGrpcClientMock struct {}

func (g fleetGrpcClientMock) RetrieveAgent(ctx context.Context, in *pb.AgentByIDReq, opts ...grpc.CallOption) (*pb.AgentRes, error) {
	return &pb.AgentRes{}, nil
}

func (g fleetGrpcClientMock) RetrieveAgentGroup(ctx context.Context, in *pb.AgentGroupByIDReq, opts ...grpc.CallOption) (*pb.AgentGroupRes, error) {
	return &pb.AgentGroupRes{}, nil
}

func NewClient() pb.FleetServiceClient {
	return &fleetGrpcClientMock{}
}
