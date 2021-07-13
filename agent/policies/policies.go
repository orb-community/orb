/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
)

type PolicyManager interface {
	ManagePolicy(payload fleet.AgentPolicyRPCPayload)
}

var _ PolicyManager = (*policyManager)(nil)

type policyManager struct {
	logger *zap.Logger
	config config.Config
	db     *sqlx.DB
}

func New(logger *zap.Logger, c config.Config) (PolicyManager, error) {
	db, err := sqlx.Connect("sqlite3", "orb-agent.db")
	if err != nil {
		return nil, err
	}
	return &policyManager{logger: logger, config: c, db: db}, nil
}

func (a *policyManager) ManagePolicy(payload fleet.AgentPolicyRPCPayload) {

	a.logger.Info("managing agent policy from core",
		zap.String("name", payload.Name),
		zap.String("backend", payload.Backend),
		zap.String("id", payload.ID),
		zap.Int32("version", payload.Version))

	if !backend.HaveBackend(payload.Backend) {
		a.logger.Warn("policy for a backend we do not have, ignoring", zap.String("id", payload.ID))
		return
	}

	be := backend.GetBackend(payload.Backend)
	err := be.ApplyPolicy(payload.Data)
	if err != nil {
		a.logger.Warn("policy failed to apply", zap.String("id", payload.ID), zap.Error(err))
	}

}
