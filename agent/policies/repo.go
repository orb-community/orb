/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"errors"
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

	db map[string]policyData
}

func NewMemRepo(logger *zap.Logger) (PolicyRepo, error) {
	r := &policyMemRepo{
		logger: logger,
		db:     make(map[string]policyData),
	}
	return r, nil
}

func (p policyMemRepo) EnsureDataset(policyID string, datasetID string) error {
	policy, ok := p.db[policyID]
	if !ok {
		return errors.New("unknown policy ID")
	}
	policy.Datasets[datasetID] = true
	return nil
}

func (p policyMemRepo) Exists(policyID string) bool {
	_, ok := p.db[policyID]
	return ok
}

func (p policyMemRepo) Get(policyID string) (policyData, error) {
	policy, ok := p.db[policyID]
	if !ok {
		return policyData{}, errors.New("unknown policy ID")
	}
	return policy, nil
}

func (p policyMemRepo) Remove(policyID string) error {
	delete(p.db, policyID)
	return nil
}

func (p policyMemRepo) Add(data policyData) error {
	p.db[data.ID] = data
	return nil
}

func (p policyMemRepo) GetAll() (ret []policyData, err error) {
	ret = make([]policyData, len(p.db))
	for _, v := range p.db {
		ret = append(ret, v)
	}
	err = nil
	return ret, err
}

var _ PolicyRepo = (*policyMemRepo)(nil)
