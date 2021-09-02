package grpc

import (
	"github.com/ns1labs/orb/fleet"
)

type accessByIDReq struct {
	AgentID string
	OwnerID string
}

func (req accessByIDReq) validate() error {
	if req.AgentID == "" || req.OwnerID == "" {
		return fleet.ErrMalformedEntity
	}

	return nil
}
