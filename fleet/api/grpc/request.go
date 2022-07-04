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

type accessAgByIDReq struct {
	AgentGroupID string
	OwnerID      string
}

func (req accessAgByIDReq) validate() error {
	if req.AgentGroupID == "" || req.OwnerID == "" {
		return fleet.ErrMalformedEntity
	}

	return nil
}

type accessOwnerByChannelIDReq struct {
	ChannelID string
}

func (req accessOwnerByChannelIDReq) validate() error {
	if req.ChannelID == "" {
		return fleet.ErrMalformedEntity
	}
	return nil
}

type accessAgentInfoByChannelIDReq struct {
	ChannelID string
}

func (req accessAgentInfoByChannelIDReq) validate() error {
	if req.ChannelID == "" {
		return fleet.ErrMalformedEntity
	}
	return nil
}
