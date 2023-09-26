package redis

import (
	"github.com/orb-community/orb/pkg/types"
	"time"
)

const (
	SinkerPrefix        = "sinker."
	SinkerUpdate        = SinkerPrefix + "update"
	SinksActivityStream = "orb.sink_activity"
	SinksIdleStream     = "orb.sink_idle"
	GroupMaestro        = "orb.maestro"
	Exists              = "BUSYGROUP Consumer Group name already exists"
)

type SinksUpdateEvent struct {
	SinkID    string
	Owner     string
	Config    types.Metadata
	Backend   string
	Timestamp time.Time
}

type SinkerUpdateEvent struct {
	OwnerID   string
	SinkID    string
	State     string
	Size      string
	Timestamp time.Time
}

func (sue SinksUpdateEvent) Decode(values map[string]interface{}) {
	sue.SinkID = values["sink_id"].(string)
	sue.Owner = values["owner"].(string)
	sue.Config = types.FromMap(values["config"].(map[string]interface{}))
	sue.Backend = values["backend"].(string)
	sue.Timestamp = values["timestamp"].(time.Time)
}

func (cse SinkerUpdateEvent) Decode(values map[string]interface{}) {
	cse.OwnerID = values["owner_id"].(string)
	cse.SinkID = values["sink_id"].(string)
	cse.State = values["state"].(string)
	cse.Size = values["size"].(string)
	cse.Timestamp = values["timestamp"].(time.Time)
}

func (cse SinkerUpdateEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"sink_id":   cse.SinkID,
		"owner":     cse.OwnerID,
		"state":     cse.State,
		"timestamp": cse.Timestamp.Unix(),
		"operation": SinkerUpdate,
	}
}

type DeploymentEvent struct {
	SinkID         string
	DeploymentYaml string
}
