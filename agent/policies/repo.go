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
	Get(policyID string) (PolicyData, error)
	Remove(policyID string) error
	Add(data PolicyData) error
	Edit(data PolicyData) error
	GetAll() ([]PolicyData, error)
	EnsureDataset(policyID string, datasetID string) error
}

type policyMemRepo struct {
	logger *zap.Logger

	db map[string]PolicyData
}

func NewMemRepo(logger *zap.Logger) (PolicyRepo, error) {
	r := &policyMemRepo{
		logger: logger,
		db:     make(map[string]PolicyData),
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

func (p policyMemRepo) Get(policyID string) (PolicyData, error) {
	policy, ok := p.db[policyID]
	if !ok {
		return PolicyData{}, errors.New("unknown policy ID")
	}
	return policy, nil
}

func (p policyMemRepo) Remove(policyID string) error {
	delete(p.db, policyID)
	return nil
}

func (p policyMemRepo) Add(data PolicyData) error {
	p.db[data.ID] = data
	return nil
}

func (p policyMemRepo) GetAll() (ret []PolicyData, err error) {
	ret = make([]PolicyData, len(p.db))
	i := 0
	for _, v := range p.db {
		ret[i] = v
		i++
	}
	err = nil
	return ret, err
}

func (p policyMemRepo) Edit(data PolicyData) error {
	p.db[data.ID] = data
	return nil
}

var _ PolicyRepo = (*policyMemRepo)(nil)
