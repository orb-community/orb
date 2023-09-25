package config

import (
	"github.com/orb-community/orb/maestro/password"
	"github.com/orb-community/orb/pkg/types"
	"go.uber.org/zap"
)

type ConfigBuilder interface {
	BuildDeploymentConfig(deployment *DeploymentRequest) (string, error)
}

type DeploymentRequest struct {
	OwnerID string
	SinkID  string
	Config  types.Metadata
	Backend string
	Status  string
}

type configBuilder struct {
	logger            *zap.Logger
	kafkaUrl          string
	encryptionService password.EncryptionService
}

var _ ConfigBuilder = (*configBuilder)(nil)

func NewConfigBuilder(logger *zap.Logger, kafkaUrl string, encryptionService password.EncryptionService) ConfigBuilder {
	return &configBuilder{logger: logger, kafkaUrl: kafkaUrl, encryptionService: encryptionService}
}
