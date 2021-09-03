package grpc

import (
	"context"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.FleetServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	pb.UnimplementedFleetServiceServer
	retrieveAgent kitgrpc.Handler
}

func NewServer(tracer opentracing.Tracer, svc fleet.Service) pb.FleetServiceServer {
	return &grpcServer{
		retrieveAgent: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_agent")(retrieveAgentEndpoint(svc)),
			decodeRetrieveAgentRequest,
			encodeAgentResponse,
		),
	}
}

func (gs *grpcServer) RetrieveAgent(ctx context.Context, req *pb.AgentByIDReq) (*pb.AgentRes, error) {
	_, res, err := gs.retrieveAgent.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.AgentRes), nil
}

func decodeRetrieveAgentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AgentByIDReq)
	return accessByIDReq{AgentID: req.AgentID, OwnerID: req.OwnerID}, nil
}

func encodeAgentResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(agentRes)
	return &pb.AgentRes{
		Id:      res.id,
		Name:    res.name,
		Channel: res.channel,
	}, nil
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case fleet.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid can access request")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
