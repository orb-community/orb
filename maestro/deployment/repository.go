package deployment

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"go.uber.org/zap"
	"time"
)

type Repository interface {
	FetchAll(ctx context.Context) ([]Deployment, error)
	Add(ctx context.Context, deployment Deployment) (Deployment, error)
	Update(ctx context.Context, deployment Deployment) (Deployment, error)
	Remove(ctx context.Context, ownerId string, sinkId string) error
	FindByOwnerAndSink(ctx context.Context, ownerId string, sinkId string) (Deployment, error)
}

type Deployment struct {
	Id                      string         `db:"id" json:"id,omitempty"`
	OwnerID                 string         `db:"owner_id" json:"ownerID,omitempty"`
	SinkID                  string         `db:"sink_id" json:"sinkID,omitempty"`
	Config                  types.Metadata `db:"config" json:"config,omitempty"`
	LastStatus              string         `db:"last_status" json:"lastStatus,omitempty"`
	LastStatusUpdate        time.Time      `db:"last_status_update" json:"lastStatusUpdate"`
	LastErrorMessage        string         `db:"last_error_message" json:"lastErrorMessage,omitempty"`
	LastErrorTime           time.Time      `db:"last_error_time" json:"lastErrorTime"`
	CollectorName           string         `db:"collector_name" json:"collectorName,omitempty"`
	LastCollectorDeployTime time.Time      `db:"last_collector_deploy_time" json:"lastCollectorDeployTime"`
	LastCollectorStopTime   time.Time      `db:"last_collector_stop_time" json:"lastCollectorStopTime"`
}

var _ Repository = (*repositoryService)(nil)

func NewRepositoryService(db *sqlx.DB, logger *zap.Logger) Repository {
	namedLogger := logger.Named("deployment-repository")
	return &repositoryService{db: db, logger: namedLogger}
}

type repositoryService struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func (r *repositoryService) FetchAll(ctx context.Context) ([]Deployment, error) {
	tx := r.db.MustBeginTx(ctx, nil)
	var deployments []Deployment
	err := tx.SelectContext(ctx, &deployments, "SELECT * FROM deployments", nil)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	r.logger.Info("fetched all deployments", zap.Int("count", len(deployments)))
	return deployments, nil
}

func (r *repositoryService) Add(ctx context.Context, deployment Deployment) (Deployment, error) {
	tx := r.db.MustBeginTx(ctx, nil)
	_, err := tx.NamedExecContext(ctx,
		`INSERT INTO deployments (id, owner_id, sink_id, config, last_status, last_status_update, last_error_message, 
				last_error_time, collector_name, last_collector_deploy_time, last_collector_stop_time) 
				VALUES (:id, :owner_id, :sink_id, :config, :last_status, :last_status_update, :last_error_message, 
				        :last_error_time, :collector_name, :last_collector_deploy_time, :last_collector_stop_time)`,
		deployment)
	if err != nil {
		_ = tx.Rollback()
		return Deployment{}, err
	}
	r.logger.Info("added deployment", zap.String("owner-id", deployment.OwnerID), zap.String("sink-id", deployment.SinkID))
	return deployment, tx.Commit()
}

func (r *repositoryService) Update(ctx context.Context, deployment Deployment) (Deployment, error) {
	tx := r.db.MustBeginTx(ctx, nil)
	_, err := tx.NamedExecContext(ctx,
		`UPDATE deployments 
				SET 
                       owner_id = :owner_id,
                       sink_id = :sink_id,
                       config = :config,
                       last_status = :last_status, 
                       last_status_update = :last_status_update, 
                       last_error_message = :last_error_message,
					   last_error_time = :last_error_time, 
					   collector_name = :collector_name, 
					   last_collector_deploy_time = :last_collector_deploy_time, 
					   last_collector_stop_time = :last_collector_stop_time 
				WHERE id = :id`,
		deployment)
	if err != nil {
		_ = tx.Rollback()
		return Deployment{}, err
	}
	r.logger.Info("update deployment", zap.String("owner-id", deployment.OwnerID), zap.String("sink-id", deployment.SinkID))
	return deployment, tx.Commit()
}

func (r *repositoryService) Remove(ctx context.Context, ownerId string, sinkId string) error {
	tx := r.db.MustBeginTx(ctx, nil)
	tx.MustExecContext(ctx, "DELETE FROM deployments WHERE owner_id = $1 AND sink_id = $2", ownerId, sinkId)
	err := tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func (r *repositoryService) FindByOwnerAndSink(ctx context.Context, ownerId string, sinkId string) (Deployment, error) {
	tx := r.db.MustBeginTx(ctx, nil)
	var rows []Deployment
	err := tx.SelectContext(ctx, &rows, "SELECT * FROM deployments WHERE owner_id = :owner_id AND sink_id = :sink_id",
		map[string]interface{}{"owner_id": ownerId, "sink_id": sinkId})
	if err != nil {
		_ = tx.Rollback()
		return Deployment{}, err
	}
	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return Deployment{}, err
	}
	if len(rows) == 0 {
		return Deployment{}, errors.New("")
	}
	deployment := rows[0]

	return deployment, nil
}
