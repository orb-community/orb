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
		copy := copyDeploy(deploy)
		allDeployments = append(allDeployments, copy)
	}
	return allDeployments, nil
}

func (f *fakeRepository) Add(_ context.Context, deployment *deployment.Deployment) (*deployment.Deployment, error) {
	deployment.Id = "fake-id"
	copy := copyDeploy(deployment)
	f.inMemoryDict[deployment.SinkID] = &copy
	return deployment, nil
}

func (f *fakeRepository) Update(_ context.Context, deployment *deployment.Deployment) (*deployment.Deployment, error) {
	copy := copyDeploy(deployment)
	f.inMemoryDict[deployment.SinkID] = &copy
	return deployment, nil
}

func (f *fakeRepository) UpdateStatus(_ context.Context, _ string, _ string, _ string, _ string) error {
	return nil
}

func (f *fakeRepository) Remove(_ context.Context, _ string, sinkId string) error {
	delete(f.inMemoryDict, sinkId)
	return nil
}

func (f *fakeRepository) FindByOwnerAndSink(ctx context.Context, _ string, sinkId string) (*deployment.Deployment, error) {
	deploy, ok := f.inMemoryDict[sinkId]
	if ok {
		copy := copyDeploy(deploy)
		return &copy, nil
	}
	return nil, errors.New("not found")
}

func (f *fakeRepository) FindByCollectorName(_ context.Context, _ string) (*deployment.Deployment, error) {
	return nil, nil
}

func copyDeploy(src *deployment.Deployment) deployment.Deployment {
	deploy := deployment.Deployment{
		Id:                      src.Id,
		OwnerID:                 src.OwnerID,
		SinkID:                  src.SinkID,
		Backend:                 src.Backend,
		Config:                  src.Config,
		LastStatus:              src.LastStatus,
		LastStatusUpdate:        src.LastStatusUpdate,
		LastErrorMessage:        src.LastErrorMessage,
		LastErrorTime:           src.LastErrorTime,
		CollectorName:           src.CollectorName,
		LastCollectorDeployTime: src.LastCollectorDeployTime,
		LastCollectorStopTime:   src.LastCollectorStopTime,
	}
	return deploy
}
