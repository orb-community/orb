package deployment

import (
	"context"
	"errors"
	"time"

	"github.com/orb-community/orb/maestro/config"
	"github.com/orb-community/orb/maestro/kubecontrol"
	"github.com/orb-community/orb/maestro/password"
	"github.com/orb-community/orb/maestro/redis/producer"
	"github.com/orb-community/orb/pkg/types"
	"go.uber.org/zap"
)

const AuthenticationKey = "authentication"

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
	NotifyCollector(ctx context.Context, ownerID string, sinkId string, operation string, status string, errorMessage string) (string, error)
}

type deploymentService struct {
	dbRepository      Repository
	logger            *zap.Logger
	kafkaUrl          string
	maestroProducer   producer.Producer
	kubecontrol       kubecontrol.Service
	configBuilder     config.ConfigBuilder
	encryptionService password.EncryptionService
}

var _ Service = (*deploymentService)(nil)

func NewDeploymentService(logger *zap.Logger, repository Repository, kafkaUrl string, encryptionKey string,
	maestroProducer producer.Producer, kubecontrol kubecontrol.Service) Service {
	namedLogger := logger.Named("deployment-service")
	es := password.NewEncryptionService(logger, encryptionKey)
	cb := config.NewConfigBuilder(namedLogger, kafkaUrl, es)
	return &deploymentService{logger: namedLogger,
		dbRepository:      repository,
		configBuilder:     cb,
		encryptionService: es,
		maestroProducer:   maestroProducer,
		kubecontrol:       kubecontrol,
	}
}

func (d *deploymentService) CreateDeployment(ctx context.Context, deployment *Deployment) error {
	if deployment == nil {
		return errors.New("deployment is nil")
	}
	codedConfig, err := d.encodeConfig(deployment)
	if err != nil {
		return err
	}
	err = deployment.SetConfig(codedConfig)
	if err != nil {
		return err
	}
	// store with config encrypted
	added, err := d.dbRepository.Add(ctx, deployment)
	if err != nil {
		return err
	}
	d.logger.Info("added deployment", zap.String("id", added.Id),
		zap.String("ownerID", added.OwnerID), zap.String("sinkID", added.SinkID))
	err = d.maestroProducer.PublishSinkStatus(ctx, added.OwnerID, added.SinkID, "unknown", "")
	if err != nil {
		return err
	}
	return nil
}

func (d *deploymentService) getAuthBuilder(authType string) config.AuthBuilderService {
	return config.GetAuthService(authType, d.encryptionService)
}

func (d *deploymentService) encodeConfig(deployment *Deployment) (types.Metadata, error) {
	authType := deployment.GetConfig()
	if authType == nil {
		return nil, errors.New("deployment do not have authentication information")
	}
	value := authType.GetSubMetadata(AuthenticationKey)["type"].(string)
	authBuilder := d.getAuthBuilder(value)
	if authBuilder == nil {
		return nil, errors.New("deployment do not have authentication information")
	}
	return authBuilder.EncodeAuth(deployment.GetConfig())
}

func (d *deploymentService) GetDeployment(ctx context.Context, ownerID string, sinkId string) (*Deployment, string, error) {
	deployment, err := d.dbRepository.FindByOwnerAndSink(ctx, ownerID, sinkId)
	if err != nil {
		return nil, "", err
	}
	authType := deployment.GetConfig()
	if authType == nil {
		return nil, "", errors.New("deployment do not have authentication information")
	}
	value := authType.GetSubMetadata(AuthenticationKey)["type"].(string)
	authBuilder := d.getAuthBuilder(value)
	decodedDeployment, err := authBuilder.DecodeAuth(deployment.GetConfig())
	if err != nil {
		return nil, "", err
	}
	err = deployment.SetConfig(decodedDeployment)
	if err != nil {
		return nil, "", err
	}
	deployReq := &config.DeploymentRequest{
		OwnerID: ownerID,
		SinkID:  sinkId,
		Config:  deployment.GetConfig(),
		Backend: deployment.Backend,
		Status:  deployment.LastStatus,
	}
	manifest, err := d.configBuilder.BuildDeploymentConfig(deployReq)
	if err != nil {
		return nil, "", err
	}
	return deployment, manifest, nil
}

// UpdateDeployment will stop the running collector if any, and change the deployment, it will not spin the collector back up,
// it will wait for the next sink.activity
func (d *deploymentService) UpdateDeployment(ctx context.Context, deployment *Deployment) error {
	now := time.Now()
	got, _, err := d.GetDeployment(ctx, deployment.OwnerID, deployment.SinkID)
	if err != nil {
		return errors.New("could not find deployment to update")
	}
	// Spin down the collector if it is running
	err = d.kubecontrol.KillOtelCollector(ctx, got.CollectorName, got.SinkID)
	if err != nil {
		d.logger.Warn("could not stop running collector, will try to update anyway", zap.Error(err))
	}
	err = deployment.Merge(*got)
	if err != nil {
		d.logger.Error("error during merge of deployments", zap.Error(err))
		return err
	}
	deployment.LastCollectorStopTime = &now
	codedConfig, err := d.encodeConfig(deployment)
	if err != nil {
		return err
	}
	err = deployment.SetConfig(codedConfig)
	if err != nil {
		return err
	}
	updated, err := d.dbRepository.Update(ctx, deployment)
	if err != nil {
		return err
	}
	d.logger.Info("updated deployment", zap.String("ownerID", updated.OwnerID),
		zap.String("sinkID", updated.SinkID))
	return nil
}

func (d *deploymentService) NotifyCollector(ctx context.Context, ownerID string, sinkId string, operation string,
	status string, errorMessage string) (string, error) {
	got, manifest, err := d.GetDeployment(ctx, ownerID, sinkId)
	if err != nil {
		return "", errors.New("could not find deployment to update")
	}
	now := time.Now()
	if operation == "delete" {
		got.LastCollectorStopTime = &now
		err = d.kubecontrol.KillOtelCollector(ctx, got.CollectorName, got.SinkID)
		if err != nil {
			d.logger.Warn("could not stop running collector, will try to update anyway", zap.Error(err))
		}
	} else if operation == "deploy" {
		// Spin up the collector
		if got.LastCollectorDeployTime == nil || got.LastCollectorDeployTime.Before(now) {
			if got.LastCollectorStopTime == nil || got.LastCollectorStopTime.Before(now) {
				d.logger.Debug("collector is not running deploying")
				got.CollectorName, err = d.kubecontrol.CreateOtelCollector(ctx, got.OwnerID, got.SinkID, manifest)
				got.LastCollectorDeployTime = &now
			} else {
				d.logger.Info("collector is already running")
			}
		}

	}
	if status != "" {
		got.LastStatus = status
		got.LastStatusUpdate = &now
	}
	if errorMessage != "" {
		got.LastErrorMessage = errorMessage
		got.LastErrorTime = &now
	}
	codedConfig, err := d.encodeConfig(got)
	if err != nil {
		return "", err
	}
	err = got.SetConfig(codedConfig)
	if err != nil {
		return "", err
	}
	updated, err := d.dbRepository.Update(ctx, got)
	if err != nil {
		return "", err
	}
	d.logger.Info("updated deployment information for collector and status or error",
		zap.String("ownerID", updated.OwnerID), zap.String("sinkID", updated.SinkID),
		zap.String("collectorName", updated.CollectorName),
		zap.String("status", updated.LastStatus), zap.String("errorMessage", updated.LastErrorMessage))
	return updated.CollectorName, nil
}

// UpdateStatus this will change the status in postgres and notify sinks service to show new status to user
func (d *deploymentService) UpdateStatus(ctx context.Context, ownerID string, sinkId string, status string, errorMessage string) error {
	got, _, err := d.GetDeployment(ctx, ownerID, sinkId)
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

	codedConfig, err := d.encodeConfig(got)
	if err != nil {
		return err
	}
	err = got.SetConfig(codedConfig)
	if err != nil {
		return err
	}
	updated, err := d.dbRepository.Update(ctx, got)
	if err != nil {
		return err
	}
	d.logger.Info("updated deployment status",
		zap.String("ownerID", updated.OwnerID), zap.String("sinkID", updated.SinkID),
		zap.String("status", updated.LastStatus), zap.String("errorMessage", updated.LastErrorMessage))
	err = d.maestroProducer.PublishSinkStatus(ctx, updated.OwnerID, updated.SinkID, updated.LastStatus, "")
	if err != nil {
		return err
	}
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
