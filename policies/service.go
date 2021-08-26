package policies

import (
	"context"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies/backend/orb"
	"github.com/ns1labs/orb/policies/backend/pktvisor"
	"time"
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
	auth mainflux.AuthServiceClient
	repo Repository
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

func New(auth mainflux.AuthServiceClient, repo Repository) Service {

	orb.Register()
	pktvisor.Register()

	return &policiesService{
		auth: auth,
		repo: repo,
	}
}
