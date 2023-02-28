/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package manager

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/orb-community/orb/agent/backend"
	"github.com/orb-community/orb/agent/config"
	"github.com/orb-community/orb/agent/policies"
	"github.com/orb-community/orb/fleet"
	"github.com/orb-community/orb/pkg/errors"
	"go.uber.org/zap"
)

type PolicyManager interface {
	ManagePolicy(payload fleet.AgentPolicyRPCPayload)
	RemovePolicyDataset(policyID string, datasetID string, be backend.Backend)
	GetPolicyState() ([]policies.PolicyData, error)
	GetRepo() policies.PolicyRepo
	ApplyBackendPolicies(be backend.Backend) error
	RemoveBackendPolicies(be backend.Backend, permanently bool) error
	RemovePolicy(policyID string, policyName string, beName string) error
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

			if payload.AgentGroupID != "" {
				err := a.repo.EnsureGroupID(payload.ID, payload.AgentGroupID)
				if err != nil {
					a.logger.Warn("policy failed to ensure agent group id", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name), zap.String("agent_group_id", payload.AgentGroupID), zap.Error(err))
				}
			}

			// if policy already exist and has no version upgrade, has no need to apply it again
			currentPolicy, err := a.repo.Get(payload.ID)
			if err != nil {
				a.logger.Error("failed to retrieve policy", zap.String("policy_id", payload.ID), zap.Error(err))
				return
			}
			if currentPolicy.Version >= pd.Version && currentPolicy.State == policies.Running {
				a.logger.Info("a better version of this policy has already been applied, skipping", zap.String("policy_id", pd.ID), zap.String("policy_name", pd.Name), zap.String("attempted_version", fmt.Sprint(pd.Version)), zap.String("current_version", fmt.Sprint(currentPolicy.Version)))
				return
			} else {
				updatePolicy = true
			}
			if currentPolicy.Name != pd.Name {
				pd.PreviousPolicyData = &policies.PolicyData{Name: currentPolicy.Name}
			}
			pd.Datasets = currentPolicy.Datasets
			pd.GroupIds = currentPolicy.GroupIds
		} else {
			// new policy we have not seen before, associate with this dataset
			// on first time we see policy, we *require* dataset
			if payload.DatasetID == "" {
				a.logger.Error("policy RPC for unseen policy did not include dataset ID, skipping", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name))
				return
			}
			pd.Datasets = map[string]bool{payload.DatasetID: true}

			if payload.AgentGroupID != "" {
				pd.GroupIds = map[string]bool{payload.AgentGroupID: true}
			}

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
		err := a.repo.Update(pd)
		if err != nil {
			a.logger.Error("got error in update last status", zap.Error(err))
			return
		}
		return
	case "remove":
		err := a.RemovePolicy(payload.ID, payload.Name, payload.Backend)
		if err != nil {
			a.logger.Error("policy failed to be removed", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name), zap.Error(err))
		}
		break
	default:
		a.logger.Error("unknown policy action, ignored", zap.String("action", payload.Action))
	}
}

func (a *policyManager) RemovePolicy(policyID string, policyName string, beName string) error {
	var pd = policies.PolicyData{
		ID:   policyID,
		Name: policyName,
	}
	if !backend.HaveBackend(beName) {
		return errors.New("policy remove for a backend we do not have, ignoring")
	}
	be := backend.GetBackend(beName)
	err := be.RemovePolicy(pd)
	if err != nil {
		a.logger.Error("backend remove policy failed: will still remove from PolicyManager", zap.String("policy_id", policyID), zap.Error(err))
	}
	// Remove policy from orb-agent local repo
	err = a.repo.Remove(pd.ID)
	if err != nil {
		return err
	}
	return nil
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
		switch {
		case strings.Contains(err.Error(), "422"):
			pd.State = policies.NoTapMatch
		default:
			pd.State = policies.FailedToApply
		}
		pd.BackendErr = err.Error()
	} else {
		a.logger.Info("policy applied successfully", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name))
		pd.State = policies.Running
		pd.BackendErr = ""
	}
}

func (a *policyManager) RemoveBackendPolicies(be backend.Backend, permanently bool) error {
	plcies, err := a.repo.GetAll()
	if err != nil {
		a.logger.Error("failed to retrieve list of policies", zap.Error(err))
		return err
	}

	for _, plcy := range plcies {
		err := be.RemovePolicy(plcy)
		if err != nil {
			a.logger.Error("failed to remove policy from backend", zap.String("policy_id", plcy.ID), zap.String("policy_name", plcy.Name), zap.Error(err))
			// note we continue here: even if the backend failed to remove, we update our policy repo to remove it
		}
		if permanently {
			err = a.repo.Remove(plcy.ID)
			if err != nil {
				return err
			}
		} else {
			plcy.State = policies.Unknown
			err = a.repo.Update(plcy)
			if err != nil {
				return err
			}
		}
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
		err := be.ApplyPolicy(policy, false)
		if err != nil {
			a.logger.Warn("policy failed to apply", zap.String("policy_id", policy.ID), zap.String("policy_name", policy.Name), zap.Error(err))
			policy.State = policies.FailedToApply
			policy.BackendErr = err.Error()
		} else {
			a.logger.Info("policy applied successfully", zap.String("policy_id", policy.ID), zap.String("policy_name", policy.Name))
			policy.State = policies.Running
			policy.BackendErr = ""
		}
		err = a.repo.Update(policy)
		if err != nil {
			return err
		}
	}
	return nil
}
