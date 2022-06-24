package migration

import (
	"context"
	"github.com/ns1labs/orb/pkg/config"
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

type querySinks struct {
	id     string
	config types.Metadata
}

func (m M2SinksCredentials) Up() (err error) {
	ctx := context.Background()
	q := "SELECT id, metadata FROM sinks"
	var querySinks []querySinks
	err = m.dbSinks.QueryRowxContext(ctx, q).StructScan(&querySinks)
	if err != nil {
		return
	}
	for _, qSink := range querySinks {
		sink := sinks.Sink{
			ID:     qSink.id,
			Config: qSink.config,
		}
		sink, err = m.encryptMetadata(sink)
		if err != nil {
			m.logger.Error("failed to encrypt data for id", zap.String("id", qSink.id), zap.Error(err))
			return
		}
		updateQuery := "UPDATE sinks SET metadata VALUES(:metadata)"
		m.dbSinks.NamedQueryContext(ctx, updateQuery, sink.Config)
	}

	return
}

func (m M2SinksCredentials) Down() (err error) {
	ctx := context.Background()
	q := "SELECT id, metadata FROM sinks"
	var querySinks []querySinks
	err = m.dbSinks.QueryRowxContext(ctx, q).StructScan(&querySinks)
	if err != nil {
		return
	}
	for _, qSink := range querySinks {
		sink := sinks.Sink{
			ID:     qSink.id,
			Config: qSink.config,
		}
		sink = m.decryptMetadata(sink)
		params := map[string]interface{}{
			"id":       sink.ID,
			"metadata": sink.Config,
		}
		updateQuery := "UPDATE sinks SET metadata VALUES(:metadata) WHERE id = :id"
		m.dbSinks.NamedExecContext(ctx, updateQuery, params)
	}

	return
}

func NewM2SinksCredentials(log *zap.Logger, dbSinks postgres.Database, config config.EncryptionKey) Plan {
	pwdSvc := sinks.NewPasswordService(log, config.Key)
	return &M2SinksCredentials{log, dbSinks, pwdSvc}
}

func (svc M2SinksCredentials) encryptMetadata(sink sinks.Sink) (sinks.Sink, error) {
	var err error
	sink.Config.FilterMap(func(key string) bool {
		return key == backend.ConfigFeatureTypePassword
	}, func(key string, value interface{}) (string, interface{}) {
		newValue, err2 := svc.pwdSvc.EncodePassword(value.(string))
		if err2 != nil {
			err = err2
			return key, value
		}
		return key, newValue
	})
	return sink, err
}

func (svc M2SinksCredentials) decryptMetadata(sink sinks.Sink) sinks.Sink {
	sink.Config.FilterMap(func(key string) bool {
		return key == backend.ConfigFeatureTypePassword
	}, func(key string, value interface{}) (string, interface{}) {
		newValue := svc.pwdSvc.GetPassword(value.(string))
		return key, newValue
	})
	return sink
}
