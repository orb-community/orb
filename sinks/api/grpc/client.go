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
	"github.com/ns1labs/orb/sinks/pb"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

var _ pb.SinkServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	logger        *zap.Logger
	timeout       time.Duration
	retrieveSink  endpoint.Endpoint
	retrieveSinks endpoint.Endpoint
}

func (client grpcClient) RetrieveSinks(ctx context.Context, in *pb.SinksFilterReq, _ ...grpc.CallOption) (*pb.SinksRes, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	sinksFilter := sinksFilter{
		isOtel: in.OtelEnabled,
	}

	res, err := client.retrieveSinks(ctx, sinksFilter)
	if err != nil {
		client.logger.Error("error during retrieve sinks", zap.Error(err))
		return nil, err
	}
	ir := res.(sinksRes)
	sinkList := make([]*pb.SinkRes, len(ir.sinks))
	for i, sinkResponse := range ir.sinks {
		sinkList[i] = &pb.SinkRes{
			Id:          sinkResponse.id,
			Name:        sinkResponse.name,
			Description: sinkResponse.description,
			Tags:        sinkResponse.tags,
			State:       sinkResponse.state,
			Error:       sinkResponse.error,
			Backend:     sinkResponse.backend,
			Config:      sinkResponse.config,
		}
	}
	return &pb.SinksRes{Sinks: sinkList}, nil
}

func (client grpcClient) RetrieveSink(ctx context.Context, in *pb.SinkByIDReq, _ ...grpc.CallOption) (*pb.SinkRes, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	ar := accessByIDReq{
		SinkID:  in.SinkID,
		OwnerID: in.OwnerID,
	}

	res, err := client.retrieveSink(ctx, ar)
	if err != nil {
		return nil, err
	}

	ir := res.(sinkRes)
	return &pb.SinkRes{
		Id:          ir.id,
		Name:        ir.name,
		Description: ir.description,
		Tags:        ir.tags,
		State:       ir.state,
		Error:       ir.error,
		Backend:     ir.backend,
		Config:      ir.config,
	}, nil
}

func NewClient(tracer opentracing.Tracer, conn *grpc.ClientConn, timeout time.Duration, logger *zap.Logger) pb.SinkServiceClient {
	svcName := "sinks.SinkService"

	return &grpcClient{
		logger:  logger,
		timeout: timeout,
		retrieveSink: kitot.TraceClient(tracer, "retrieve_sink")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrieveSink",
			encodeRetrieveSinkRequest,
			decodeSinkResponse,
			pb.SinkRes{},
		).Endpoint()),
		retrieveSinks: kitot.TraceClient(tracer, "retrieve_sinks_internal")(kitgrpc.NewClient(
			conn,
			svcName,
			"RetrieveSinksInternal",
			encodeRetrieveSinksRequest,
			decodeSinksResponse,
			pb.SinksRes{},
		).Endpoint()),
	}
}

func encodeRetrieveSinksRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(sinksFilter)
	return &pb.SinksFilterReq{OtelEnabled: req.isOtel}, nil
}

func decodeSinksResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.SinksRes)
	sinkList := make([]sinkRes, len(res.Sinks))
	for i, sink := range res.Sinks {
		sinkList[i] = sinkRes{
			id:          sink.Id,
			name:        sink.Name,
			description: sink.Description,
			tags:        sink.Tags,
			state:       sink.State,
			error:       sink.Error,
			backend:     sink.Backend,
			config:      sink.Config,
		}
	}
	return sinksRes{sinks: sinkList}, nil
}

func encodeRetrieveSinkRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(accessByIDReq)
	return &pb.SinkByIDReq{SinkID: req.SinkID, OwnerID: req.OwnerID}, nil
}

func decodeSinkResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*pb.SinkRes)
	return sinkRes{
		id:          res.GetId(),
		name:        res.GetName(),
		description: res.GetDescription(),
		tags:        res.GetTags(),
		state:       res.GetState(),
		error:       res.GetError(),
		backend:     res.GetBackend(),
		config:      res.GetConfig(),
	}, nil
}
