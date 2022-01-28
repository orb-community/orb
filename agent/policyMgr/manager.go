/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package manager

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/agent/policies"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
)

type PolicyManager interface {
	ManagePolicy(payload fleet.AgentPolicyRPCPayload)
	RemovePolicyDataset(policyID string, datasetID string, be backend.Backend)
	GetPolicyState() ([]policies.PolicyData, error)
	GetRepo() policies.PolicyRepo
	ApplyBackendPolicies(be backend.Backend) error
	RemoveBackendPolicies(be backend.Backend) error
}

var _ PolicyManager = (*policyManager)(nil)

type policyManager struct {
	logger *zap.Logger
	config config.Config

	repo policies.PolicyRepo
}

func (a *policyManager) GetRepo() policies.PolicyRepo {
	return a.repo
}

func (a *policyManager) GetPolicyState() ([]policies.PolicyData, error) {
	return a.repo.GetAll()
}

func New(logger *zap.Logger, c config.Config, db *sqlx.DB) (PolicyManager, error) {
	repo, err := policies.NewMemRepo(logger)
	if err != nil {
		return nil, err
	}
	return &policyManager{logger: logger, config: c, repo: repo}, nil
}

func (a *policyManager) ManagePolicy(payload fleet.AgentPolicyRPCPayload) {

	a.logger.Info("managing agent policy from core",
		zap.String("action", payload.Action),
		zap.String("name", payload.Name),
		zap.String("dataset", payload.DatasetID),
		zap.String("backend", payload.Backend),
		zap.String("id", payload.ID),
		zap.Int32("version", payload.Version))

	switch payload.Action {
	case "manage":
		var pd = policies.PolicyData{
			ID:      payload.ID,
			Name:    payload.Name,
			Backend: payload.Backend,
			Version: payload.Version,
			Data:    payload.Data,
			State:   policies.Unknown,
		}
		var updatePolicy bool
		if a.repo.Exists(payload.ID) {
			// we have already processed this policy id before (it may be running or failed)
			// ensure we are associating this dataset with this policy, if one was specified
			// note the usual case is dataset id is NOT passed during policy updates
			if payload.DatasetID != "" {
				err := a.repo.EnsureDataset(payload.ID, payload.DatasetID)
				if err != nil {
					a.logger.Warn("policy failed to ensure dataset id", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name), zap.String("dataset_id", payload.DatasetID), zap.Error(err))
				}
			}
			// if policy already exist and has no version upgrade, has no need to apply it again
			currentPolicy, err := a.repo.Get(payload.ID)
			if err != nil {
				a.logger.Error("failed to retrieve policy", zap.String("policy_id", payload.ID), zap.Error(err))
				return
			}
			if currentPolicy.Version >= pd.Version {
				a.logger.Info("a better version of this policy has already been applied, skipping", zap.String("policy_id", pd.ID), zap.String("policy_name", pd.Name), zap.String("attempted_version", fmt.Sprint(pd.Version)), zap.String("current_version", fmt.Sprint(currentPolicy.Version)))
				return
			} else {
				updatePolicy = true
			}
		} else {
			// new policy we have not seen before, associate with this dataset
			// on first time we see policy, we *require* dataset
			if payload.DatasetID == "" {
				a.logger.Error("policy RPC for unseen policy did not include dataset ID, skipping", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name))
				return
			}
			pd.Datasets = map[string]bool{payload.DatasetID: true}
		}
		if !backend.HaveBackend(payload.Backend) {
			a.logger.Warn("policy failed to apply because backend is not available", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name))
			pd.State = policies.FailedToApply
			pd.BackendErr = "backend not available"
		} else {
			// attempt to apply the policy to the backend. status of policy application (running/failed) is maintained there.
			be := backend.GetBackend(payload.Backend)
			a.applyPolicy(payload, be, &pd, updatePolicy)
		}
		// save policy (with latest status) to local policy db
		a.repo.Update(pd)
		return
	case "remove":
		var pd = policies.PolicyData{
			ID:   payload.ID,
			Name: payload.Name,
		}
		if !backend.HaveBackend(payload.Backend) {
			a.logger.Warn("policy remove for a backend we do not have, ignoring", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name))
			return
		}
		be := backend.GetBackend(payload.Backend)
		// Remove policy via http request
		err := be.RemovePolicy(pd)
		if err != nil {
			a.logger.Warn("policy failed to remove", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name), zap.Error(err))
		}
		// Remove policy from orb-agent local repo
		err = a.repo.Remove(pd.ID)
		if err != nil {
			a.logger.Warn("policy failed to remove local", zap.String("policy_id", pd.ID), zap.String("policy_name", pd.Name), zap.Error(err))
		}
		break
	default:
		a.logger.Error("unknown policy action, ignored", zap.String("action", payload.Action))
	}

}

func (a *policyManager) RemovePolicyDataset(policyID string, datasetID string, be backend.Backend) {
	policyData, err := a.repo.Get(policyID)
	if err != nil {
		a.logger.Warn("failed to retrieve policy data", zap.String("policy_id", policyID), zap.String("policy_name", policyData.Name), zap.Error(err))
		return
	}
	removePolicy, err := a.repo.RemoveDataset(policyID, datasetID)
	if err != nil {
		a.logger.Warn("failed to remove policy dataset", zap.String("dataset_id", datasetID), zap.String("policy_name", policyData.Name), zap.Error(err))
		return
	}
	if removePolicy {
		// Remove policy via http request
		err := be.RemovePolicy(policyData)
		if err != nil {
			a.logger.Warn("policy failed to remove", zap.String("policy_id", policyID), zap.String("policy_name", policyData.Name), zap.Error(err))
		}
		// Remove policy from orb-agent local repo
		err = a.repo.Remove(policyData.ID)
		if err != nil {
			a.logger.Warn("policy failed to remove local", zap.String("policy_id", policyData.ID), zap.String("policy_name", policyData.Name), zap.Error(err))
		}
	}
}

func (a *policyManager) applyPolicy(payload fleet.AgentPolicyRPCPayload, be backend.Backend, pd *policies.PolicyData, updatePolicy bool) {
	err := be.ApplyPolicy(*pd, updatePolicy)
	if err != nil {
		a.logger.Warn("policy failed to apply", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name), zap.Error(err))
		pd.State = policies.FailedToApply
		pd.BackendErr = err.Error()
	} else {
		a.logger.Info("policy applied successfully", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name))
		pd.State = policies.Running
		pd.BackendErr = ""
	}
}

func (a *policyManager) RemoveBackendPolicies(be backend.Backend) error {
	plcies, err := a.repo.GetAll()
	if err != nil {
		a.logger.Error("failed to retrieve list of policies", zap.Error(err))
		return err
	}

	for _, plcy := range plcies {
		err := be.RemovePolicy(plcy)
		if err != nil {
			a.logger.Error("failed to remove policy from backend", zap.String("policy_id", plcy.ID), zap.String("policy_name", plcy.Name), zap.Error(err))
			return err
		}
		plcy.State = policies.Unknown
		a.repo.Update(plcy)
	}
	return nil
}

func (a *policyManager) ApplyBackendPolicies(be backend.Backend) error {
	plcies, err := a.repo.GetAll()
	if err != nil {
		a.logger.Error("failed to retrieve list of policies", zap.Error(err))
		return err
	}

	for _, policy := range plcies {
		be.ApplyPolicy(policy, false)
		if err != nil {
			a.logger.Warn("policy failed to apply", zap.String("policy_id", policy.ID), zap.String("policy_name", policy.Name), zap.Error(err))
			policy.State = policies.FailedToApply
			policy.BackendErr = err.Error()
		} else {
			a.logger.Info("policy applied successfully", zap.String("policy_id", policy.ID), zap.String("policy_name", policy.Name))
			policy.State = policies.Running
			policy.BackendErr = ""
		}
		a.repo.Update(policy)
	}
	return nil
}
