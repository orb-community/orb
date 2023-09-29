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

func TestEventService_HandleSinkActivity(t *testing.T) {
	type args struct {
		event redis.SinkerUpdateEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "activity on a sink that does not exist",
			args: args{
				event: redis.SinkerUpdateEvent{
					OwnerID:   "owner1",
					SinkID:    "sink1",
					State:     "active",
					Size:      "22",
					Timestamp: time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name: "activity success",
			args: args{
				event: redis.SinkerUpdateEvent{
					OwnerID:   "owner2",
					SinkID:    "sink22",
					State:     "active",
					Size:      "22",
					Timestamp: time.Now(),
				},
			}, wantErr: false,
		},
	}
	logger := zap.NewNop()
	deploymentService := deployment.NewDeploymentService(logger, NewFakeRepository(logger), "kafka:9092",
		"MY_SECRET", NewTestProducer(logger), NewTestKubeCtr(logger))
	d := NewEventService(logger, deploymentService, nil)
	err := d.HandleSinkCreate(context.Background(), redis.SinksUpdateEvent{
		SinkID:  "sink22",
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
			if err := d.HandleSinkActivity(ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("HandleSinkActivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_HandleSinkIdle(t *testing.T) {
	type args struct {
		event redis.SinkerUpdateEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sink idle on a sink that does not exist",
			args: args{
				event: redis.SinkerUpdateEvent{
					OwnerID:   "owner1",
					SinkID:    "sink1",
					State:     "idle",
					Size:      "22",
					Timestamp: time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name: "sink idle success",
			args: args{
				event: redis.SinkerUpdateEvent{
					OwnerID:   "owner2",
					SinkID:    "sink222",
					State:     "idle",
					Size:      "22",
					Timestamp: time.Now(),
				},
			}, wantErr: false,
		},
	}
	logger := zap.NewNop()
	deploymentService := deployment.NewDeploymentService(logger, NewFakeRepository(logger), "kafka:9092", "MY_SECRET", NewTestProducer(logger),
		NewTestKubeCtr(logger))
	d := NewEventService(logger, deploymentService, NewTestKubeCtr(logger))
	err := d.HandleSinkCreate(context.Background(), redis.SinksUpdateEvent{
		SinkID:  "sink222",
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
			if err := d.HandleSinkIdle(ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("HandleSinkIdle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
