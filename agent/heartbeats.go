/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"encoding/json"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

const HeartbeatFreq = 60 * time.Second

func (a *orbAgent) sendSingleHeartbeat(t time.Time, state fleet.State) {

	a.logger.Debug("heartbeat")

	bes := make(map[string]fleet.BackendStateInfo)
	for name, be := range a.backends {
		state, errmsg, err := be.GetState()
		if err != nil {
			a.logger.Error("failed to retrieve backend state", zap.String("backend", name), zap.Error(err))
			bes[name] = fleet.BackendStateInfo{State: backend.AgentError.String(), Error: err.Error()}
			continue
		}
		bes[name] = fleet.BackendStateInfo{State: state.String(), Error: errmsg}
	}

	ps := make(map[string]fleet.PolicyStateInfo)
	pdata, err := a.policyManager.GetPolicyState()
	if err == nil {
		for _, pd := range pdata {
			ps[pd.ID] = fleet.PolicyStateInfo{
				State: pd.State.String(),
				Error: pd.BackendErr,
			}
		}
	} else {
		a.logger.Error("unable to retrieved policy state", zap.Error(err))
	}

	hbData := fleet.Heartbeat{
		SchemaVersion: fleet.CurrentHeartbeatSchemaVersion,
		TimeStamp:     t,
		BackendState:  bes,
		PolicyState:   ps,
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
