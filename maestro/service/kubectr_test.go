package service

import (
	"context"
	"github.com/orb-community/orb/maestro/kubecontrol"
	"go.uber.org/zap"
)

type testKubeCtr struct {
	logger *zap.Logger
}

func NewTestKubeCtr(logger *zap.Logger) kubecontrol.Service {
	return &testKubeCtr{logger: logger}
}

func (t *testKubeCtr) CreateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) (string, error) {
	name := "test-collector"
	return name, nil
}

func (t *testKubeCtr) KillOtelCollector(ctx context.Context, deploymentName, sinkID string) error {
	return nil
}
