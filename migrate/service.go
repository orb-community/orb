package migrate

import (
	"github.com/ns1labs/orb/migrate/migration"
)

var _ Service = (*serviceMigrate)(nil)

type serviceMigrate struct {
	migrations []migration.Plan
}

func (sm *serviceMigrate) AddMigration(plan migration.Plan) {
	sm.migrations = append(sm.migrations, plan)
}

func New(plans ...migration.Plan) Service {
	return &serviceMigrate{
		migrations: plans,
	}
}
