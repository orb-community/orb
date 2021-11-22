// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"context"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/policies"
	"github.com/ns1labs/orb/policies/pb"

	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.PolicyServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	pb.UnimplementedPolicyServiceServer
	retrievePolicy             kitgrpc.Handler
	retrievePoliciesByGroups   kitgrpc.Handler
	retrieveDataset            kitgrpc.Handler
	retrieveDatasetsByPolicyID kitgrpc.Handler
}

// NewServer returns new PolicyServiceServer instance.
func NewServer(tracer opentracing.Tracer, svc policies.Service) pb.PolicyServiceServer {
	return &grpcServer{
		retrievePolicy: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_policy")(retrievePolicyEndpoint(svc)),
			decodeRetrievePolicyRequest,
			encodePolicyResponse,
		),
		retrievePoliciesByGroups: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_policies_by_groups")(retrievePoliciesByGroupsEndpoint(svc)),
			decodeRetrievePoliciesByGroupRequest,
			encodePolicyInDSListResponse,
		),
		retrieveDataset: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_dataset")(retrieveDatasetEnpoint(svc)),
			decodeRetrieveDatasetRequest,
			encodeDatasetResponse,
		),
		retrieveDatasetsByPolicyID: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_datasets_by_policy_id")(retrieveDatasetsByPolicyIDEndpoint(svc)),
			decodeRetrieveDatasetsByPolicyIDRequest,
			encodeDatasetsByPolicyIDResponse,
		),
	}
}

func (gs *grpcServer) RetrievePoliciesByGroups(ctx context.Context, req *pb.PoliciesByGroupsReq) (*pb.PolicyInDSListRes, error) {
	_, res, err := gs.retrievePoliciesByGroups.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.PolicyInDSListRes), nil
}

func (gs *grpcServer) RetrievePolicy(ctx context.Context, req *pb.PolicyByIDReq) (*pb.PolicyRes, error) {
	_, res, err := gs.retrievePolicy.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.PolicyRes), nil
}

func (gs *grpcServer) RetrieveDataset(ctx context.Context, req *pb.DatasetByIDReq) (*pb.DatasetRes, error) {
	_, res, err := gs.retrieveDataset.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.DatasetRes), nil
}

func (gs *grpcServer) RetrieveDatasetsByPolicyID(ctx context.Context, req *pb.PolicyByIDReq) (*pb.DatasetListRes, error) {
	_, res, err := gs.retrieveDatasetsByPolicyID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.DatasetListRes), nil
}

func decodeRetrievePolicyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PolicyByIDReq)
	return accessByIDReq{PolicyID: req.PolicyID, OwnerID: req.OwnerID}, nil
}

func decodeRetrievePoliciesByGroupRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PoliciesByGroupsReq)
	return accessByGroupIDReq{GroupIDs: req.GroupIDs, OwnerID: req.OwnerID}, nil
}

func decodeRetrieveDatasetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DatasetByIDReq)
	return accessDatasetByIDReq{
		datasetID: req.DatasetID,
		ownerID:   req.OwnerID,
	}, nil
}

func decodeRetrieveDatasetsByPolicyIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PolicyByIDReq)
	return accessByIDReq{PolicyID: req.PolicyID, OwnerID: req.OwnerID}, nil
}

func encodePolicyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(policyRes)
	return &pb.PolicyRes{
		Id:      res.id,
		Name:    res.name,
		Backend: res.backend,
		Version: res.version,
		Data:    res.data,
	}, nil
}

func encodePolicyInDSListResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(policyInDSListRes)

	plist := make([]*pb.PolicyInDSRes, len(res.policies))
	for i, p := range res.policies {
		plist[i] = &pb.PolicyInDSRes{Id: p.id, Name: p.name, Data: p.data, Backend: p.backend, Version: p.version, DatasetId: p.datasetID}
	}
	return &pb.PolicyInDSListRes{Policies: plist}, nil
}

func encodeDatasetResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(datasetRes)

	return &pb.DatasetRes{
		Id:           res.id,
		AgentGroupId: res.agentGroupID,
		PolicyId:     res.policyID,
		Valid:        res.valid,
		SinkIds:      res.sinkIDs,
	}, nil
}

func encodeDatasetsByPolicyIDResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(datasetListRes)

	dsList := make([]*pb.DatasetRes, len(res.datasets))
	for i, p := range res.datasets {
		dsList[i] = &pb.DatasetRes{
			Id:           p.id,
			AgentGroupId: p.agentGroupID,
			PolicyId:     p.policyID,
			Valid:        p.valid,
			SinkIds:      p.sinkIDs,
		}
	}
	return &pb.DatasetListRes{Datasets: dsList}, nil
}

func encodeError(err error) error {
	switch errorVal := err.(type) {
	case errors.Error:
		switch {
		case errors.Contains(errorVal, policies.ErrSelectEntity):
			return status.Error(codes.InvalidArgument, err.Error())
		case errors.Contains(errorVal, policies.ErrMalformedEntity):
			return status.Error(codes.InvalidArgument, "received invalid can access request")
		}
	default:
		return status.Error(codes.Internal, "internal server error")
	}
	return nil
}
