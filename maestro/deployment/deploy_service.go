package deployment

import (
	"context"
	"github.com/orb-community/orb/maestro/config"
	"github.com/orb-community/orb/maestro/kubecontrol"
	"go.uber.org/zap"
	"time"
)

type DeployService interface {
	HandleSinkCreate(ctx context.Context, sink config.SinkData) error
	HandleSinkUpdate(ctx context.Context, sink config.SinkData) error
	HandleSinkDelete(ctx context.Context, sink config.SinkData) error
	HandleSinkActivity(ctx context.Context, sink config.SinkData) error
}

type deployService struct {
	logger            *zap.Logger
	deploymentService Service
	kubecontrol       kubecontrol.Service

	// Configuration for KafkaURL from Orb Deployment
	kafkaUrl string
}

var _ DeployService = (*deployService)(nil)

func NewDeployService(logger *zap.Logger, service Service, kubecontrol kubecontrol.Service) DeployService {
	namedLogger := logger.Named("deploy-service")
	return &deployService{logger: namedLogger, deploymentService: service}
}

// HandleSinkCreate will create deployment entry in postgres, will create deployment in Redis, to prepare for SinkActivity
func (d *deployService) HandleSinkCreate(ctx context.Context, sink config.SinkData) error {
	now := time.Now()
	// Create Deployment Entry
	entry := Deployment{
		OwnerID:                 sink.OwnerID,
		SinkID:                  sink.SinkID,
		Config:                  sink.Config,
		LastStatus:              "provisioning",
		LastStatusUpdate:        &now,
		LastErrorMessage:        "",
		LastErrorTime:           nil,
		CollectorName:           "",
		LastCollectorDeployTime: nil,
		LastCollectorStopTime:   nil,
	}
	// Use deploymentService, which will create deployment in both postgres and redis
	err := d.deploymentService.CreateDeployment(ctx, &entry)
	if err != nil {
		d.logger.Error("error trying to create deployment entry", zap.Error(err))
		return err
	}
	return nil
}

func (d *deployService) HandleSinkUpdate(ctx context.Context, sink config.SinkData) error {
	now := time.Now()
	// check if exists deployment entry from postgres
	entry, manifest, err := d.deploymentService.GetDeployment(ctx, sink.OwnerID, sink.SinkID)
	if err != nil {
		d.logger.Error("error trying to get deployment entry", zap.Error(err))
		return err
	}
	// update sink status to provisioning
	err = d.deploymentService.UpdateStatus(ctx, sink.OwnerID, sink.SinkID, "provisioning", "")
	if err != nil {
		return err
	}
	err = d.kubecontrol.DeleteOtelCollector(ctx, sink.OwnerID, sink.SinkID, manifest)
	if err != nil {
		return err
	}
	entry.

	return nil
}

func (d *deployService) HandleSinkDelete(ctx context.Context, sink config.SinkData) error {
	//TODO implement me
	panic("implement me")
}

func (d *deployService) HandleSinkActivity(ctx context.Context, sink config.SinkData) error {
	//TODO implement me
	panic("implement me")
}
