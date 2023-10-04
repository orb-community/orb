package deployment

import (
	"encoding/json"
	"time"

	"github.com/orb-community/orb/pkg/types"
)

type Deployment struct {
	Id                      string     `db:"id" json:"id,omitempty"`
	OwnerID                 string     `db:"owner_id" json:"ownerID,omitempty"`
	SinkID                  string     `db:"sink_id" json:"sinkID,omitempty"`
	Backend                 string     `db:"backend" json:"backend,omitempty"`
	Config                  []byte     `db:"config" json:"config,omitempty"`
	LastStatus              string     `db:"last_status" json:"lastStatus,omitempty"`
	LastStatusUpdate        *time.Time `db:"last_status_update" json:"lastStatusUpdate"`
	LastErrorMessage        string     `db:"last_error_message" json:"lastErrorMessage,omitempty"`
	LastErrorTime           *time.Time `db:"last_error_time" json:"lastErrorTime"`
	CollectorName           string     `db:"collector_name" json:"collectorName,omitempty"`
	LastCollectorDeployTime *time.Time `db:"last_collector_deploy_time" json:"lastCollectorDeployTime"`
	LastCollectorStopTime   *time.Time `db:"last_collector_stop_time" json:"lastCollectorStopTime"`
}

func NewDeployment(ownerID string, sinkID string, config types.Metadata, backend string) Deployment {
	now := time.Now()
	deploymentName := "otel-" + sinkID
	configAsByte := toByte(config)
	return Deployment{
		OwnerID:          ownerID,
		SinkID:           sinkID,
		Backend:          backend,
		Config:           configAsByte,
		LastStatus:       "unknown",
		LastStatusUpdate: &now,
		CollectorName:    deploymentName,
	}
}

func (d *Deployment) Merge(other Deployment) error {
	if other.Id != "" {
		d.Id = other.Id
	}
	if other.LastErrorMessage != d.LastErrorMessage {
		d.LastErrorMessage = other.LastErrorMessage
		d.LastErrorTime = other.LastErrorTime
	}
	if other.CollectorName != "" {
		d.CollectorName = other.CollectorName
		d.LastCollectorDeployTime = other.LastCollectorDeployTime
		d.LastCollectorStopTime = other.LastCollectorStopTime
	}
	if other.LastStatus != d.LastStatus {
		d.LastStatus = other.LastStatus
		d.LastStatusUpdate = other.LastStatusUpdate
	}
	return nil
}

func (d *Deployment) GetConfig() types.Metadata {
	var config types.Metadata
	err := json.Unmarshal(d.Config, &config)
	if err != nil {
		return nil
	}
	return config
}

func (d *Deployment) SetConfig(config types.Metadata) error {
	configAsByte := toByte(config)
	d.Config = configAsByte
	return nil
}

func toByte(config types.Metadata) []byte {
	configAsByte, err := json.Marshal(config)
	if err != nil {
		return nil
	}
	return configAsByte
}
