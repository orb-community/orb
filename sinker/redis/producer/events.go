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
