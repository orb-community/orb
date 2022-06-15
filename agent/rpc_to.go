/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"encoding/json"
	"github.com/ns1labs/orb/buildinfo"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

func (a *orbAgent) sendCapabilities() error {

	capabilities := fleet.Capabilities{
		SchemaVersion: fleet.CurrentCapabilitiesSchemaVersion,
		AgentTags:     a.config.OrbAgent.Tags,
		OrbAgent: fleet.OrbAgentInfo{
			Version: buildinfo.GetVersion(),
		},
	}

	capabilities.Backends = make(map[string]fleet.BackendInfo)
	for name, be := range a.backends {
		ver, err := be.Version()
		if err != nil {
			a.logger.Error("backend failed to retrieve version, skipping", zap.String("backend", name), zap.Error(err))
			continue
		}
		cp, err := be.GetCapabilities()
		if err != nil {
			a.logger.Error("backend failed to retrieve capabilities, skipping", zap.String("backend", name), zap.Error(err))
			continue
		}
		capabilities.Backends[name] = fleet.BackendInfo{
			Version: ver,
			Data:    cp,
		}
	}

	body, err := json.Marshal(capabilities)
	if err != nil {
		a.logger.Error("backend failed to marshal capabilities, skipping", zap.Error(err))
		return err
	}

	a.logger.Info("sending capabilities", zap.ByteString("value", body))
	if token := a.client.Publish(a.capabilitiesTopic, 1, false, body); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (a *orbAgent) sendGroupMembershipReq() error {
	a.logger.Debug("sending group membership request")
	payload := fleet.GroupMembershipReqRPCPayload{}

	data := fleet.RPC{
		SchemaVersion: fleet.CurrentRPCSchemaVersion,
		Func:          fleet.GroupMembershipReqRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if token := a.client.Publish(a.rpcToCoreTopic, 1, false, body); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil

}

func (a *orbAgent) checkGroupMembershipResponse() error {
	if a.groupRequestTicker == nil {
		a.groupRequestTicker = time.NewTicker(retryRequestFixedTime * retryRequestDuration)
	}
	for {
		calls := 0
		select {
		case <-a.groupRequestSucceeded:
			return nil
		case _ = <-a.groupRequestTicker.C:
			a.logger.Info("agent not received any group membership from fleet, re-requesting")
			duration := retryRequestFixedTime + (calls * retryDurationIncrPerAttempts)
			a.groupRequestTicker.Reset(time.Duration(duration) * retryRequestDuration)
			calls++
			err := a.sendGroupMembershipReq()
			if err != nil {
				a.logger.Error("failed to send group membership request", zap.Error(err))
			}
		}
	}
}

func (a *orbAgent) sendAgentPoliciesReq() error {
	a.logger.Debug("sending agent policies request")
	payload := fleet.AgentPoliciesReqRPCPayload{}

	data := fleet.RPC{
		SchemaVersion: fleet.CurrentRPCSchemaVersion,
		Func:          fleet.AgentPoliciesReqRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if token := a.client.Publish(a.rpcToCoreTopic, 1, false, body); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (a orbAgent) checkAgentPolicyResponse() {
	if a.policyRequestTicker == nil {
		a.policyRequestTicker = time.NewTicker(retryRequestFixedTime * retryRequestDuration)
	}
	for {
		calls := 0
		select {
		case <-a.policyRequestSucceeded:
		case _ = <-a.policyRequestTicker.C:
			a.logger.Info("agent not received any policy from fleet, re-requesting")
			duration := retryRequestFixedTime + (calls * retryDurationIncrPerAttempts)
			a.policyRequestTicker.Reset(time.Duration(duration) * retryRequestDuration)
			calls++
			err := a.sendAgentPoliciesReq()
			if err != nil {

			}
		}
	}
}
