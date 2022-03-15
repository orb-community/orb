/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

type ConfigRepo interface {
	Exists(sinkID string) bool
	Add(config SinkConfig) error
	Get(sinkID string) (SinkConfig, error)
	Edit(config SinkConfig) error
	GetAll() ([]SinkConfig, error)
}

//type sinkConfigMemRepo struct {
//	logger *zap.Logger
//	db     map[string]SinkConfig
//	mu     sync.Mutex
//}

//func NewMemRepo(logger *zap.Logger) ConfigRepo {
//	repo := &sinkConfigMemRepo{
//		logger: logger,
//		db:     make(map[string]SinkConfig),
//	}
//	return repo
//}
//
//func (s sinkConfigMemRepo) Exists(sinkID string) bool {
//	_, ok := s.db[sinkID]
//	return ok
//}
//
//func (s *sinkConfigMemRepo) Add(config SinkConfig) error {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	s.db[config.SinkID] = config
//	return nil
//}
//
//func (s *sinkConfigMemRepo) Edit(config SinkConfig) error {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	if _, ok := s.db[config.SinkID]; ok {
//		s.db[config.SinkID] = config
//	}
//	return nil
//}
//
//func (s sinkConfigMemRepo) Get(sinkID string) (SinkConfig, error) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	config, ok := s.db[sinkID]
//	if !ok {
//		return SinkConfig{}, errors.New("unknown sink ID")
//	}
//	return config, nil
//}
//
//func (s sinkConfigMemRepo) GetAll() ([]SinkConfig, error) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	configs := []SinkConfig{}
//	for _, v := range s.db {
//		configs = append(configs, v)
//	}
//	return configs, nil
//}
