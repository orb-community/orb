package deployment

import (
	"context"
	"github.com/orb-community/orb/pkg/types"
	"time"
)

type Repository interface {
	FetchAll(ctx context.Context) ([]Deployment, error)
	Add(ctx context.Context, deployment Deployment) (Deployment, error)
	Update(ctx context.Context, deployment Deployment) (Deployment, error)
	Remove(ctx context.Context, ownerId string, sinkId string) error
	FindByOwnerAndSink(ctx context.Context, ownerId string, sinkId string) (Deployment, error)
}

type Deployment struct {
	Id                      string
	OwnerID                 string
	SinkID                  string
	Config                  types.Metadata
	LastStatus              string
	LastStatusUpdate        time.Time
	LastErrorMessage        string
	LastErrorTime           time.Time
	CollectorName           string
	LastCollectorDeployTime time.Time
	LastCollectorStopTime   time.Time
}

var _ Repository = (*repositoryService)(nil)

type repositoryService struct {
}

func (r *repositoryService) FetchAll(ctx context.Context) ([]Deployment, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repositoryService) Add(ctx context.Context, deployment Deployment) (Deployment, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repositoryService) Update(ctx context.Context, deployment Deployment) (Deployment, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repositoryService) Remove(ctx context.Context, ownerId string, sinkId string) error {
	//TODO implement me
	panic("implement me")
}

func (r *repositoryService) FindByOwnerAndSink(ctx context.Context, ownerId string, sinkId string) (Deployment, error) {
	//TODO implement me
	panic("implement me")
}
