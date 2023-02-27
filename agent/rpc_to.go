/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/orb-community/orb/buildinfo"
	"github.com/orb-community/orb/fleet"
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

func (a *orbAgent) sendGroupMembershipRequest() error {
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

func (a *orbAgent) sendGroupMembershipReq() error {
	defer a.retryGroupMembershipRequest()
	return a.sendGroupMembershipRequest()
}

func (a *orbAgent) retryGroupMembershipRequest() {
	if a.groupRequestTicker == nil {
		a.groupRequestTicker = time.NewTicker(retryRequestFixedTime * retryRequestDuration)
	}
	var ctx context.Context
	ctx, a.groupRequestSucceeded = a.extendContext("retryGroupMembershipRequest")
	go func(ctx context.Context) {
		defer a.groupRequestTicker.Stop()
		defer func(t time.Time) {
			a.logger.Info("execution period of the re-request of retryGroupMembership", zap.Duration("waiting period", time.Now().Sub(t)))
		}(time.Now())
		lastT := time.Now()
		for calls := 1; calls <= retryMaxAttempts; calls++ {
			select {
			case <-ctx.Done():
				return
			case t := <-a.groupRequestTicker.C:
				a.logger.Info("agent did not receive any group membership from fleet, re-requesting", zap.Duration("waiting period", lastT.Sub(t)))
				duration := retryRequestFixedTime + (calls * retryDurationIncrPerAttempts)
				a.groupRequestTicker.Reset(time.Duration(duration) * retryRequestDuration)
				err := a.sendGroupMembershipRequest()
				if err != nil {
					a.logger.Error("failed to send group membership request", zap.Error(err))
					return
				}
				lastT = t
			}
		}
		a.logger.Warn(fmt.Sprintf("retryGroupMembership retried %d times and still got no response from fleet", retryMaxAttempts))
		return
	}(ctx)
}

func (a *orbAgent) sendAgentPoliciesReq() error {
	defer a.retryAgentPolicyResponse()
	return a.sendAgentPoliciesRequest()
}

func (a *orbAgent) sendAgentPoliciesRequest() error {
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

func (a *orbAgent) retryAgentPolicyResponse() {
	if a.policyRequestTicker == nil {
		a.policyRequestTicker = time.NewTicker(retryRequestFixedTime * retryRequestDuration)
	}
	var ctx context.Context
	ctx, a.policyRequestSucceeded = a.extendContext("retryAgentPolicyResponse")
	go func(ctx context.Context) {
		defer a.policyRequestTicker.Stop()
		defer func(t time.Time) {
			a.logger.Info("execution period of the re-request of retryAgentPolicy", zap.Duration("period", time.Now().Sub(t)))
		}(time.Now())
		lastT := time.Now()
		for calls := 1; calls <= retryMaxAttempts; calls++ {
			select {
			case <-ctx.Done():
				a.policyRequestTicker.Stop()
				return
			case t := <-a.policyRequestTicker.C:
				a.logger.Info("agent did not receive any policy from fleet, re-requesting", zap.Duration("waiting period", lastT.Sub(t)))
				duration := retryRequestFixedTime + (calls * retryDurationIncrPerAttempts)
				a.policyRequestTicker.Reset(time.Duration(duration) * retryRequestDuration)
				err := a.sendAgentPoliciesRequest()
				if err != nil {
					a.logger.Error("failed to send agent policies request", zap.Error(err))
					return
				}
				lastT = t
			}
		}
		a.logger.Warn(fmt.Sprintf("retryAgentPolicy retried %d times and still got no response from fleet", retryMaxAttempts))
		return
	}(ctx)
}
