package migrate

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/migrate/postgres"
	"github.com/ns1labs/orb/pkg/errors"
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

	errorIndex := int64(0)
	index := current
	for index < latest {
		s.logger.Info(fmt.Sprintf("applying migration %d of %d", index+1, latest))
		err = s.migrations[index].Up()
		if err != nil {
			s.logger.Error(fmt.Sprintf("error on migration up %d of %d", index+1, latest), zap.Error(err))
			errorIndex = index
			break
		}
		index++
	}

	for errorIndex > current {
		s.logger.Info(fmt.Sprintf("rolling back migration %d of %d", errorIndex+1, latest))
		errMigration := s.migrations[errorIndex].Down()
		if errMigration != nil {
			s.logger.Error(fmt.Sprintf("error on migration down %d of %d", errorIndex+1, latest), zap.Error(errMigration))
			break
		}
		errorIndex--
	}

	if errSchema = s.SetSchemaVersion(index); errSchema != nil {
		return errSchema
	}
	return
}

func (s *serviceMigrate) Down() error {
	return errors.New("not implemented")
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
	return int64(len(s.migrations))
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
