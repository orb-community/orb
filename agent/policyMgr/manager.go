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

	if !backend.HaveBackend(payload.Backend) {
		a.logger.Warn("policy for a backend we do not have, ignoring", zap.String("id", payload.ID))
		return
	}

	be := backend.GetBackend(payload.Backend)

	switch payload.Action {
	case "manage":
		if a.repo.Exists(payload.ID) {
			a.logger.Info("policy already exists, ensuring dataset", zap.String("id", payload.ID), zap.String("dataset_id", payload.DatasetID))
			err := a.repo.EnsureDataset(payload.ID, payload.DatasetID)
			if err != nil {
				a.logger.Warn("policy failed to ensure dataset id", zap.String("id", payload.ID), zap.String("dataset_id", payload.DatasetID), zap.Error(err))
			}
			pd, err := applyPolicy(payload, be, a)
			if err != nil {
				a.repo.Edit(pd)
			}
		} else {
			pd, err := applyPolicy(payload, be, a)
			if err != nil {
				a.repo.Add(pd)
			}
		}
		return
	case "remove":
		err := be.RemovePolicy(payload.ID)
		if err != nil {
			a.logger.Warn("policy failed to remove", zap.String("id", payload.ID), zap.Error(err))
		}
		break
	default:
		a.logger.Error("unknown policy action, ignored", zap.String("action", payload.Action))
	}

}

func applyPolicy(payload fleet.AgentPolicyRPCPayload, be backend.Backend, a *policyManager) (policies.PolicyData, error) {
	pd := policies.PolicyData{
		ID:      payload.ID,
		Name:    payload.Name,
		Backend: payload.Backend,
		Version: payload.Version,
		Data:    payload.Data,
		State:   policies.Unknown,
	}
	err := be.ApplyPolicy(pd)
	if err != nil {
		a.logger.Warn("policy failed to apply", zap.String("id", payload.ID), zap.Error(err))
		pd.State = policies.FailedToApply
		pd.BackendErr = err.Error()
		return pd, err
	} else {
		pd.State = policies.Running
	}
	return pd, nil
}
