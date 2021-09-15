/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"database/sql/driver"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
)

const (
	Unknown State = iota
	Applied
	Failed
)

type State int

type policyData struct {
	ID         string
	Datasets   []string
	Name       string
	Backend    string
	Version    int32
	Data       interface{}
	State      State
	BackendErr string
}

var stateMap = [...]string{
	"unknown",
	"applied",
	"failed",
}

var stateRevMap = map[string]State{
	"unknown": Unknown,
	"applied": Applied,
	"failed":  Failed,
}

func (s State) String() string {
	return stateMap[s]
}

func (s *State) Scan(value interface{}) error { *s = stateRevMap[string(value.([]byte))]; return nil }
func (s State) Value() (driver.Value, error)  { return s.String(), nil }

type PolicyManager interface {
	ManagePolicy(payload fleet.AgentPolicyRPCPayload)
}

var _ PolicyManager = (*policyManager)(nil)

type policyManager struct {
	logger *zap.Logger
	config config.Config
	db     *sqlx.DB

	repo PolicyRepo
}

func New(logger *zap.Logger, c config.Config, db *sqlx.DB) (PolicyManager, error) {
	return &policyManager{logger: logger, config: c, db: db}, nil
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
			err := a.repo.EnsureDataset(payload.ID, payload.DatasetID)
			if err != nil {
				a.logger.Warn("policy failed to ensure dataset id", zap.String("id", payload.ID), zap.String("dataset_id", payload.DatasetID), zap.Error(err))
			}
		} else {
			pd := policyData{
				ID:      payload.ID,
				Name:    payload.Name,
				Backend: payload.Backend,
				Version: payload.Version,
				Data:    payload.Data,
				State:   Unknown,
			}
			err := be.ApplyPolicy(payload.ID, payload.Data)
			if err != nil {
				a.logger.Warn("policy failed to apply", zap.String("id", payload.ID), zap.Error(err))
				pd.State = Failed
				pd.BackendErr = err.Error()
			} else {
				pd.State = Applied
			}
			a.repo.Add(pd)
			err = a.repo.EnsureDataset(payload.ID, payload.DatasetID)
			if err != nil {
				a.logger.Warn("policy failed to ensure dataset id", zap.String("id", payload.ID), zap.String("dataset_id", payload.DatasetID), zap.Error(err))
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
