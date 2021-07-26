package sinks

import (
	"context"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/sinks/backend"
)

var (
	ErrCreateSink                 = errors.New("failed to create Sink")
	ErrConflictSink               = errors.New("entity already exists")
	ErrUnsupportedContentTypeSink = errors.New("unsupported content type")
)

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

func (svc sinkService) UpdateSink(ctx context.Context, token string, sink Sink) error {
	skOwnerID, err := svc.identify(token)
	if err != nil {
		return err
	}
	sink.MFOwnerID = skOwnerID
	return svc.sinkRepo.Update(ctx, sink)
}

func (svc sinkService) ListBackends(ctx context.Context, token string) ([]string, error) {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return []string{}, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}
	return backend.GetList(), nil
}

func (svc sinkService) ViewBackend(ctx context.Context, token string, key string) (backend.Backend, error) {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}
	res := backend.GetBackend(key)
	if res == nil {
		return nil, errors.Wrap(errors.ErrNotFound, err)
	}
	return res, nil
}

func (svc sinkService) ViewSink(ctx context.Context, token string, key string) (Sink, error) {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Sink{}, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}
	res, err := svc.sinkRepo.RetrieveById(ctx, key)
	if err != nil {
		return Sink{}, errors.Wrap(errors.ErrNotFound, err)
	}
	return res, nil
}

func (svc sinkService) ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Page{}, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return svc.sinkRepo.RetrieveAll(ctx, res.GetId(), pm)
}
