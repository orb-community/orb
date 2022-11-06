package migrate

import (
	"github.com/etaques/orb/migrate/migration"
	"github.com/etaques/orb/migrate/postgres"
	"go.uber.org/zap"
)

var _ Service = (*serviceMigrate)(nil)

type serviceMigrate struct {
	logger     *zap.Logger
	dbs        map[string]postgres.Database
	migrations []migration.Plan
}

func (sm *serviceMigrate) AddMigration(plan migration.Plan) {
	sm.migrations = append(sm.migrations, plan)
}

func New(logger *zap.Logger, dbs map[string]postgres.Database, plans ...migration.Plan) Service {
	return &serviceMigrate{
		logger:     logger,
		dbs:        dbs,
		migrations: plans,
	}
}
