/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

const HeartbeatFreq = 60 * time.Second

func (svc *sinkerService) sendSingleHeartbeat(t time.Time) {
	configs, err := svc.configRepo.GetAll()
	if err != nil {
		svc.logger.Error("unable to retrieved policy state", zap.Error(err))
		return
	}
	for cfg := range configs {
		// Todo implement a redis event
		fmt.Print(cfg)
	}
}

func (svc *sinkerService) sendHeartbeats() {
	svc.sendSingleHeartbeat(time.Now())
	for {
		select {
		case <-svc.hbDone:
			return
		case t := <-svc.hbTicker.C:
			svc.sendSingleHeartbeat(t)
		}
	}
}
