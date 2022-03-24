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
	CheckerFreq    = 300 * time.Second
	DefaultTimeout = 30 * time.Minute
)

func (svc *sinkerService) checkState(t time.Time) {
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
			// Set idle if the sinker is more then 30 minutes not been sending metrics
			if cfg.LastRemoteWrite.Add(DefaultTimeout).Before(time.Now()) {
				cfg.State = config.Idle
				if err := svc.sinkerCache.Edit(cfg); err != nil {
					svc.logger.Error("error updating sink config cache", zap.Error(err))
					return
				}
			}
		}
	}
}

func (svc *sinkerService) checkSinker() {
	svc.checkState(time.Now())
	for {
		select {
		case <-svc.hbDone:
			return
		case t := <-svc.hbTicker.C:
			svc.checkState(t)
		}
	}
}
