package service

import (
	"context"
	"github.com/orb-community/orb/maestro/deployment"
	"github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/pkg/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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
					SinkID: "crt-sink1",
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
		{
			name: "create event without config",
			args: args{
				event: redis.SinksUpdateEvent{
					SinkID:    "crt-sink1",
					Owner:     "owner1",
					Config:    nil,
					Backend:   "prometheus",
					Timestamp: time.Now(),
				},
			},
			wantErr: true,
		},
	}
	logger := zap.NewNop()
	deploymentService := deployment.NewDeploymentService(logger, NewFakeRepository(logger), "kafka:9092", "MY_SECRET", NewTestProducer(logger), nil)
	d := NewEventService(logger, deploymentService, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "test", tt.name)
			if err := d.HandleSinkCreate(ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("HandleSinkCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_HandleSinkUpdate(t *testing.T) {
	type args struct {
		event redis.SinksUpdateEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update event when there is none in db",
			args: args{
				event: redis.SinksUpdateEvent{
					SinkID: "upd-sink1",
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
		{
			name: "update event success",
			args: args{
				event: redis.SinksUpdateEvent{
					SinkID:  "upd-sink1",
					Owner:   "owner1",
					Backend: "prometheus",
					Config: types.Metadata{
						"exporter": types.Metadata{
							"remote_host": "https://acme.com/prom/push",
						},
						"authentication": types.Metadata{
							"type":     "basicauth",
							"username": "prom-user-2",
							"password": "dbpass-2",
						},
					},
					Timestamp: time.Now(),
				},
			},
			wantErr: false,
		},
	}
	logger := zap.NewNop()
	deploymentService := deployment.NewDeploymentService(logger, NewFakeRepository(logger), "kafka:9092", "MY_SECRET", NewTestProducer(logger),
		NewTestKubeCtr(logger))
	d := NewEventService(logger, deploymentService, NewTestKubeCtr(logger))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "test", tt.name)
			if err := d.HandleSinkUpdate(ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("HandleSinkUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_HandleSinkDelete(t *testing.T) {
	type args struct {
		event redis.SinksUpdateEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete event when there is none in db",
			args: args{
				event: redis.SinksUpdateEvent{
					SinkID:  "sink1",
					Owner:   "owner1",
					Backend: "prometheus",
					Config: types.Metadata{
						"exporter": types.Metadata{
							"remote_host": "https://acme.com/prom/push",
						},
						"authentication": types.Metadata{
							"type":     "basicauth",
							"username": "prom-user-2",
							"password": "dbpass-2",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "delete event success",
			args: args{
				event: redis.SinksUpdateEvent{
					SinkID:  "sink2-1",
					Owner:   "owner2",
					Backend: "prometheus",
					Config: types.Metadata{
						"exporter": types.Metadata{
							"remote_host": "https://acme.com/prom/push",
						},
						"authentication": types.Metadata{
							"type":     "basicauth",
							"username": "prom-user-2",
							"password": "dbpass-2",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	logger := zap.NewNop()
	deploymentService := deployment.NewDeploymentService(logger, NewFakeRepository(logger), "kafka:9092", "MY_SECRET", NewTestProducer(logger), nil)
	d := NewEventService(logger, deploymentService, nil)
	err := d.HandleSinkCreate(context.Background(), redis.SinksUpdateEvent{
		SinkID:  "sink2-1",
		Owner:   "owner2",
		Backend: "prometheus",
		Config: types.Metadata{
			"exporter": types.Metadata{
				"remote_host": "https://acme.com/prom/push",
			},
			"authentication": types.Metadata{
				"type":     "basicauth",
				"username": "prom-user-2",
				"password": "dbpass-2",
			},
		},
	})
	require.NoError(t, err, "should not error")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "test", tt.name)
			if err := d.HandleSinkDelete(ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("HandleSinkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
