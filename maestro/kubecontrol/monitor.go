package kubecontrol

import (
	"context"
	"go.uber.org/zap"
	"time"
)

func (svc *deployService) StartMonitor(ctx context.Context, cancelFunc context.CancelFunc) {
	ticker := time.NewTicker(5 * time.Minute)
	svc.logger.Debug("start monitor routine", zap.Any("routine", ctx.Value("#routine")))
	defer func() {
		cancelFunc()
		svc.logger.Debug("stopping monitor routine")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case _ = <-ticker.C:
			svc.monitorSink()
		}
	}
}

func (svc *deployService) monitorSink() {
	for sink, up := range svc.deploymentState {
		if up {
			svc.logger.Info("check sink", zap.Any("sink", sink))
		}
	}
}
