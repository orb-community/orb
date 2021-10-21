// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"context"
	"github.com/ns1labs/orb/sinks/pb"
	"google.golang.org/grpc"
)

var _ pb.SinkServiceClient = (*grpcClient)(nil)

type grpcClient struct {}

func (client grpcClient) RetrieveSink(ctx context.Context, in *pb.SinkByIDReq, opts ...grpc.CallOption) (*pb.SinkRes, error) {
	return &pb.SinkRes{}, nil
}

func NewClient() pb.SinkServiceClient {
	return &grpcClient{}
}
