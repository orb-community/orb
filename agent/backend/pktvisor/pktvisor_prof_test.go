package pktvisor

import (
	"context"
	"testing"
)

func Test_pktvisorBackend_scrapeOpenTelemetry(t *testing.T) {

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p pktvisorBackend
			p.scrapeOpenTelemetry(tt.args.ctx)
		})
	}
}
