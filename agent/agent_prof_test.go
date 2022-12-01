package agent

import (
	"context"
	"testing"
)

func Test_orbAgent_startBackends(t *testing.T) {

	type args struct {
		agentCtx context.Context
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := orbAgent{}
			if err := a.startBackends(tt.args.agentCtx); (err != nil) != tt.wantErr {
				t.Errorf("startBackends() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
