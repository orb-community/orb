package redis

import (
	"github.com/orb-community/orb/pkg/types"
	"time"
)

const (
	SinkerPrefix = "sinker."
	SinkerUpdate = SinkerPrefix + "update"
	SinkPrefix   = "sinks."
	SinkCreate   = SinkPrefix + "create"
	SinkDelete   = SinkPrefix + "remove"
	SinkUpdate   = SinkPrefix + "update"
	StreamSinks  = "orb.sinks"
	GroupMaestro = "orb.maestro"
	Exists       = "BUSYGROUP Consumer Group name already exists"
)

type SinksUpdateEvent struct {
	SinkID    string
	Owner     string
	Config    types.Metadata
	Backend   string
	Timestamp time.Time
}

type SinkerUpdateEvent struct {
	SinkID    string
	Owner     string
	State     string
	Timestamp time.Time
}

func (sue SinksUpdateEvent) Decode(values map[string]interface{}) {
	sue.SinkID = values["sink_id"].(string)
	sue.Owner = values["owner"].(string)
	sue.Config = values["config"].(types.Metadata)
	sue.Backend = values["backend"].(string)
	sue.Timestamp = time.Unix(values["timestamp"].(int64), 0)
}

func (cse SinkerUpdateEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"sink_id":   cse.SinkID,
		"owner":     cse.Owner,
		"state":     cse.State,
		"timestamp": cse.Timestamp.Unix(),
		"operation": SinkerUpdate,
	}
}

type DeploymentEvent struct {
	SinkID         string
	DeploymentYaml string
}
