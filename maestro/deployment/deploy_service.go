package deployment

import (
	"context"
	"github.com/orb-community/orb/maestro/kubecontrol"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	"go.uber.org/zap"
	"time"
)

type DeployService interface {
	HandleSinkCreate(ctx context.Context, event maestroredis.SinksUpdateEvent) error
	HandleSinkUpdate(ctx context.Context, event maestroredis.SinksUpdateEvent) error
	HandleSinkDelete(ctx context.Context, event maestroredis.SinksUpdateEvent) error
	HandleSinkActivity(ctx context.Context, event maestroredis.SinksUpdateEvent) error
}

type deployService struct {
	logger            *zap.Logger
	deploymentService Service
	// Configuration for KafkaURL from Orb Deployment
	kafkaUrl string
}

var _ DeployService = (*deployService)(nil)

func NewDeployService(logger *zap.Logger, service Service, kubecontrol kubecontrol.Service) DeployService {
	namedLogger := logger.Named("deploy-service")
	return &deployService{logger: namedLogger, deploymentService: service}
}

// HandleSinkCreate will create deployment entry in postgres, will create deployment in Redis, to prepare for SinkActivity
func (d *deployService) HandleSinkCreate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	now := time.Now()
	// Create Deployment Entry
	entry := Deployment{
		OwnerID:                 event.Owner,
		SinkID:                  event.SinkID,
		Config:                  event.Config,
		Backend:                 event.Backend,
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

func (d *deployService) HandleSinkUpdate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	now := time.Now()
	// check if exists deployment entry from postgres
	entry, _, err := d.deploymentService.GetDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		d.logger.Error("error trying to get deployment entry", zap.Error(err))
		return err
	}
	// async update sink status to provisioning
	go func() {
		_ = d.deploymentService.UpdateStatus(ctx, event.Owner, event.SinkID, "provisioning", "")
	}()
	// update deployment entry in postgres
	entry.Config = event.Config
	entry.LastCollectorStopTime = &now
	entry.LastStatus = "provisioning"
	entry.LastStatusUpdate = &now
	err = d.deploymentService.UpdateDeployment(ctx, entry)

	return nil
}

func (d *deployService) HandleSinkDelete(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	deploymentEntry, _, err := d.deploymentService.GetDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		d.logger.Warn("did not find collector entry for sink", zap.String("sink-id", event.SinkID))
		return err
	}
	if deploymentEntry.LastCollectorDeployTime != nil || deploymentEntry.LastCollectorDeployTime.Before(time.Now()) {
		if deploymentEntry.LastCollectorStopTime != nil || deploymentEntry.LastCollectorStopTime.Before(time.Now()) {
			d.logger.Warn("collector is not running, skipping")
		} else {
			//
		}
	}
	err = d.deploymentService.RemoveDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		return err
	}
	return nil
}

func (d *deployService) HandleSinkActivity(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	//TODO implement me
	panic("implement me")
}
