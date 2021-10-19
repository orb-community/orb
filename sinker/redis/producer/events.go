package producer

import (
	"time"
)

const (
	SinkerPrefix = "sinker."
	ChangeState  = SinkerPrefix + "change_state"
)

type event interface {
	Encode() map[string]interface{}
}

var (
	_ event = (*ChangeSinkerStateEvent)(nil)
)

type ChangeSinkerStateEvent struct {
	SinkID    string
	Owner     string
	State     string
	Msg       string
	Timestamp time.Time
}

func (cse ChangeSinkerStateEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"sink_id":   cse.SinkID,
		"owner":     cse.Owner,
		"state":     cse.State,
		"msg":       cse.Msg,
		"timestamp": cse.Timestamp.Unix(),
		"operation": ChangeState,
	}
}
