// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"context"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/pb"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.SinkServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	pb.UnimplementedSinkServiceServer
	retrieveSink    kitgrpc.Handler
	passwordService sinks.PasswordService
}

func NewServer(tracer opentracing.Tracer, svc sinks.SinkService) pb.SinkServiceServer {
	return &grpcServer{
		retrieveSink: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_sink")(retrieveSinkEndpoint(svc)),
			decodeRetrieveSinkRequest,
			encodeSinkResponse,
		),
	}
}

func (gs *grpcServer) RetrieveSink(ctx context.Context, req *pb.SinkByIDReq) (*pb.SinkRes, error) {
	_, res, err := gs.retrieveSink.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.SinkRes), nil
}

func decodeRetrieveSinkRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SinkByIDReq)
	return accessByIDReq{SinkID: req.SinkID, OwnerID: req.OwnerID}, nil
}

func encodeSinkResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(sinkRes)
	return &pb.SinkRes{
		Id:          res.id,
		Name:        res.name,
		Description: res.description,
		Tags:        res.tags,
		State:       res.state,
		Error:       res.error,
		Backend:     res.backend,
		Config:      res.config,
	}, nil
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case sinks.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid can access request")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
