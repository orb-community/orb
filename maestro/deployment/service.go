package deployment

import (
	"context"
	"errors"
	"github.com/orb-community/orb/maestro/redis/consumer"
	"github.com/orb-community/orb/maestro/redis/producer"
	"go.uber.org/zap"
	"time"
)

type Service interface {
	// CreateDeployment to be used to create the deployment when there is a sink.create
	CreateDeployment(ctx context.Context, deployment *Deployment) error
	// GetDeployment to be used to get the deployment information for creating the collector or monitoring the collector
	GetDeployment(ctx context.Context, ownerID string, sinkId string) (*Deployment, string, error)
	// UpdateDeployment to be used to update the deployment when there is a sink.update
	UpdateDeployment(ctx context.Context, deployment *Deployment) error
	// UpdateStatus to be used to update the status of the sink, when there is an error or when the sink is running
	UpdateStatus(ctx context.Context, ownerID string, sinkId string, status string, errorMessage string) error
	// RemoveDeployment to be used to remove the deployment when there is a sink.delete
	RemoveDeployment(ctx context.Context, ownerID string, sinkId string) error
	// GetDeploymentByCollectorName to be used to get the deployment information for creating the collector or monitoring the collector
	GetDeploymentByCollectorName(ctx context.Context, collectorName string) (*Deployment, error)
	// NotifyCollector add collector information to deployment
	NotifyCollector(ctx context.Context, ownerID string, sinkId string, collectorName string, operation string, status string, errorMessage string) error
}

type deploymentService struct {
	dbRepository    Repository
	logger          *zap.Logger
	cacheRepository consumer.DeploymentHashsetRepository
	maestroProducer producer.Producer
}

var _ Service = (*deploymentService)(nil)

func NewDeploymentService(logger *zap.Logger, repository Repository) Service {
	namedLogger := logger.Named("deployment-service")
	return &deploymentService{logger: namedLogger, dbRepository: repository}
}

func (d *deploymentService) CreateDeployment(ctx context.Context, deployment *Deployment) error {
	if deployment == nil {
		return errors.New("deployment is nil")
	}
	added, err := d.dbRepository.Add(ctx, deployment)
	if err != nil {
		return err
	}
	d.logger.Info("added deployment", zap.String("id", added.Id),
		zap.String("ownerID", added.OwnerID), zap.String("sinkID", added.SinkID))
	err = d.cacheRepository.CreateDeploymentEntry(ctx, deployment)
	if err != nil {
		return err
	}
	return nil
}

func (d *deploymentService) GetDeployment(ctx context.Context, ownerID string, sinkId string) (*Deployment, string, error) {
	deployment, err := d.dbRepository.FindByOwnerAndSink(ctx, ownerID, sinkId)
	if err != nil {
		return nil, "", err
	}
	manifest := d.cacheRepository.GetDeploymentEntryFromSinkId(ctx, sinkId)
	return deployment, nil
}

func (d *deploymentService) UpdateDeployment(ctx context.Context, deployment *Deployment) error {
	got, err := d.dbRepository.FindByOwnerAndSink(ctx, deployment.OwnerID, deployment.SinkID)
	if err != nil {
		return errors.New("could not find deployment to update")
	}
	err = deployment.Merge(*got)
	if err != nil {
		d.logger.Error("error during merge of deployments", zap.Error(err))
		return err
	}
	if deployment == nil {
		return errors.New("deployment is nil")
	}
	updated, err := d.dbRepository.Update(ctx, deployment)
	if err != nil {
		return err
	}
	d.logger.Info("updated deployment", zap.String("ownerID", updated.OwnerID),
		zap.String("sinkID", updated.SinkID))
	return nil
}

func (d *deploymentService) NotifyCollector(ctx context.Context, ownerID string, sinkId string, collectorName string, operation string, status string, errorMessage string) error {
	got, err := d.dbRepository.FindByOwnerAndSink(ctx, ownerID, sinkId)
	if err != nil {
		return errors.New("could not find deployment to update")
	}
	now := time.Now()
	got.CollectorName = collectorName
	if operation == "delete" {
		got.LastCollectorStopTime = &now
	} else if operation == "deploy" {
		got.LastCollectorDeployTime = &now
	}
	if status != "" {
		got.LastStatus = status
		got.LastStatusUpdate = &now
	}
	if errorMessage != "" {
		got.LastErrorMessage = errorMessage
		got.LastErrorTime = &now
	}
	updated, err := d.dbRepository.Update(ctx, got)
	if err != nil {
		return err
	}
	d.logger.Info("updated deployment information for collector and status or error",
		zap.String("ownerID", updated.OwnerID), zap.String("sinkID", updated.SinkID),
		zap.String("collectorName", updated.CollectorName),
		zap.String("status", updated.LastStatus), zap.String("errorMessage", updated.LastErrorMessage))
	return nil
}

// UpdateStatus this will change the status in postgres and notify sinks service to show new status to user
func (d *deploymentService) UpdateStatus(ctx context.Context, ownerID string, sinkId string, status string, errorMessage string) error {
	got, err := d.dbRepository.FindByOwnerAndSink(ctx, ownerID, sinkId)
	if err != nil {
		return errors.New("could not find deployment to update")
	}
	now := time.Now()
	if status != "" {
		got.LastStatus = status
		got.LastStatusUpdate = &now
	}
	if errorMessage != "" {
		got.LastErrorMessage = errorMessage
		got.LastErrorTime = &now
	}
	updated, err := d.dbRepository.Update(ctx, got)
	if err != nil {
		return err
	}
	d.logger.Info("updated deployment status",
		zap.String("ownerID", updated.OwnerID), zap.String("sinkID", updated.SinkID),
		zap.String("status", updated.LastStatus), zap.String("errorMessage", updated.LastErrorMessage))

	return nil
}

// RemoveDeployment this will remove the deployment from postgres and redis
func (d *deploymentService) RemoveDeployment(ctx context.Context, ownerID string, sinkId string) error {
	err := d.dbRepository.Remove(ctx, ownerID, sinkId)
	if err != nil {
		return err
	}
	d.logger.Info("removed deployment", zap.String("ownerID", ownerID), zap.String("sinkID", sinkId))
	return nil
}

func (d *deploymentService) GetDeploymentByCollectorName(ctx context.Context, collectorName string) (*Deployment, error) {
	deployment, err := d.dbRepository.FindByCollectorName(ctx, collectorName)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}
