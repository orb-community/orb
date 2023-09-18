package config

import (
	"github.com/orb-community/orb/maestro/deployment"
	"go.uber.org/zap"
)

type ConfigBuilder interface {
	BuildDeploymentConfig(deployment *deployment.Deployment) (string, error)
}

type configBuilder struct {
	logger            *zap.Logger
	kafkaUrl          string
	encryptionService deployment.EncryptionService
}

var _ ConfigBuilder = (*configBuilder)(nil)

func NewConfigBuilder(logger *zap.Logger, kafkaUrl string, encryptionService deployment.EncryptionService) ConfigBuilder {
	return &configBuilder{logger: logger, kafkaUrl: kafkaUrl, encryptionService: encryptionService}
}
