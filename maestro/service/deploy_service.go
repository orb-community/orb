package service

import (
	"context"
	"github.com/orb-community/orb/maestro/deployment"
	"github.com/orb-community/orb/maestro/kubecontrol"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// EventService will hold the business logic of the handling events from both Listeners
type EventService interface {
	HandleSinkCreate(ctx context.Context, event maestroredis.SinksUpdateEvent) error
	HandleSinkUpdate(ctx context.Context, event maestroredis.SinksUpdateEvent) error
	HandleSinkDelete(ctx context.Context, event maestroredis.SinksUpdateEvent) error
	HandleSinkActivity(ctx context.Context, event maestroredis.SinkerUpdateEvent) error
	HandleSinkIdle(ctx context.Context, event maestroredis.SinkerUpdateEvent) error
}

type eventService struct {
	logger            *zap.Logger
	deploymentService deployment.Service
	// Configuration for KafkaURL from Orb Deployment
	kafkaUrl string
}

var _ EventService = (*eventService)(nil)

func NewEventService(logger *zap.Logger, service deployment.Service, _ kubecontrol.Service) EventService {
	namedLogger := logger.Named("deploy-service")
	return &eventService{logger: namedLogger, deploymentService: service}
}

// HandleSinkCreate will create deployment entry in postgres, will create deployment in Redis, to prepare for SinkActivity
func (d *eventService) HandleSinkCreate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	// Create Deployment Entry
	entry := deployment.NewDeployment(event.Owner, event.SinkID, event.Config)
	// Use deploymentService, which will create deployment in both postgres and redis
	err := d.deploymentService.CreateDeployment(ctx, &entry)
	if err != nil {
		d.logger.Error("error trying to create deployment entry", zap.Error(err))
		return err
	}
	return nil
}

func (d *eventService) HandleSinkUpdate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
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

func (d *eventService) HandleSinkDelete(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	deploymentEntry, _, err := d.deploymentService.GetDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		d.logger.Warn("did not find collector entry for sink", zap.String("sink-id", event.SinkID))
		return err
	}
	if deploymentEntry.LastCollectorDeployTime != nil || deploymentEntry.LastCollectorDeployTime.Before(time.Now()) {
		if deploymentEntry.LastCollectorStopTime != nil || deploymentEntry.LastCollectorStopTime.Before(time.Now()) {
			d.logger.Warn("collector is not running, skipping")
		}
	}
	err = d.deploymentService.RemoveDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		return err
	}
	return nil
}

func (d *eventService) HandleSinkActivity(ctx context.Context, event maestroredis.SinkerUpdateEvent) error {
	if event.State != "active" {
		return errors.New("trying to deploy sink that is not active")
	}
	// check if exists deployment entry from postgres
	_, _, err := d.deploymentService.GetDeployment(ctx, event.OwnerID, event.SinkID)
	if err != nil {
		d.logger.Error("error trying to get deployment entry", zap.Error(err))
		return err
	}
	// async update sink status to provisioning
	go func() {
		_ = d.deploymentService.UpdateStatus(ctx, event.OwnerID, event.SinkID, "provisioning", "")
	}()
	_, err = d.deploymentService.NotifyCollector(ctx, event.OwnerID, event.SinkID, "deploy", "", "")
	if err != nil {
		d.logger.Error("error trying to notify collector", zap.Error(err))
		return err
	}
	err2 := d.deploymentService.UpdateStatus(ctx, event.OwnerID, event.SinkID, "provisioning_error", err.Error())
	if err2 != nil {
		d.logger.Warn("error during notifying provisioning error, customer will not be notified of error")
		d.logger.Error("error during update status", zap.Error(err))
		return err
	}

	return nil
}

func (d *eventService) HandleSinkIdle(ctx context.Context, event maestroredis.SinkerUpdateEvent) error {
	// check if exists deployment entry from postgres
	_, _, err := d.deploymentService.GetDeployment(ctx, event.OwnerID, event.SinkID)
	if err != nil {
		d.logger.Error("error trying to get deployment entry", zap.Error(err))
		return err
	}
	// async update sink status to idle
	go func() {
		_ = d.deploymentService.UpdateStatus(ctx, event.OwnerID, event.SinkID, "idle", "")
	}()
	_, err = d.deploymentService.NotifyCollector(ctx, event.OwnerID, event.SinkID, "deploy", "", "")
	if err != nil {
		d.logger.Error("error trying to notify collector", zap.Error(err))
		return err
	}
	err2 := d.deploymentService.UpdateStatus(ctx, event.OwnerID, event.SinkID, "provisioning_error", err.Error())
	if err2 != nil {
		d.logger.Warn("error during notifying provisioning error, customer will not be notified of error")
		d.logger.Error("error during update status", zap.Error(err))
		return err
	}
	return nil
}
