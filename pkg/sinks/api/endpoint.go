// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package api

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/ns1labs/orb/pkg/sinks"
)

func addEndpoint(svc sinks.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		saved, err := svc.Add()
		if err != nil {
			return nil, err
		}

		res := sinkRes{
			id:      saved.SinkID,
			created: true,
		}

		return res, nil
	}
}
