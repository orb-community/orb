package migration

import (
	"context"
	"github.com/ns1labs/orb/migrate/postgres"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	"go.uber.org/zap"
)

type M3SinksOpenTelemetry struct {
	logger  *zap.Logger
	dbSinks postgres.Database
	pwdSvc  sinks.PasswordService
}

func NewM3SinksOpenTelemetry(log *zap.Logger, dbSinks postgres.Database) Plan {
	return &M3SinksOpenTelemetry{logger: log, dbSinks: dbSinks}
}

func (m M3SinksOpenTelemetry) Up() (err error) {
	ctx := context.Background()
	q := "SELECT Id, Metadata FROM sinks"
	params := map[string]interface{}{}
	rows, err := m.dbSinks.NamedQueryContext(ctx, q, params)
	if err != nil {
		return
	}
	for rows.Next() {
		qSink := querySink{}
		if err = rows.StructScan(&qSink); err != nil {
			return err
		}
		sink := sinks.Sink{
			ID:     qSink.Id,
			Config: qSink.Metadata,
		}
		sink, err = m.addOpenTelemetryFlag(sink)
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
	return nil
}

func (m M3SinksOpenTelemetry) Down() (err error) {
	ctx := context.Background()
	q := "SELECT Id, Metadata FROM sinks"
	params := map[string]interface{}{}
	rows, err := m.dbSinks.NamedQueryContext(ctx, q, params)
	if err != nil {
		return
	}
	for rows.Next() {
		qSink := querySink{}
		if err = rows.StructScan(&qSink); err != nil {
			return err
		}
		sink := sinks.Sink{
			ID:     qSink.Id,
			Config: qSink.Metadata,
		}
		sink, err = m.rollbackOpenTelemetryFlag(sink)
		if err != nil {
			if err.Error() != "skip" {
				m.logger.Error("failed to encrypt data for id", zap.String("id", qSink.Id), zap.Error(err))
				return err
			}
			continue
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
	return nil
}

func (m M3SinksOpenTelemetry) addOpenTelemetryFlag(sink sinks.Sink) (sinks.Sink, error) {
	newMetadata := types.Metadata{
		"opentelemetry": "enabled",
		"migrated":      "m3",
	}
	sink.Config.Merge(newMetadata)
	return sink, nil
}

func (m M3SinksOpenTelemetry) rollbackOpenTelemetryFlag(sink sinks.Sink) (sinks.Sink, error) {
	if _, ok := sink.Config["migrated"]; !ok {
		return sinks.Sink{}, errors.New("skip")
	}
	sink.Config.RemoveKeys([]string{"opentelemetry", "migrated"})
	return sink, nil
}
