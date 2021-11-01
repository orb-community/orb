/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package manager

import (
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/agent/policies"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
)

type PolicyManager interface {
	ManagePolicy(payload fleet.AgentPolicyRPCPayload)
	GetPolicyState() ([]policies.PolicyData, error)
}

var _ PolicyManager = (*policyManager)(nil)

type policyManager struct {
	logger *zap.Logger
	config config.Config

	repo policies.PolicyRepo
}

func (a *policyManager) GetPolicyState() ([]policies.PolicyData, error) {
	d, e := a.repo.GetAll()
	return d, e
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
		if a.repo.Exists(payload.ID) {
			// we have already processed this policy id before (it may be running or failed)
			// ensure we are associating this dataset with this policy
			err := a.repo.EnsureDataset(payload.ID, payload.DatasetID)
			if err != nil {
				a.logger.Warn("policy failed to ensure dataset id", zap.String("policy_id", payload.ID), zap.String("dataset_id", payload.DatasetID), zap.Error(err))
			}
		} else {
			// new policy we have not seen before, associate with this dataset
			pd.Datasets = map[string]bool{payload.DatasetID: true}
		}
		if !backend.HaveBackend(payload.Backend) {
			a.logger.Warn("policy failed to apply because backend is not available", zap.String("policy_id", payload.ID), zap.String("policy_name", payload.Name))
			pd.State = policies.FailedToApply
			pd.BackendErr = "backend not available"
		} else {
			// attempt to apply the policy to the backend. status of policy application (running/failed) is maintained there.
			be := backend.GetBackend(payload.Backend)
			a.applyPolicy(payload, be, &pd)
		}
		// save policy (with latest status) to local policy db
		a.repo.Update(pd)
		return
	case "remove":
		if !backend.HaveBackend(payload.Backend) {
			a.logger.Warn("policy remove for a backend we do not have, ignoring", zap.String("policy_id", payload.ID))
			return
		}
		be := backend.GetBackend(payload.Backend)
		err := be.RemovePolicy(payload.ID)
		if err != nil {
			a.logger.Warn("policy failed to remove", zap.String("policy_id", payload.ID), zap.Error(err))
		}
		break
	default:
		a.logger.Error("unknown policy action, ignored", zap.String("action", payload.Action))
	}

}

func (a *policyManager) applyPolicy(payload fleet.AgentPolicyRPCPayload, be backend.Backend, pd *policies.PolicyData) {
	err := be.ApplyPolicy(*pd)
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
