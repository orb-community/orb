/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"encoding/json"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/policies"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

// HeartbeatFreq how often to heartbeat
const HeartbeatFreq = 50 * time.Second

// RestartTimeMin minimum time to wait between restarts
const RestartTimeMin = 5 * time.Minute

func (a *orbAgent) sendSingleHeartbeat(t time.Time, state fleet.State) {

	a.logger.Debug("heartbeat")

	bes := make(map[string]fleet.BackendStateInfo)
	for name, be := range a.backends {
		if state == fleet.Offline {
			bes[name] = fleet.BackendStateInfo{State: backend.Offline.String()}
		} else {
			state, errmsg, err := be.GetState()
			if err != nil {
				a.logger.Error("failed to retrieve backend state", zap.String("backend", name), zap.Error(err))
				bes[name] = fleet.BackendStateInfo{State: backend.AgentError.String(), Error: err.Error()}
				if time.Now().Sub(be.GetStartTime()) >= RestartTimeMin {
					a.logger.Info("attempting backend restart due to failed status during heartbeat")
					err := a.RestartBackend(name, "failed during heartbeat")
					if err != nil {
						a.logger.Error("failed to restart backend", zap.Error(err), zap.String("backend", name))
					}
				}
				continue
			}
			bes[name] = fleet.BackendStateInfo{State: state.String(), Error: errmsg}
		}
	}

	ps := make(map[string]fleet.PolicyStateInfo)
	pdata, err := a.policyManager.GetPolicyState()
	if err == nil {
		for _, pd := range pdata {
			if state == fleet.Offline {
				ps[pd.ID] = fleet.PolicyStateInfo{
					Name:     pd.Name,
					State:    policies.Offline.String(),
					Error:    pd.BackendErr,
					Datasets: pd.GetDatasetIDs(),
				}
			} else {
				ps[pd.ID] = fleet.PolicyStateInfo{
					Name:     pd.Name,
					State:    pd.State.String(),
					Error:    pd.BackendErr,
					Datasets: pd.GetDatasetIDs(),
				}
			}
		}
	} else {
		a.logger.Error("unable to retrieved policy state", zap.Error(err))
	}

	ag := make(map[string]fleet.GroupStateInfo)
	for id, groupInfo := range a.groupsInfos {
		ag[id] = fleet.GroupStateInfo{
			GroupName:    groupInfo.Name,
			GroupChannel: groupInfo.ChannelID,
		}
	}

	hbData := fleet.Heartbeat{
		SchemaVersion: fleet.CurrentHeartbeatSchemaVersion,
		State:         state,
		TimeStamp:     t,
		BackendState:  bes,
		PolicyState:   ps,
		GroupState:    ag,
	}

	body, err := json.Marshal(hbData)
	if err != nil {
		a.logger.Error("error marshalling heartbeat", zap.Error(err))
		return
	}

	if token := a.client.Publish(a.heartbeatsTopic, 1, false, body); token.Wait() && token.Error() != nil {
		a.logger.Error("error sending heartbeat", zap.Error(token.Error()))
		return
	}
}

func (a *orbAgent) sendHeartbeats() {
	a.logger.Debug("start heartbeats routine")
	a.sendSingleHeartbeat(time.Now(), fleet.Online)
	defer func() {
		a.logger.Debug("stopping heartbeats routine")
		a.sendSingleHeartbeat(time.Now(), fleet.Offline)
	}()
	for {
		select {
		case <-a.hbDone:
			return
		case t := <-a.hbTicker.C:
			a.sendSingleHeartbeat(t, fleet.Online)
		}
	}
}
