package fleet

import (
	"time"
)

const (
	HeartbeatFreq  = 60 * time.Second
	DefaultTimeout = 30 * time.Minute
)

func (svc *fleetService) checkState(t time.Time) {
	svc.logger.Info("checking stale agents...")
	// TODO Create a query to update to stale agents with more than 30 minutes without a heartbeat
}


func (svc *fleetService) checkAgents() {
	svc.checkState(time.Now())
	for {
		select {
		case <-svc.aDone:
			svc.logger.Info("stopping ticker")
			return
		case t := <-svc.aTicker.C:
			svc.checkState(t)
		}
	}
}
