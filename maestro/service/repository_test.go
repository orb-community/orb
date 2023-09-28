package service

import (
	"context"
	"errors"
	"github.com/orb-community/orb/maestro/deployment"
	"go.uber.org/zap"
)

type fakeRepository struct {
	logger       *zap.Logger
	inMemoryDict map[string]*deployment.Deployment
}

func NewFakeRepository(logger *zap.Logger) deployment.Repository {
	return &fakeRepository{logger: logger, inMemoryDict: make(map[string]*deployment.Deployment)}
}

func (f *fakeRepository) FetchAll(ctx context.Context) ([]deployment.Deployment, error) {
	var allDeployments []deployment.Deployment
	for _, deploy := range f.inMemoryDict {
		allDeployments = append(allDeployments, *deploy)
	}
	return allDeployments, nil
}

func (f *fakeRepository) Add(ctx context.Context, deployment *deployment.Deployment) (*deployment.Deployment, error) {
	deployment.Id = "fake-id"
	f.inMemoryDict[deployment.SinkID] = deployment
	return deployment, nil
}

func (f *fakeRepository) Update(ctx context.Context, deployment *deployment.Deployment) (*deployment.Deployment, error) {
	f.inMemoryDict[deployment.SinkID] = deployment
	return deployment, nil
}

func (f *fakeRepository) UpdateStatus(ctx context.Context, ownerID string, sinkId string, status string, errorMessage string) error {
	return nil
}

func (f *fakeRepository) Remove(ctx context.Context, ownerId string, sinkId string) error {
	delete(f.inMemoryDict, sinkId)
	return nil
}

func (f *fakeRepository) FindByOwnerAndSink(ctx context.Context, _ string, sinkId string) (*deployment.Deployment, error) {
	deploy, ok := f.inMemoryDict[sinkId]
	if ok {
		return deploy, nil
	}
	return nil, errors.New("not found")
}

func (f *fakeRepository) FindByCollectorName(ctx context.Context, collectorName string) (*deployment.Deployment, error) {
	return nil, nil
}
