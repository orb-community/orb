package redis

import (
	"encoding/json"
	"github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/pkg/types"
	"time"
)

const (
	SinkPrefix   = "sinks."
	SinkCreate   = SinkPrefix + "create"
	SinkDelete   = SinkPrefix + "remove"
	SinkUpdate   = SinkPrefix + "update"
	StreamSinks  = "orb.sinks"
	GroupMaestro = "orb.maestro"
	Exists       = "BUSYGROUP Consumer Group name already exists"
)

type StateUpdateEvent struct {
	OwnerID   string
	SinkID    string
	State     string
	Msg       string
	Timestamp time.Time
}

func DecodeSinksEvent(event map[string]interface{}, operation string) (redis.SinksUpdateEvent, error) {
	val := redis.SinksUpdateEvent{
		SinkID:    read(event, "sink_id", ""),
		Owner:     read(event, "owner", ""),
		Backend:   read(event, "backend", ""),
		Config:    readMetadata(event, "config"),
		Timestamp: time.Now(),
	}
	if operation != SinkDelete {
		var metadata types.Metadata
		if err := json.Unmarshal([]byte(read(event, "config", "")), &metadata); err != nil {
			return redis.SinksUpdateEvent{}, err
		}
		val.Config = metadata
		return val, nil
	}

	return val, nil
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}

func readMetadata(event map[string]interface{}, key string) types.Metadata {
	val, ok := event[key].(types.Metadata)
	if !ok {
		return types.Metadata{}
	}

	return val
}
