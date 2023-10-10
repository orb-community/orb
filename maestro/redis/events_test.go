package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSinkerUpdateEvent_Decode(t *testing.T) {
	type fields struct {
		OwnerID string
		SinkID  string
		State   string
		Size    string
	}
	type args struct {
		values map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test_decode_allfields", fields: fields{
			OwnerID: "owner-1",
			SinkID:  "sink-1",
			State:   "active",
			Size:    "111",
		}, args: args{
			values: map[string]interface{}{
				"owner_id":  "owner-1",
				"sink_id":   "sink-1",
				"state":     "active",
				"size":      "111",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cse := SinkerUpdateEvent{}
			cse.Decode(tt.args.values)
			assert.Equal(t, tt.fields.OwnerID, cse.OwnerID)
			assert.Equal(t, tt.fields.SinkID, cse.SinkID)
			assert.Equal(t, tt.fields.State, cse.State)
			assert.Equal(t, tt.fields.Size, cse.Size)
		})
	}
}
