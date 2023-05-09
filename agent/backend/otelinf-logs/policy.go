/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otelinf

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/orb-community/orb/agent/policies"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func (d *otelinfBackend) ApplyPolicy(data policies.PolicyData, updatePolicy bool) error {
	if updatePolicy {
		// To update a policy it's necessary first remove it and then apply a new version
		err := d.RemovePolicy(data)
		if err != nil {
			d.logger.Warn("policy failed to remove", zap.String("policy_id", data.ID), zap.String("policy_name", data.Name), zap.Error(err))
		}
	}
	d.logger.Debug("diode-agent policy apply", zap.String("policy_id", data.ID), zap.Any("data", data.Data))
	pol := map[string]interface{}{
		data.Name: data.Data,
	}
	policyYaml, err := yaml.Marshal(pol)
	if err != nil {
		d.logger.Warn("yaml policy marshal failure", zap.String("policy_id", data.ID), zap.Any("policy", pol))
		return err
	}
	var resp map[string]interface{}
	err = d.request("policies", &resp, http.MethodPost, bytes.NewBuffer(policyYaml), "application/x-yaml", ApplyPolicyTimeout)
	if err != nil {
		d.logger.Warn("yaml policy application failure", zap.String("policy_id", data.ID), zap.ByteString("policy", policyYaml))
		return err
	}
	return nil
}

func (d *otelinfBackend) RemovePolicy(data policies.PolicyData) error {
	d.logger.Debug("diode-agent policy remove", zap.String("policy_id", data.ID))
	var resp interface{}
	var name string
	// Since we use Name for removing policies not IDs, if there is a change, we need to remove the previous name of the policy
	if data.PreviousPolicyData != nil && data.PreviousPolicyData.Name != data.Name {
		name = data.PreviousPolicyData.Name
	} else {
		name = data.Name
	}
	err := d.request(fmt.Sprintf("policies/%s", name), &resp, http.MethodDelete, http.NoBody, "application/json", RemovePolicyTimeout)
	if err != nil {
		return err
	}
	return nil
}
