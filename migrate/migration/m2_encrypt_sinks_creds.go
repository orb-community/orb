package migration

import (
	"context"
	"github.com/ns1labs/orb/pkg/config"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/backend"
	"github.com/ns1labs/orb/sinks/postgres"
	"go.uber.org/zap"
)

type M2SinksCredentials struct {
	logger  *zap.Logger
	dbSinks postgres.Database
	pwdSvc  sinks.PasswordService
}

type querySink struct {
	Id       string
	Metadata types.Metadata
}

func (m M2SinksCredentials) Up() (err error) {
	ctx := context.Background()
	q := "SELECT Id, Metadata FROM sinks"
	params := map[string]interface{}{}
	rows, err := m.dbSinks.NamedQueryContext(ctx, q, params)
	if err != nil {
		return
	}
	for rows.Next() {
		qSink := querySink{}
		if err2 := rows.StructScan(&qSink); err2 != nil {
			return err2
		}
		sink := sinks.Sink{
			ID:     qSink.Id,
			Config: qSink.Metadata,
		}
		sink, err = m.encryptMetadata(sink)
		if err != nil {
			m.logger.Error("failed to encrypt data for id", zap.String("id", qSink.Id), zap.Error(err))
			return err
		}
		params := map[string]interface{}{
			"id":       sink.ID,
			"metadata": db.Metadata(sink.Config),
		}
		updateQuery := "UPDATE sinks SET metadata = :metadata WHERE id = :id"
		_, err := m.dbSinks.NamedQueryContext(ctx, updateQuery, params)
		if err != nil {
			m.logger.Error("failed to update data for id", zap.String("id", qSink.Id), zap.Error(err))
			return err
		}
	}

	return
}

func (m M2SinksCredentials) Down() (err error) {
	ctx := context.Background()
	q := "SELECT id, metadata FROM sinks"
	var querySinks []querySink
	params := map[string]interface{}{}
	rows, err := m.dbSinks.NamedQueryContext(ctx, q, params)
	if err != nil {
		return
	}
	err = rows.StructScan(querySinks)
	if err != nil {
		return
	}
	for rows.Next() {
		qSink := querySink{}
		if err2 := rows.StructScan(&qSink); err2 != nil {
			return err2
		}
		sink := sinks.Sink{
			ID:     qSink.Id,
			Config: qSink.Metadata,
		}
		sink, err = m.decryptMetadata(sink)
		if err != nil {
			return
		}
		params := map[string]interface{}{
			"id":       sink.ID,
			"metadata": db.Metadata(sink.Config),
		}
		updateQuery := "UPDATE sinks SET metadata = :metadata WHERE id = :id"
		_, err := m.dbSinks.NamedQueryContext(ctx, updateQuery, params)
		if err != nil {
			m.logger.Error("failed to update data for id", zap.String("id", qSink.Id), zap.Error(err))
			return err
		}
	}

	return
}

func NewM2SinksCredentials(log *zap.Logger, dbSinks postgres.Database, config config.EncryptionKey) Plan {
	pwdSvc := sinks.NewPasswordService(log, config.Key)
	return &M2SinksCredentials{log, dbSinks, pwdSvc}
}

func (m M2SinksCredentials) encryptMetadata(sink sinks.Sink) (sinks.Sink, error) {
	var err error
	sink.Config.FilterMap(func(key string) bool {
		return key == backend.ConfigFeatureTypePassword
	}, func(key string, value interface{}) (string, interface{}) {
		newValue, err2 := m.pwdSvc.EncodePassword(value.(string))
		if err2 != nil {
			err = err2
			return key, value
		}
		return key, newValue
	})
	return sink, err
}

func (m M2SinksCredentials) decryptMetadata(sink sinks.Sink) (sinks.Sink, error) {
	var err error
	sink.Config.FilterMap(func(key string) bool {
		return key == backend.ConfigFeatureTypePassword
	}, func(key string, value interface{}) (string, interface{}) {
		newValue, err2 := m.pwdSvc.GetPassword(value.(string))
		if err2 != nil {
			err = err2
			return key, value
		}
		return key, newValue
	})
	return sink, err
}
