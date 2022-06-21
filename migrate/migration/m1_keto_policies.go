package migration

import (
	"github.com/ns1labs/orb/migrate/postgres"
	"go.uber.org/zap"
)

type M1KetoPolicies struct {
	log *zap.Logger
	dbs map[string]postgres.Database
}

func NewM0KetoPolicies(log *zap.Logger, dbs map[string]postgres.Database) Plan {
	return &M1KetoPolicies{log, dbs}
}

func (m *M1KetoPolicies) Up() error {
	m.log.Info("do nothing for now")
	return nil
}

func (m *M1KetoPolicies) Down() error {
	m.log.Info("do nothing for now")
	return nil
}
