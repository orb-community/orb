// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"time"
)

var _ pb.FleetServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	timeout                  time.Duration
	retrieveAgent            endpoint.Endpoint
	retrieveAgentGroup       endpoint.Endpoint
	retrieveOwnerByChannelID endpoint.Endpoint
}

func (g grpcClient) RetrieveAgent(ctx context.Context, in *pb.AgentByIDReq, opts ...grpc.CallOption) (*pb.AgentRes, error) {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	ar := accessByIDReq{
		AgentID: in.AgentID,
		OwnerID: in.OwnerID,
	}
	res, err := g.retrieveAgent(ctx, ar)
	if err != nil {
		return nil, err
	}

	ir := res.(agentRes)
	return &pb.AgentRes{Id: ir.id, Name: ir.name, Channel: ir.channel}, nil
}

func (g grpcClient) RetrieveAgentGroup(ctx context.Context, in *pb.AgentGroupByIDReq, opts ...grpc.CallOption) (*pb.AgentGroupRes, error) {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	ar := accessAgByIDReq{
		AgentGroupID: in.AgentGroupID,
		OwnerID:      in.OwnerID,
	}
	res, err := g.retrieveAgentGroup(ctx, ar)
	if err != nil {
		return nil, err
	}

	ir := res.(agentGroupRes)
	return &pb.AgentGroupRes{Id: ir.id, Name: ir.name, Channel: ir.channel}, nil
}

func (g grpcClient) RetrieveOwnerByChannelID(ctx context.Context, in *pb.OwnerByChannelIDReq, opts ...grpc.CallOption) (*pb.OwnerRes, error) {
	//ctx, cancel := context.WithTimeout(ctx, g.timeout)
	ctx, cancel := context.WithTimeout(ctx, 10000*time.Second)
	defer cancel()

	ar := accessOwnerByChannelIDReq{ChannelID: in.Channel}

	res, err := g.retrieveOwnerByChannelID(ctx, ar)
	if err != nil {
		return nil, err
	}

	ir := res.(ownerRes)
	return &pb.OwnerRes{OwnerID: ir.ownerID}, nil
}

// NewClient returns new gRPC client instance.
func NewClient(tracer opentracing.Tracer, conn *grpc.ClientConn, timeout time.Duration) pb.FleetServiceClient {
	svcName := "fleet.FleetService"

	return &grpcClient{
		timeout: timeout,
		retrieveAgent: kitot.TraceClient(tracer, "retrieve_agent_by_id")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrieveAgent",
			encodeRetrieveAgentRequest,
			decodeAgentResponse,
			pb.AgentRes{},
		).Endpoint()),
		retrieveAgentGroup: kitot.TraceClient(tracer, "retrieve_agent_group_by_id")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrieveAgentGroup",
			encodeRetrieveAgentGroupRequest,
			decodeAgentGroupResponse,
			pb.AgentGroupRes{},
		).Endpoint()),
		retrieveOwnerByChannelID: kitot.TraceClient(tracer, "retrieve_owner_id_by_channel_id")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrieveOwnerByChannelID",
			encodeRetrieveOwnerByChannelIDRequest,
			decodeOwnerResponse,
			pb.OwnerRes{},
		).Endpoint()),
	}
}

func encodeRetrieveAgentRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByIDReq)
	return &pb.AgentByIDReq{
		AgentID: req.AgentID,
		OwnerID: req.OwnerID,
	}, nil
}

func decodeAgentResponse(ctx context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.AgentRes)
	return agentRes{
		id:      res.GetId(),
		name:    res.GetName(),
		channel: res.GetChannel(),
	}, nil
}

func encodeRetrieveAgentGroupRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessAgByIDReq)
	return &pb.AgentGroupByIDReq{
		AgentGroupID: req.AgentGroupID,
		OwnerID:      req.OwnerID,
	}, nil
}

func decodeAgentGroupResponse(ctx context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.AgentGroupRes)
	return agentGroupRes{
		id:      res.GetId(),
		name:    res.GetName(),
		channel: res.GetChannel(),
	}, nil
}

func encodeRetrieveOwnerByChannelIDRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessOwnerByChannelIDReq)
	return &pb.OwnerByChannelIDReq{
		Channel: req.ChannelID,
	}, nil
}

func decodeOwnerResponse(ctx context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.OwnerRes)
	return ownerRes{
		ownerID: res.GetOwnerID(),
	}, nil
}
