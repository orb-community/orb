package policies

import (
	"context"
	"github.com/orb-community/orb/policies/backend/otel"
	"time"

	"github.com/mainflux/mainflux"
	fleetpb "github.com/orb-community/orb/fleet/pb"
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/policies/backend/diode"
	"github.com/orb-community/orb/policies/backend/orb"
	"github.com/orb-community/orb/policies/backend/pktvisor"
	sinkpb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
)

type PageMetadata struct {
	Total    uint64
	Offset   uint64         `json:"offset,omitempty"`
	Limit    uint64         `json:"limit,omitempty"`
	Name     string         `json:"name,omitempty"`
	Order    string         `json:"order,omitempty"`
	Dir      string         `json:"dir,omitempty"`
	Metadata types.Metadata `json:"metadata,omitempty"`
	Tags     types.Tags     `json:"tags,omitempty"`
}

var _ Service = (*policiesService)(nil)

type policiesService struct {
	logger          *zap.Logger
	auth            mainflux.AuthServiceClient
	fleetGrpcClient fleetpb.FleetServiceClient
	sinksGrpcClient sinkpb.SinkServiceClient
	repo            Repository
}

func (s policiesService) identify(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := s.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return res.GetId(), nil
}

func New(logger *zap.Logger, auth mainflux.AuthServiceClient, repo Repository, fleetGrpcClient fleetpb.FleetServiceClient, sinksGrpcclient sinkpb.SinkServiceClient) Service {

	orb.Register()
	otel.Register()
	pktvisor.Register()
	diode.Register()

	return &policiesService{
		logger:          logger,
		auth:            auth,
		repo:            repo,
		fleetGrpcClient: fleetGrpcClient,
		sinksGrpcClient: sinksGrpcclient,
	}
}
