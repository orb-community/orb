package service

import (
	"context"
	"github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type mockSinksPb struct {
	logger *zap.Logger
}

func NewSinksPb(logger *zap.Logger) pb.SinkServiceClient {
	return &mockSinksPb{logger: logger}
}

var _ pb.SinkServiceClient = (*mockSinksPb)(nil)

func (m mockSinksPb) RetrieveSink(ctx context.Context, in *pb.SinkByIDReq, opts ...grpc.CallOption) (*pb.SinkRes, error) {
	return nil, nil
}

func (m mockSinksPb) RetrieveSinks(ctx context.Context, in *pb.SinksFilterReq, opts ...grpc.CallOption) (*pb.SinksRes, error) {
	return nil, nil
}
