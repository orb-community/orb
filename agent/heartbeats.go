/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/orb-community/orb/agent/backend"
	"github.com/orb-community/orb/agent/policies"
	"github.com/orb-community/orb/fleet"
	"go.uber.org/zap"
	"time"
)

// HeartbeatFreq how often to heartbeat
const HeartbeatFreq = 50 * time.Second

// RestartTimeMin minimum time to wait between restarts
const RestartTimeMin = 5 * time.Minute

func (a *orbAgent) sendSingleHeartbeat(ctx context.Context, t time.Time, agentsState fleet.State) {

	if a.heartbeatsTopic == "" {
		a.logger.Debug("heartbeat topic not yet set, skipping")
		return
	}

	a.logger.Debug("heartbeat", zap.String("state", agentsState.String()))

	bes := make(map[string]fleet.BackendStateInfo)
	for name, be := range a.backends {
		if agentsState == fleet.Offline {
			bes[name] = fleet.BackendStateInfo{State: backend.Offline.String()}
			continue
		}
		besi := fleet.BackendStateInfo{}
		backendStatus, errMsg, err := be.GetRunningStatus()
		a.backendState[name].Status = backendStatus
		besi.State = backendStatus.String()
		if backendStatus != backend.Running {
			a.logger.Error("backend not ready", zap.String("backend", name), zap.String("status", backendStatus.String()), zap.String("errMsg", errMsg), zap.Error(err))
			if err != nil {
				a.backendState[name].LastError = fmt.Sprintf("failed to retrieve backend status: %v", err)
			} else if errMsg != "" {
				a.backendState[name].LastError = errMsg
			}
			// status is not running so we have a current error
			besi.Error = a.backendState[name].LastError
			if time.Now().Sub(be.GetStartTime()) >= RestartTimeMin {
				a.logger.Info("attempting backend restart due to failed status during heartbeat")
				if a.config.OrbAgent.Cloud.MQTT.Id != "" {
					ctx = context.WithValue(ctx, "agent_id", a.config.OrbAgent.Cloud.MQTT.Id)
				} else {
					ctx = context.WithValue(ctx, "agent_id", "auto-provisioning-without-id")
				}
				err := a.RestartBackend(ctx, name, "failed during heartbeat")
				if err != nil {
					a.logger.Error("failed to restart backend", zap.Error(err), zap.String("backend", name))
				}
			} else {
				a.logger.Info("waiting to attempt backend restart due to failed status", zap.Duration("remaining_secs", RestartTimeMin-(time.Now().Sub(be.GetStartTime()))))
			}
		} else {
			// status is Running so no current error
			besi.Error = ""
		}
		if a.backendState[name].LastError != "" {
			besi.LastError = a.backendState[name].LastError
		}
		if !a.backendState[name].LastRestartTS.IsZero() {
			besi.LastRestartTS = a.backendState[name].LastRestartTS
		}
		if a.backendState[name].RestartCount > 0 {
			besi.RestartCount = a.backendState[name].RestartCount
		}
		if a.backendState[name].LastRestartReason != "" {
			besi.LastRestartReason = a.backendState[name].LastRestartReason
		}
		bes[name] = besi
	}

	ps := make(map[string]fleet.PolicyStateInfo)
	pdata, err := a.policyManager.GetPolicyState()
	if err == nil {
		for _, pd := range pdata {
			pstate := policies.Offline.String()
			// if agent is not offline, default to status that policy manager believes we should be in
			if agentsState != fleet.Offline {
				pstate = pd.State.String()
			}
			// but if the policy backend is not running, policy isn't either
			if bestate, ok := a.backendState[pd.Backend]; ok && bestate.Status != backend.Running {
				pstate = policies.Unknown.String()
				pd.BackendErr = "backend is unreachable"
			}
			ps[pd.ID] = fleet.PolicyStateInfo{
				Name:            pd.Name,
				Version:         pd.Version,
				State:           pstate,
				Error:           pd.BackendErr,
				Datasets:        pd.GetDatasetIDs(),
				LastScrapeTS:    pd.LastScrapeTS,
				LastScrapeBytes: pd.LastScrapeBytes,
				Backend:         pd.Backend,
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
		State:         agentsState,
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
		err = a.restartComms(ctx)
		if err != nil {
			a.logger.Error("error reconnecting with MQTT, stopping agent")
			a.Stop(ctx)
		}
	}
}

func (a *orbAgent) sendHeartbeats(ctx context.Context, cancelFunc context.CancelFunc) {
	a.logger.Debug("start heartbeats routine", zap.Any("routine", ctx.Value("routine")))
	a.sendSingleHeartbeat(ctx, time.Now(), fleet.Online)
	defer func() {
		cancelFunc()
	}()
	for {
		select {
		case <-ctx.Done():
			a.logger.Debug("context done, stopping heartbeats routine")
			a.sendSingleHeartbeat(ctx, time.Now(), fleet.Offline)
			a.heartbeatCtx = nil
			return
		case t := <-a.hbTicker.C:
			a.sendSingleHeartbeat(ctx, t, fleet.Online)
		}
	}
}
