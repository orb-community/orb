package deployment

import (
	"github.com/orb-community/orb/pkg/types"
	"time"
)

type Deployment struct {
	Id                      string         `db:"id" json:"id,omitempty"`
	OwnerID                 string         `db:"owner_id" json:"ownerID,omitempty"`
	SinkID                  string         `db:"sink_id" json:"sinkID,omitempty"`
	Config                  types.Metadata `db:"config" json:"config,omitempty"`
	LastStatus              string         `db:"last_status" json:"lastStatus,omitempty"`
	LastStatusUpdate        *time.Time     `db:"last_status_update" json:"lastStatusUpdate"`
	LastErrorMessage        string         `db:"last_error_message" json:"lastErrorMessage,omitempty"`
	LastErrorTime           *time.Time     `db:"last_error_time" json:"lastErrorTime"`
	CollectorName           string         `db:"collector_name" json:"collectorName,omitempty"`
	LastCollectorDeployTime *time.Time     `db:"last_collector_deploy_time" json:"lastCollectorDeployTime"`
	LastCollectorStopTime   *time.Time     `db:"last_collector_stop_time" json:"lastCollectorStopTime"`
}

func NewDeployment(ownerID string, sinkID string, config types.Metadata) Deployment {
	now := time.Now()
	return Deployment{
		OwnerID:          ownerID,
		SinkID:           sinkID,
		Config:           config,
		LastStatus:       "pending",
		LastStatusUpdate: &now,
	}
}

func (d *Deployment) Merge(other Deployment) error {
	if other.Id != "" {
		d.Id = other.Id
	}
	if other.LastErrorMessage != "" {
		d.LastErrorMessage = other.LastErrorMessage
		d.LastErrorTime = other.LastErrorTime
	}
	if other.CollectorName != "" {
		d.CollectorName = other.CollectorName
		d.LastCollectorDeployTime = other.LastCollectorDeployTime
		d.LastCollectorStopTime = other.LastCollectorStopTime
	}
	if other.LastStatus != "" {
		d.LastStatus = other.LastStatus
		d.LastStatusUpdate = other.LastStatusUpdate
	}
	return nil
}
