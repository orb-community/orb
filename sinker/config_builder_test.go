package sinker

import (
	"context"
	"testing"
)

func TestReturnConfigYamlFromSink(t *testing.T) {
	type args struct {
		in0            context.Context
		kafkaUrlConfig string
		sinkId         string
		sinkUrl        string
		sinkUsername   string
		sinkPassword   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "simple test", args: args{
			in0:            context.Background(),
			kafkaUrlConfig: "kafka:9092",
			sinkId:         "sink-id-222",
			sinkUrl:        "https://mysinkurl:9922",
			sinkUsername:   "wile.e.coyote",
			sinkPassword:   "CarnivorousVulgaris",
		}, want: "---\nreceivers:\n  kafka:\n    brokers:\n    - kafka:9092\n    topic: otlp_metrics-sink-id-222\n    protocol_version: 2.0.0\nextensions:\n  health_check: {}\n  pprof:\n    endpoint: :1888\n  basicauth/exporter:\n    client_auth:\n      username: wile.e.coyote\n      password: CarnivorousVulgaris\nexporters:\n  prometheusremotewrite:\n    endpoint: https://mysinkurl:9922\nservice:\n  extensions:\n  - pprof\n  - health_check\n  - basicauth/exporter\n  pipelines:\n    metrics:\n      receivers:\n      - kafka\n      exporters:\n      - prometheusremotewrite\n",
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReturnConfigYamlFromSink(tt.args.in0, tt.args.kafkaUrlConfig, tt.args.sinkId, tt.args.sinkUrl, tt.args.sinkUsername, tt.args.sinkPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReturnConfigYamlFromSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReturnConfigYamlFromSink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
