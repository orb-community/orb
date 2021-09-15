/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PolicyRepo interface {
	Exists(policyID string) bool
	Get(policyID string) (policyData, error)
	Remove(policyID string) error
	Add(data policyData) error
	GetAll() ([]policyData, error)
	EnsureDataset(policyID string, datasetID string) error
}

type policyMemRepo struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func (p policyMemRepo) EnsureDataset(policyID string, datasetID string) error {
	panic("implement me")
}

func (p policyMemRepo) Exists(policyID string) bool {
	panic("implement me")
}

func (p policyMemRepo) Get(policyID string) (policyData, error) {
	panic("implement me")
}

func (p policyMemRepo) Remove(policyID string) error {
	panic("implement me")
}

func (p policyMemRepo) Add(data policyData) error {
	panic("implement me")
}

func (p policyMemRepo) GetAll() ([]policyData, error) {
	panic("implement me")
}

var _ PolicyRepo = (*policyMemRepo)(nil)
