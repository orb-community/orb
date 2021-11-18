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
	timeout                    time.Duration
	retrievePolicy             endpoint.Endpoint
	retrievePoliciesByGroups   endpoint.Endpoint
	retrieveDataset            endpoint.Endpoint
	retrieveDatasetsByPolicyID endpoint.Endpoint
}

func (client grpcClient) RetrieveDatasetsByPolicyID(ctx context.Context, in *pb.PolicyByIDReq, opts ...grpc.CallOption) (*pb.DatasetListRes, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	ar := accessByIDReq{
		PolicyID: in.PolicyID,
		OwnerID:  in.OwnerID,
	}
	res, err := client.retrieveDatasetsByPolicyID(ctx, ar)
	if err != nil {
		return nil, err
	}

	ir := res.(datasetListRes)

	dsList := make([]*pb.DatasetRes, len(ir.datasets))
	for i, p := range ir.datasets {
		dsList[i] = &pb.DatasetRes{
			Id:           p.id,
			AgentGroupId: p.agentGroupID,
			PolicyId:     p.policyID,
			SinkIds:      p.sinkIDs,
		}
	}
	return &pb.DatasetListRes{Datasets: dsList}, nil
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

func (client grpcClient) RetrievePoliciesByGroups(ctx context.Context, in *pb.PoliciesByGroupsReq, opts ...grpc.CallOption) (*pb.PolicyInDSListRes, error) {
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

	ir := res.(policyInDSListRes)

	plist := make([]*pb.PolicyInDSRes, len(ir.policies))
	for i, p := range ir.policies {
		plist[i] = &pb.PolicyInDSRes{Id: p.id, Name: p.name, Data: p.data, Backend: p.backend, Version: p.version, DatasetId: p.datasetID}
	}
	return &pb.PolicyInDSListRes{Policies: plist}, nil
}

func (client grpcClient) RetrieveDataset(ctx context.Context, in *pb.DatasetByIDReq, opts ...grpc.CallOption) (*pb.DatasetRes, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	ar := accessDatasetByIDReq{
		datasetID: in.DatasetID,
		ownerID:   in.OwnerID,
	}
	res, err := client.retrieveDataset(ctx, ar)
	if err != nil {
		return nil, err
	}
	ir := res.(datasetRes)
	return &pb.DatasetRes{
		Id:           ir.id,
		AgentGroupId: ir.agentGroupID,
		PolicyId:     ir.policyID,
		SinkIds:      ir.sinkIDs,
	}, nil
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
			pb.PolicyInDSListRes{},
		).Endpoint()),
		retrieveDataset: kitot.TraceClient(tracer, "retrieve_dataset")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrieveDataset",
			encodeRetrieveDatasetRequest,
			decodeDatasetResponse,
			pb.DatasetRes{},
		).Endpoint()),
		retrieveDatasetsByPolicyID: kitot.TraceClient(tracer, "retrieve_datasets_by_policy_id")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrieveDatasetsByPolicyID",
			encodeRetrieveDatasetsByPolicyIDRequest,
			decodeDatasetsByPolicyIDResponse,
			pb.DatasetListRes{},
		).Endpoint()),
	}
}

func encodeRetrievePolicyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByIDReq)
	return &pb.PolicyByIDReq{PolicyID: req.PolicyID, OwnerID: req.OwnerID}, nil
}

func encodeRetrievePoliciesByGroupsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByGroupIDReq)
	return &pb.PoliciesByGroupsReq{GroupIDs: req.GroupIDs, OwnerID: req.OwnerID}, nil
}

func encodeRetrieveDatasetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessDatasetByIDReq)
	return &pb.DatasetByIDReq{
		DatasetID: req.datasetID,
		OwnerID:   req.ownerID,
	}, nil
}

func encodeRetrieveDatasetsByPolicyIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByIDReq)
	return &pb.PolicyByIDReq{PolicyID: req.PolicyID, OwnerID: req.OwnerID}, nil
}

func decodePolicyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.PolicyRes)
	return policyRes{id: res.GetId(), name: res.GetName(), data: res.GetData(), version: res.GetVersion(), backend: res.GetBackend()}, nil
}

func decodeDatasetResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.DatasetRes)
	return datasetRes{
		id:           res.GetId(),
		agentGroupID: res.GetAgentGroupId(),
		policyID:     res.GetPolicyId(),
		sinkIDs:      res.GetSinkIds(),
	}, nil
}

func decodePolicyListResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.PolicyInDSListRes)
	policies := make([]policyInDSRes, len(res.Policies))
	for i, p := range res.Policies {
		policies[i] = policyInDSRes{id: p.GetId(), name: p.GetName(), data: p.GetData(), version: p.GetVersion(), backend: p.GetBackend(), datasetID: p.GetDatasetId()}
	}
	return policyInDSListRes{policies: policies}, nil
}

func decodeDatasetsByPolicyIDResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.DatasetListRes)
	datasets := make([]datasetRes, len(res.Datasets))
	for i, p := range res.Datasets {
		datasets[i] = datasetRes{
			id:           p.GetId(),
			agentGroupID: p.GetAgentGroupId(),
			policyID:     p.GetPolicyId(),
			sinkIDs:      p.GetSinkIds(),
		}
	}
	return datasetListRes{datasets: datasets}, nil
}
