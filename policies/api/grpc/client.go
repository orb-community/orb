// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/ns1labs/orb/policies/pb"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var _ pb.PolicyServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	timeout        time.Duration
	retrievePolicy endpoint.Endpoint
}

func (client grpcClient) RetrievePolicyDataByGroups(ctx context.Context, in *pb.PolicyByGroupsReq, opts ...grpc.CallOption) (*pb.PolicyDataListRes, error) {

}

func (client grpcClient) RetrievePolicyData(ctx context.Context, in *pb.PolicyByIDReq, opts ...grpc.CallOption) (*pb.PolicyDataRes, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	ar := accessByIDReq{
		PolicyID: in.PolicyID,
		OwnerID:  in.OwnerID,
	}
	res, err := client.retrievePolicy(ctx, ar)
	if err != nil {
		return nil, err
	}

	ir := res.(policyRes)
	return &pb.PolicyDataRes{Id: ir.id, Name: ir.name, Data: ir.data, Backend: ir.backend, Version: ir.version}, nil
}

// NewClient returns new gRPC client instance.
func NewClient(tracer opentracing.Tracer, conn *grpc.ClientConn, timeout time.Duration) pb.PolicyServiceClient {
	svcName := "policies.PolicyService"

	return &grpcClient{
		timeout: timeout,
		retrievePolicy: kitot.TraceClient(tracer, "retrieve_policy")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrievePolicyData",
			encodeRetrievePolicyRequest,
			decodePolicyResponse,
			pb.PolicyDataRes{},
		).Endpoint()),
	}
}

func encodeRetrievePolicyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByIDReq)
	return &pb.PolicyByIDReq{PolicyID: req.PolicyID, OwnerID: req.OwnerID}, nil
}

func decodePolicyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.PolicyDataRes)
	return policyRes{id: res.GetId(), name: res.GetName(), data: res.GetData(), version: res.GetVersion(), backend: res.GetBackend()}, nil
}
