/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
	"time"
)

const (
	streamID       = "orb.sinker"
	streamLen      = 1000
	CheckerFreq    = 5 * time.Minute
	DefaultTimeout = 30 * time.Minute
)

func (svc *SinkerService) checkState(_ time.Time) {
	owners, err := svc.sinkerCache.GetAllOwners()
	if err != nil {
		svc.logger.Error("failed to retrieve the list of owners")
		return
	}

	for _, ownerID := range owners {
		configs, err := svc.sinkerCache.GetAll(ownerID)
		if err != nil {
			svc.logger.Error("unable to retrieve policy state", zap.Error(err))
			return
		}
		for _, cfg := range configs {
			// Set idle if the sinker is more than 30 minutes not sending metrics (Remove from Redis)
			if cfg.LastRemoteWrite.Add(DefaultTimeout).Before(time.Now()) {
				if cfg.State == config.Active {
					if err := svc.sinkerCache.Remove(cfg.OwnerID, cfg.SinkID); err != nil {
						svc.logger.Error("error updating sink config cache", zap.Error(err))
						return
					}
				}
			}
		}
	}
}

func (svc *SinkerService) checkSinker() {
	svc.checkState(time.Now())
	for {
		select {
		case <-svc.hbDone:
			svc.otelCancelFunct()
			svc.cancelAsyncContext()
			return
		case t := <-svc.hbTicker.C:
			svc.checkState(t)
		}
	}
}
