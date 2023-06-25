package config

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/orb-community/orb/pkg/types"
)

func TestReturnConfigYamlFromSink(t *testing.T) {
	type args struct {
		in0            context.Context
		kafkaUrlConfig string
		sink           SinkData
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "prometheus, basicauth",
			args: args{
				in0:            context.Background(),
				kafkaUrlConfig: "kafka:9092",
				sink: SinkData{
					SinkID:  "sink-id-11",
					OwnerID: "11",
					Backend: "prometheus",
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
					State:           0,
					Msg:             "",
					LastRemoteWrite: time.Time{},
				},
			},
			want:    `---\nreceivers:\n  kafka:\n    brokers:\n    - kafka:9092\n    topic: otlp_metrics-sink-id-11\n    protocol_version: 2.0.0\nextensions:\n  pprof:\n    endpoint: 0.0.0.0:1888\n  basicauth/exporter:\n    client_auth:\n      username: prom-user\n      password: dbpass\nexporters:\n  prometheusremotewrite:\n    endpoint: https://acme.com/prom/push\n    auth:\n      authenticator: basicauth/exporter\nservice:\n  extensions:\n  - pprof\n  - basicauth/exporter\n  pipelines:\n    metrics:\n      receivers:\n      - kafka\n      exporters:\n      - prometheusremotewrite\n`,
			wantErr: false,
		},
		{
			name: "otlp, basicauth",
			args: args{
				in0:            context.Background(),
				kafkaUrlConfig: "kafka:9092",
				sink: SinkData{
					SinkID:  "sink-id-22",
					OwnerID: "22",
					Backend: "otlphttp",
					Config: types.Metadata{
						"exporter": types.Metadata{
							"endpoint": "https://acme.com/otlphttp/push",
						},
						"authentication": types.Metadata{
							"type":     "basicauth",
							"username": "otlp-user",
							"password": "dbpass",
						},
					},
					State:           0,
					Msg:             "",
					LastRemoteWrite: time.Time{},
				},
			},
			want:    `---\nreceivers:\n  kafka:\n    brokers:\n    - kafka:9092\n    topic: otlp_metrics-sink-id-22\n    protocol_version: 2.0.0\nextensions:\n  pprof:\n    endpoint: 0.0.0.0:1888\n  basicauth/exporter:\n    client_auth:\n      username: otlp-user\n      password: dbpass\nexporters:\n  otlphttp:\n    endpoint: https://acme.com/otlphttp/push\n    auth:\n      authenticator: basicauth/exporter\nservice:\n  extensions:\n  - pprof\n  - basicauth/exporter\n  pipelines:\n    metrics:\n      receivers:\n      - kafka\n      exporters:\n      - otlphttp\n`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		got, err := ReturnConfigYamlFromSink(tt.args.in0, tt.args.kafkaUrlConfig, tt.args.sink)
		if (err != nil) != tt.wantErr {
			t.Errorf("ReturnConfigYamlFromSink() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		fmt.Printf("%s\n", got)
	// 		if got != tt.want {
	// 			t.Errorf("ReturnConfigYamlFromSink() got = \n%v\n, want \n%v", got, tt.want)
	// 		}
		})
	}
}
