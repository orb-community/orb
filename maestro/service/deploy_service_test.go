package service

import (
	"github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/pkg/types"
	"testing"
	"time"
)

func Test_eventService_HandleSinkCreate(t *testing.T) {

	type args struct {
		event redis.SinksUpdateEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create event",
			args: args{
				event: redis.SinksUpdateEvent{
					SinkID: "sink1",
					Owner:  "owner1",
					Config: types.Metadata{
						"exporter": types.Metadata{
							"remote_host": "https://acme.com/prom/push",
						},
						"authentication": types.Metadata{
							"type":     "basicauth",
							"username": "prom-user",
							"password": "dbpass",
						},
					},
					Backend:   "prometheus",
					Timestamp: time.Now(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//logger := zap.NewNop()
			//deploymentService := deployment.NewDeploymentService(logger,
			//d := NewEventService(logger, )
			//if err := d.HandleSinkCreate(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
			//	t.Errorf("HandleSinkCreate() error = %v, wantErr %v", err, tt.wantErr)
			//}
		})
	}
}
