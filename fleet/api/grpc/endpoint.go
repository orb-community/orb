package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/ns1labs/orb/fleet"
)

func retrieveAgentEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}
