package migrate

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/orb-community/orb/migrate/postgres"
	"github.com/orb-community/orb/pkg/errors"
	"go.uber.org/zap"
)

type SchemaVersion struct {
	Version int64 `db:"version"`
}

func (s *serviceMigrate) Up() (err error) {
	current, errSchema := s.CurrentSchemaVersion()
	latest := s.LatestSchemaVersion()

	if errSchema != nil {
		return errSchema
	}

	if current == latest {
		s.logger.Info(fmt.Sprintf("current on most recent schema version (%d)", current))
		return
	}
	if current > latest {
		s.logger.Warn(fmt.Sprintf("current schema version %d is greater than latest (version %d)", current, latest))
		return
	}

	s.logger.Info(fmt.Sprintf("applying last migration version %d", latest))
	err = s.migrations[0].Up()
	if err != nil {
		s.logger.Error(fmt.Sprintf("error on migration"), zap.Error(err))
		s.logger.Info(fmt.Sprintf("rolling back migration to %d", current))
		errMigration := s.migrations[0].Down()
		if errMigration != nil {
			s.logger.Error(fmt.Sprintf("error on migration down to %d", current), zap.Error(errMigration))
			return errMigration
		}
		return err
	}

	if errSchema = s.SetSchemaVersion(latest); errSchema != nil {
		return errSchema
	}
	return
}

func (s *serviceMigrate) Down() (err error) {
	_, errSchema := s.CurrentSchemaVersion()
	latest := s.LatestSchemaVersion()
	if errSchema != nil {
		return errSchema
	}

	s.logger.Info(fmt.Sprintf("applying last migration version %d", latest))
	err = s.migrations[0].Down()
	if err != nil {
		s.logger.Error(fmt.Sprintf("error on migration"), zap.Error(err))
		return err
	}

	if errSchema = s.SetSchemaVersion(latest - 1); errSchema != nil {
		return errSchema
	}
	return
}

func (s *serviceMigrate) Drop() error {
	return errors.New("not implemented")
}

func (s *serviceMigrate) SetSchemaVersion(version int64) error {
	return s.doOnTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("UPDATE schema_version SET version = $1", version)
		s.logger.Info(fmt.Sprintf("updated schema version to %d", version))
		return err
	})
}

func (s *serviceMigrate) CurrentSchemaVersion() (int64, error) {
	var schemaVersion int64 = 0
	var m []SchemaVersion

	err := s.doOnTx(func(tx *sqlx.Tx) (err error) {
		err = tx.Select(&m, "SELECT version FROM schema_version")
		if err != nil {
			return
		}
		schemaVersion = m[0].Version
		s.logger.Info(fmt.Sprintf("current schema version %d", schemaVersion))
		return
	})

	return schemaVersion, err
}

func (s *serviceMigrate) LatestSchemaVersion() int64 {
	return 3
}

func (s *serviceMigrate) doOnTx(f func(tx *sqlx.Tx) error) error {
	tx, err := s.dbs[postgres.DbKeto].BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	if err := f(tx); err != nil {
		s.logger.Error("error on tx", zap.Error(err))
		if err = tx.Rollback(); err != nil {
			s.logger.Error("error on tx rollback", zap.Error(err))
		}
	} else {
		s.logger.Debug("tx ok")
		if err = tx.Commit(); err != nil {
			s.logger.Error("error on tx commit", zap.Error(err))
		}
	}
	return err
}
