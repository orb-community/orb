package deployment

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func Test_repositoryService_FindByOwnerAndSink(t *testing.T) {
	now := time.Now()
	deployCreate := &Deployment{
		OwnerID: "owner-1",
		SinkID:  "sink-1",
		Backend: "prometheus",
		Config: []byte(`{
			"authentication": {
				"username": "user",
				"password": "pass"
			},
			"exporter" : {
				"remote_host": "http://localhost:9090"
			}
		}`),
		LastStatus:              "pending",
		LastStatusUpdate:        &now,
		LastErrorMessage:        "",
		LastErrorTime:           &now,
		CollectorName:           "",
		LastCollectorDeployTime: &now,
		LastCollectorStopTime:   &now,
	}
	type args struct {
		ownerId string
		sinkId  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Deployment
		wantErr bool
	}{
		{
			name: "FindByOwnerAndSink_success",
			args: args{
				ownerId: "owner-1",
				sinkId:  "sink-1",
			},
			want:    deployCreate,
			wantErr: false,
		},
	}

	r := &repositoryService{
		logger: zap.NewNop(),
		db:     pg,
	}
	_, err := r.Add(context.Background(), deployCreate)
	if err != nil {
		t.Fatalf("error adding deployment: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "test", tt.name)
			got, err := r.FindByOwnerAndSink(ctx, tt.args.ownerId, tt.args.sinkId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByOwnerAndSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want.SinkID, got.SinkID)
			require.Equal(t, tt.want.OwnerID, got.OwnerID)
			require.Equal(t, tt.want.Backend, got.Backend)
			var gotInterface map[string]interface{}
			err = json.Unmarshal(got.Config, &gotInterface)
			require.NoError(t, err)
			var wantInterface map[string]interface{}
			err = json.Unmarshal(tt.want.Config, &wantInterface)
			require.NoError(t, err)
			require.Equal(t, wantInterface, gotInterface)
		})
	}
}

func Test_repositoryService_AddUpdateRemove(t *testing.T) {
	now := time.Now()
	type args struct {
		create *Deployment
		update *Deployment
	}
	tests := []struct {
		name    string
		args    args
		want    *Deployment
		wantErr bool
	}{
		{
			name: "update_success",
			args: args{
				create: &Deployment{
					OwnerID: "owner-1",
					SinkID:  "sink-1",
					Backend: "prometheus",
					Config: []byte(`{
			"authentication": {
				"username": "user",
				"password": "pass"
			},
			"exporter" : {
				"remote_host": "http://localhost:9090"
			}
		}`),
					LastStatus:              "pending",
					LastStatusUpdate:        &now,
					LastErrorMessage:        "",
					LastErrorTime:           &now,
					CollectorName:           "",
					LastCollectorDeployTime: &now,
					LastCollectorStopTime:   &now,
				},
				update: &Deployment{
					OwnerID: "owner-1",
					SinkID:  "sink-1",
					Backend: "prometheus",
					Config: []byte(`{
			"authentication": {
				"username": "user2",
				"password": "pass2"
			},
			"exporter" : {
				"remote_host": "http://localhost:9090"
			}
		}`),
					LastStatus:              "pending",
					LastStatusUpdate:        &now,
					LastErrorMessage:        "",
					LastErrorTime:           &now,
					CollectorName:           "",
					LastCollectorDeployTime: &now,
					LastCollectorStopTime:   &now,
				},
			},
			want: &Deployment{
				OwnerID: "owner-1",
				SinkID:  "sink-1",
				Backend: "prometheus",
				Config: []byte(`{
			"authentication": {
				"username": "user2",
				"password": "pass2"
			},
			"exporter" : {
				"remote_host": "http://localhost:9090"
			}
		}`),
				LastStatus:              "pending",
				LastStatusUpdate:        &now,
				LastErrorMessage:        "",
				LastErrorTime:           &now,
				CollectorName:           "",
				LastCollectorDeployTime: &now,
				LastCollectorStopTime:   &now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "test", tt.name)
			r := &repositoryService{
				logger: logger,
				db:     pg,
			}
			got, err := r.Add(ctx, tt.args.create)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotEmptyf(t, got.Id, "id should not be empty")
			var gotInterface map[string]interface{}
			var wantInterface map[string]interface{}

			tt.args.update.Id = got.Id

			got, err = r.Update(ctx, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want.SinkID, got.SinkID)
			require.Equal(t, tt.want.OwnerID, got.OwnerID)
			require.Equal(t, tt.want.Backend, got.Backend)
			err = json.Unmarshal(got.Config, &gotInterface)
			require.NoError(t, err)
			err = json.Unmarshal(tt.want.Config, &wantInterface)
			require.NoError(t, err)
			require.Equal(t, wantInterface, gotInterface)

			if err := r.Remove(ctx, tt.want.OwnerID, tt.want.SinkID); (err != nil) != tt.wantErr {
				t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
