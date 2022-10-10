package producer

import (
	"time"
)

const (
	SinkerPrefix = "sinker."
	SinkerUpdate = SinkerPrefix + "update"
)

type event interface {
	Encode() map[string]interface{}
}

var (
	_ event = (*SinkerUpdateEvent)(nil)
)

type SinkerUpdateEvent struct {
	SinkID    string
	Owner     string
	State     string
	Msg       string
	Timestamp time.Time
}

type SinkerOtelConfigEvent struct {
	SinkId     string
	Owner      string
	State      string
	ConfigYaml string
	Timestamp  time.Time
}

func (e SinkerOtelConfigEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"sink_id":   e.SinkId,
		"owner":     e.Owner,
		"state":     e.State,
		"config":    e.ConfigYaml,
		"timestamp": e.Timestamp.Unix(),
	}
}

func (cse SinkerUpdateEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"sink_id":   cse.SinkID,
		"owner":     cse.Owner,
		"state":     cse.State,
		"msg":       cse.Msg,
		"timestamp": cse.Timestamp.Unix(),
		"operation": SinkerUpdate,
	}
}
