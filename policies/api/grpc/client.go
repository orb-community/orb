// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
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
	timeout                  time.Duration
	retrievePolicy           endpoint.Endpoint
	retrievePoliciesByGroups endpoint.Endpoint
	inactivateDataset        endpoint.Endpoint
}

func (client grpcClient) RetrievePoliciesByGroups(ctx context.Context, in *pb.PoliciesByGroupsReq, opts ...grpc.CallOption) (*pb.PolicyListRes, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	ar := accessByGroupIDReq{
		GroupIDs: in.GroupIDs,
		OwnerID:  in.OwnerID,
	}
	res, err := client.retrievePoliciesByGroups(ctx, ar)
	if err != nil {
		return nil, err
	}

	ir := res.(policyListRes)

	plist := make([]*pb.PolicyRes, len(ir.policies))
	for i, p := range ir.policies {
		plist[i] = &pb.PolicyRes{Id: p.id, Name: p.name, Data: p.data, Backend: p.backend, Version: p.version}
	}
	return &pb.PolicyListRes{Policies: plist}, nil
}

func (client grpcClient) RetrievePolicy(ctx context.Context, in *pb.PolicyByIDReq, opts ...grpc.CallOption) (*pb.PolicyRes, error) {
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
	return &pb.PolicyRes{Id: ir.id, Name: ir.name, Data: ir.data, Backend: ir.backend, Version: ir.version}, nil
}

func (client grpcClient) InactivateDataset(ctx context.Context, in *pb.DatasetByGroupReq, otps ...grpc.CallOption) (*empty.Empty, error) {
	//ctx, cancel := context.WithTimeout(ctx, client.timeout)
	ctx, cancel := context.WithTimeout(ctx, time.Second*30000000)
	defer cancel()

	ar := accessByGroupAndOwnerID{
		GroupID: in.GroupID,
		OwnerID: in.OwnerID,
	}

	_, err := client.inactivateDataset(ctx, ar)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// NewClient returns new gRPC client instance.
func NewClient(tracer opentracing.Tracer, conn *grpc.ClientConn, timeout time.Duration) pb.PolicyServiceClient {
	svcName := "policies.PolicyService"

	return &grpcClient{
		timeout: timeout,
		retrievePolicy: kitot.TraceClient(tracer, "retrieve_policy")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrievePolicy",
			encodeRetrievePolicyRequest,
			decodePolicyResponse,
			pb.PolicyRes{},
		).Endpoint()),
		retrievePoliciesByGroups: kitot.TraceClient(tracer, "retrieve_policies_by_groups")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrievePoliciesByGroups",
			encodeRetrievePoliciesByGroupsRequest,
			decodePolicyListResponse,
			pb.PolicyListRes{},
		).Endpoint()),
		inactivateDataset: kitot.TraceClient(tracer, "inactivate_dataset")(kitgrpc.NewClient(
			conn,
			svcName,
			"InactivateDataset",
			encodeInactivateDatasetRequest,
			decodePolicyResponse,
			empty.Empty{},
		).Endpoint()),
	}
}

func encodeRetrievePolicyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByIDReq)
	return &pb.PolicyByIDReq{PolicyID: req.PolicyID, OwnerID: req.OwnerID}, nil
}

func decodePolicyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.PolicyRes)
	return policyRes{id: res.GetId(), name: res.GetName(), data: res.GetData(), version: res.GetVersion(), backend: res.GetBackend()}, nil
}

func encodeRetrievePoliciesByGroupsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByGroupIDReq)
	return &pb.PoliciesByGroupsReq{GroupIDs: req.GroupIDs, OwnerID: req.OwnerID}, nil
}

func encodeInactivateDatasetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByGroupAndOwnerID)
	return &pb.DatasetByGroupReq{
		GroupID: req.GroupID,
		OwnerID: req.OwnerID,
	}, nil
}

func decodePolicyListResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.PolicyListRes)
	policies := make([]policyRes, len(res.Policies))
	for i, p := range res.Policies {
		policies[i] = policyRes{id: p.GetId(), name: p.GetName(), data: p.GetData(), version: p.GetVersion(), backend: p.GetBackend()}
	}
	return policyListRes{policies: policies}, nil
}
