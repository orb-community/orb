package service

import (
	"context"
	"encoding/json"
	maestroerrors "github.com/orb-community/orb/maestro/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/pb"
	"time"

	"github.com/orb-community/orb/maestro/deployment"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/pkg/errors"
	"go.uber.org/zap"
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
	sinkGrpcClient    pb.SinkServiceClient
	// Configuration for KafkaURL from Orb Deployment
	kafkaUrl string
}

var _ EventService = (*eventService)(nil)

func NewEventService(logger *zap.Logger, service deployment.Service, sinksGrpcClient *pb.SinkServiceClient) EventService {
	namedLogger := logger.Named("deploy-service")
	return &eventService{logger: namedLogger, deploymentService: service, sinkGrpcClient: *sinksGrpcClient}
}

// HandleSinkCreate will create deployment entry in postgres, will create deployment in Redis, to prepare for SinkActivity
func (d *eventService) HandleSinkCreate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	d.logger.Debug("handling sink create event", zap.String("sink-id", event.SinkID), zap.String("owner-id", event.Owner))
	// Create Deployment Entry
	entry := deployment.NewDeployment(event.Owner, event.SinkID, event.Config, event.Backend)
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
	d.logger.Debug("handling sink update event", zap.String("sink-id", event.SinkID))
	// check if exists deployment entry from postgres
	entry, _, err := d.deploymentService.GetDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		if err.Error() != "not found" {
			d.logger.Error("error trying to get deployment entry", zap.Error(err))
			return err
		} else {
			newEntry := deployment.NewDeployment(event.Owner, event.SinkID, event.Config, event.Backend)
			err := d.deploymentService.CreateDeployment(ctx, &newEntry)
			if err != nil {
				d.logger.Error("error trying to recreate deployment entry", zap.Error(err))
				return err
			}
			entry = &newEntry
		}
	}
	// update deployment entry in postgres
	err = entry.SetConfig(event.Config)
	if err != nil {
		return err
	}
	entry.LastCollectorStopTime = &now
	entry.LastStatus = "unknown"
	entry.LastStatusUpdate = &now
	entry.LastErrorMessage = ""
	entry.LastErrorTime = nil
	err = d.deploymentService.UpdateDeployment(ctx, entry)

	return nil
}

func (d *eventService) HandleSinkDelete(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	d.logger.Debug("handling sink delete event", zap.String("sink-id", event.SinkID))
	deploymentEntry, _, err := d.deploymentService.GetDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		d.logger.Warn("did not find collector entry for sink", zap.String("sink-id", event.SinkID))
		return err
	}
	if deploymentEntry.LastCollectorDeployTime == nil || deploymentEntry.LastCollectorDeployTime.Before(time.Now()) {
		if deploymentEntry.LastCollectorStopTime == nil || deploymentEntry.LastCollectorStopTime.Before(time.Now()) {
			d.logger.Warn("collector is not running, skipping")
		}
	}
	err = d.deploymentService.RemoveDeployment(ctx, event.Owner, event.SinkID)
	if err != nil {
		d.logger.Warn("error removing deployment entry, deployment will be orphan", zap.Error(err))
		return err
	}
	return nil
}

func (d *eventService) HandleSinkActivity(ctx context.Context, event maestroredis.SinkerUpdateEvent) error {
	if event.State != "active" {
		d.logger.Error("trying to deploy sink that is not active", zap.String("sink-id", event.SinkID),
			zap.String("status", event.State))
		return errors.New("trying to deploy sink that is not active")
	}
	deploymentEntry, _, err := d.deploymentService.GetDeployment(ctx, event.OwnerID, event.SinkID)
	if err != nil {
		if err == maestroerrors.NotFound {
			d.logger.Info("did not find collector entry for sink, retrieving from sinks grpc", zap.String("sink-id", event.SinkID))
			sink, err := d.sinkGrpcClient.RetrieveSink(ctx, &pb.SinkByIDReq{
				SinkID:  event.SinkID,
				OwnerID: event.OwnerID,
			})
			if err != nil {
				d.logger.Error("error retrieving sink from grpc", zap.Error(err))
				return err
			}
			metadata := make(map[string]interface{})
			err = json.Unmarshal(sink.Config, &metadata)
			if err != nil {
				d.logger.Error("error unmarshalling sink metadata", zap.Error(err))
				return err
			}
			newEntry := deployment.NewDeployment(event.OwnerID, event.SinkID, types.FromMap(metadata), sink.Backend)
			err = d.deploymentService.CreateDeployment(ctx, &newEntry)
			if err != nil {
				d.logger.Error("error trying to recreate deployment entry", zap.Error(err))
				return err
			}
			deploymentEntry, _, err = d.deploymentService.GetDeployment(ctx, event.OwnerID, event.SinkID)
			if err != nil {
				d.logger.Error("error trying to recreate deployment entry", zap.Error(err))
				return err
			}
		} else {
			d.logger.Warn("did not find collector entry for sink", zap.String("sink-id", event.SinkID))
			return err
		}
	}
	d.logger.Debug("handling sink activity event", zap.String("sink-id", event.SinkID), zap.String("deployment-status", deploymentEntry.LastStatus))
	if deploymentEntry.LastStatus == "unknown" || deploymentEntry.LastStatus == "idle" {
		// async update sink status to provisioning
		go func() {
			err := d.deploymentService.UpdateStatus(ctx, event.OwnerID, event.SinkID, "provisioning", "")
			if err != nil {
				d.logger.Error("error updating status to provisioning", zap.Error(err))
			}
		}()
		_, err = d.deploymentService.NotifyCollector(ctx, event.OwnerID, event.SinkID, "deploy", "", "")
		if err != nil {
			d.logger.Error("error trying to notify collector", zap.Error(err))
			err2 := d.deploymentService.UpdateStatus(ctx, event.OwnerID, event.SinkID, "provisioning_error", err.Error())
			if err2 != nil {
				d.logger.Warn("error during notifying provisioning error, customer will not be notified of error")
				d.logger.Error("error during update provisioning error status", zap.Error(err))
				return err
			}
			return err
		}
		return nil
	} else {
		d.logger.Warn("collector is already running, skipping", zap.String("last_status", deploymentEntry.LastStatus))
		return nil
	}
}

func (d *eventService) HandleSinkIdle(ctx context.Context, event maestroredis.SinkerUpdateEvent) error {
	// check if exists deployment entry from postgres
	d.logger.Debug("handling sink idle event", zap.String("sink-id", event.SinkID), zap.String("owner-id", event.OwnerID))
	// async update sink status to idle
	go func() {
		err := d.deploymentService.UpdateStatus(ctx, event.OwnerID, event.SinkID, "idle", "")
		if err != nil {
			d.logger.Error("error updating status to idle", zap.Error(err))
		}
	}()
	// dropping idle otel collector
	_, err := d.deploymentService.NotifyCollector(ctx, event.OwnerID, event.SinkID, "delete", "idle", "")
	if err != nil {
		d.logger.Error("error trying to notify collector", zap.Error(err))
		err2 := d.deploymentService.UpdateStatus(ctx, event.OwnerID, event.SinkID, "provisioning_error", err.Error())
		if err2 != nil {
			d.logger.Warn("error during notifying provisioning error, customer will not be notified of error")
		}
		return err
	}

	return nil
}
