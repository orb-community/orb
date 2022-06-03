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
	retrieveAgent                kitgrpc.Handler
	retrieveAgentGroup           kitgrpc.Handler
	retrieveOwnerByChannelID     kitgrpc.Handler
	retrieveAgentInfoByChannelID kitgrpc.Handler
}

func NewServer(tracer opentracing.Tracer, svc fleet.Service) pb.FleetServiceServer {
	return &grpcServer{
		retrieveAgent: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_agent")(retrieveAgentEndpoint(svc)),
			decodeRetrieveAgentRequest,
			encodeAgentResponse,
		),
		retrieveAgentGroup: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_agent_group")(retrieveAgentGroupEndpoint(svc)),
			decodeRetrieveAgentGroupRequest,
			encodeAgentGroupResponse,
		),
		retrieveOwnerByChannelID: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_owner_by_channel_id")(retrieveOwnerByChannelIDEndpoint(svc)),
			decodeRetrieveOwnerByChannelIDRequest,
			encodeOwnerResponse,
		),
		retrieveAgentInfoByChannelID: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_agent_info_by_channel_id")(retrieveAgentInfoByChannelIDEndpoint(svc)),
			decodeRetrieveAgentInfoByChannelIDRequest,
			encodeAgentInfoResponse,
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

func (gs *grpcServer) RetrieveAgentGroup(ctx context.Context, req *pb.AgentGroupByIDReq) (*pb.AgentGroupRes, error) {
	_, res, err := gs.retrieveAgentGroup.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.AgentGroupRes), nil
}

func (gs *grpcServer) RetrieveOwnerByChannelID(ctx context.Context, req *pb.OwnerByChannelIDReq) (*pb.OwnerRes, error) {
	_, res, err := gs.retrieveOwnerByChannelID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*pb.OwnerRes), nil
}

func (gs *grpcServer) RetrieveAgentInfoByChannelID(ctx context.Context, req *pb.AgentInfoByChannelIDReq) (*pb.AgentInfoRes, error) {
	_, res, err := gs.retrieveAgentInfoByChannelID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.AgentInfoRes), nil
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

func decodeRetrieveAgentGroupRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AgentGroupByIDReq)
	return accessAgByIDReq{AgentGroupID: req.AgentGroupID, OwnerID: req.OwnerID}, nil
}

func encodeAgentGroupResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(agentGroupRes)
	return &pb.AgentGroupRes{
		Id:      res.id,
		Name:    res.name,
		Channel: res.channel,
	}, nil
}

func decodeRetrieveOwnerByChannelIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.OwnerByChannelIDReq)
	return accessOwnerByChannelIDReq{ChannelID: req.Channel}, nil
}

func encodeOwnerResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(ownerRes)
	return &pb.OwnerRes{
		OwnerID:   res.ownerID,
		AgentName: res.agentName,
	}, nil
}

func decodeRetrieveAgentInfoByChannelIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AgentInfoByChannelIDReq)
	return accessAgentInfoByChannelIDReq{ChannelID: req.Channel}, nil
}

func encodeAgentInfoResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(agentInfoRes)
	return &pb.AgentInfoRes{
		OwnerID:   res.ownerID,
		AgentName: res.agentName,
		AgentTags: res.agentTags,
		OrbTags:   res.orbTags,
	}, nil
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case fleet.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid can access request")
	case fleet.ErrNotFound:
		return status.Error(codes.NotFound, "not found")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
