package fleet

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const (
	HeartbeatFreq  = 60 * time.Second
	DefaultTimeout = 300 * time.Second
)

func (svc *fleetService) checkState(t time.Time) {
	svc.logger.Info("checking for stale agents")
	count, err := svc.agentRepo.SetStaleStatus(context.Background(), DefaultTimeout)
	if err != nil {
		svc.logger.Error("failed to change agents status to stale", zap.Error(err))
	}
	if count > 0 {
		svc.logger.Info(fmt.Sprintf("%d agents with more than %v without heartbeats had their state changed to stale", count, DefaultTimeout))
	}
}


func (svc *fleetService) checkAgents() {
	svc.checkState(time.Now())
	for {
		select {
		case <-svc.aDone:
			svc.logger.Info("stopping stale agent routine")
			return
		case t := <-svc.aTicker.C:
			svc.checkState(t)
		}
	}
}
