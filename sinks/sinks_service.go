package sinks

import (
	"context"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/sinks/backend"
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

func (svc sinkService) CreateSink(ctx context.Context, token string, sink Sink) (Sink, error) {

	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}

	sink.MFOwnerID = mfOwnerID

	id, err := svc.sinkRepo.Save(ctx, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}
	sink.ID = id
	return sink, nil
}

func (svc sinkService) ListBackends(ctx context.Context, token string)([]string, error) {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return []string{}, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}
	return backend.GetList(), nil
}


func (svc sinkService) GetBackend(ctx context.Context, token string, key string)(backend.Backend, error) {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return backend.GetBackend(""), errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}
	return backend.GetBackend(key), nil
}