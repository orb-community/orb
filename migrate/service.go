package migrate

import (
	"github.com/ns1labs/orb/migrate/migration"
	"github.com/ns1labs/orb/migrate/postgres"
	"github.com/ns1labs/orb/pkg/config"
	"go.uber.org/zap"
)

var _ Service = (*serviceMigrate)(nil)

type serviceMigrate struct {
	logger     *zap.Logger
	dbs        map[string]postgres.Database
	migrations []migration.Plan
}

func New(logger *zap.Logger, dbs map[string]postgres.Database, config config.EncryptionKey) Service {
	return &serviceMigrate{
		logger: logger,
		dbs:    dbs,
		migrations: []migration.Plan{
			migration.NewM1KetoPolicies(logger, dbs),
			migration.NewM2SinksCredentials(logger, dbs["sinks"], config),
		},
	}
}
