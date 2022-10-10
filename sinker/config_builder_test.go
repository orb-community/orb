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
		}, want: "",
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
