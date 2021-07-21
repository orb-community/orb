package sinks

import (
	"context"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
)

var (
	ErrCreateSink = errors.New("failed to create Sink")

	ErrConflictSink = errors.New("entity already exists")
	)

func (svc sinkService) ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Page{}, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return svc.sinkRepo.RetrieveAll(ctx, res.GetId(), pm)
}

func (s sinkService) CreateSink(ctx context.Context, token string, sink Sink) (Sink, error) {

	mfOwnerID, err := s.identify(token)
	if err != nil {
		return Sink{}, err
	}

	sink.MFOwnerID = mfOwnerID

	id, err := s.sinkRepo.Save(ctx, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}
	sink.ID = id
	return sink, nil
}