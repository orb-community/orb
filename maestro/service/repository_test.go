package service

import (
	"context"
	"github.com/orb-community/orb/maestro/deployment"
	"go.uber.org/zap"
)

type fakeRepository struct {
	logger       *zap.Logger
	inMemoryDict map[string]*deployment.Deployment
}

func NewFakeRepository(logger *zap.Logger) deployment.Repository {

	return &fakeRepository{logger: logger}
}

func (f *fakeRepository) FetchAll(ctx context.Context) ([]deployment.Deployment, error) {
	return nil, nil
}

func (f *fakeRepository) Add(ctx context.Context, deployment *deployment.Deployment) (*deployment.Deployment, error) {
	return nil, nil
}

func (f *fakeRepository) Update(ctx context.Context, deployment *deployment.Deployment) (*deployment.Deployment, error) {
	return nil, nil
}

func (f *fakeRepository) UpdateStatus(ctx context.Context, ownerID string, sinkId string, status string, errorMessage string) error {
	return nil
}

func (f *fakeRepository) Remove(ctx context.Context, ownerId string, sinkId string) error {
	return nil
}

func (f *fakeRepository) FindByOwnerAndSink(ctx context.Context, ownerId string, sinkId string) (*deployment.Deployment, error) {
	return nil, nil
}

func (f *fakeRepository) FindByCollectorName(ctx context.Context, collectorName string) (*deployment.Deployment, error) {
	return nil, nil
}
