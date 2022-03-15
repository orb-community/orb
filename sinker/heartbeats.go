/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinker/config"
	"github.com/ns1labs/orb/sinker/redis/producer"
	"go.uber.org/zap"
	"time"
)

const (
	streamID       = "orb.sinker"
	streamLen      = 1000
	HeartbeatFreq  = 60 * time.Second
	DefaultTimeout = 30 * time.Minute
)

func (svc *sinkerService) sendSingleHeartbeat(t time.Time) {
	svc.logger.Info("heartbeat")
	configs, err := svc.sinkerCache.GetAll()
	if err != nil {
		svc.logger.Error("unable to retrieve policy state", zap.Error(err))
		return
	}
	for _, cfg := range configs {
		// Set idle if the sinker is more then 30 minutes not been sending metrics
		if cfg.LastRemoteWrite.Add(DefaultTimeout).Before(time.Now()) {
			cfg.State = config.Idle
			event := producer.SinkerUpdateEvent{
				SinkID:    cfg.SinkID,
				Owner:     cfg.OwnerID,
				State:     cfg.State.String(),
				Msg:       cfg.Msg,
				Timestamp: t,
			}
			record := &redis.XAddArgs{
				Stream:       streamID,
				MaxLenApprox: streamLen,
				Values:       event.Encode(),
			}
			err = svc.esclient.XAdd(context.Background(), record).Err()
			if err != nil {
				svc.logger.Error("error sending event to event store", zap.Error(err))
				continue
			}
		}
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
