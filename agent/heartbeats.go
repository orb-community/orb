/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"encoding/json"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

const HeartbeatFreq = 60 * time.Second

func (a *orbAgent) sendSingleHeartbeat(t time.Time, state fleet.State) {

	a.logger.Debug("heartbeat")

	hbData := fleet.Heartbeat{
		SchemaVersion: fleet.CurrentHeartbeatSchemaVersion,
		TimeStamp:     t,
		State:         state,
	}

	body, err := json.Marshal(hbData)
	if err != nil {
		a.logger.Error("error creating heartbeat data", zap.Error(err))
		return
	}

	if token := a.client.Publish(a.heartbeatsTopic, 1, false, body); token.Wait() && token.Error() != nil {
		a.logger.Error("error sending heartbeat", zap.Error(token.Error()))
	}
}

func (a *orbAgent) sendHeartbeats() {
	a.sendSingleHeartbeat(time.Now(), fleet.Online)
	for {
		select {
		case <-a.hbDone:
			return
		case t := <-a.hbTicker.C:
			a.sendSingleHeartbeat(t, fleet.Online)
		}
	}
}
