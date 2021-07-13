/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"encoding/json"
	"github.com/ns1labs/orb"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
)

func (a *orbAgent) sendCapabilities() error {

	capabilities := fleet.Capabilities{
		SchemaVersion: fleet.CurrentCapabilitiesSchemaVersion,
		AgentTags:     a.config.OrbAgent.Tags,
		OrbAgent: fleet.OrbAgentInfo{
			Version: orb.GetVersion(),
		},
	}

	capabilities.Backends = make(map[string]fleet.BackendInfo)
	for name, be := range a.backends {
		ver, err := be.Version()
		if err != nil {
			a.logger.Error("backend failed to retrieve version", zap.String("backend", name), zap.Error(err))
			continue
		}
		capabilities.Backends[name] = fleet.BackendInfo{
			Version: ver,
		}
	}

	body, err := json.Marshal(capabilities)
	if err != nil {
		return err
	}

	if token := a.client.Publish(a.capabilitiesTopic, 1, false, body); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (a *orbAgent) sendGroupMembershipReq() error {

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

func (a *orbAgent) sendAgentPoliciesReq() error {

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
