package migrate

import (
	"github.com/ns1labs/orb/migrate/migration"
	"github.com/ns1labs/orb/migrate/postgres"
	"go.uber.org/zap"
)

var _ Service = (*serviceMigrate)(nil)

type serviceMigrate struct {
	logger     *zap.Logger
	dbs        map[string]postgres.Database
	migrations []migration.Plan
}

func New(logger *zap.Logger, dbs map[string]postgres.Database) Service {
	return &serviceMigrate{
		logger: logger,
		dbs:    dbs,
		migrations: []migration.Plan{
			migration.NewM0KetoPolicies(logger, dbs),
		},
	}
}
