package prometheus

import (
	"github.com/orb-community/orb/pkg/types"
	"reflect"
	"testing"
)

var (
	validConfiguration = map[string]interface{}{RemoteHostURLConfigFeature: "https://acme.com/prom/push"}
	validYaml          = "remote_host: https://acme.com/prom/push"
)

func TestBackend_ValidateConfiguration(t *testing.T) {
	type args struct {
		config types.Metadata
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid configuration",
			args: args{
				config: validConfiguration,
			},
			wantErr: false,
		},
		{
			name: "invalid host configuration",
			args: args{
				config: map[string]interface{}{RemoteHostURLConfigFeature: "acme.com/prom/push"},
			},
			wantErr: true,
		},
		{
			name: "missing host configuration",
			args: args{
				config: map[string]interface{}{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Backend{}
			if err := p.ValidateConfiguration(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackend_ParseConfig(t *testing.T) {
	type args struct {
		format string
		config string
	}
	tests := []struct {
		name             string
		args             args
		wantConfigReturn types.Metadata
		wantErr          bool
	}{
		{
			name: "valid parse",
			args: args{
				format: "yaml",
				config: validYaml,
			},
			wantConfigReturn: map[string]interface{}{"exporter": map[string]interface{}{RemoteHostURLConfigFeature: "https://acme.com/prom/push"}},
			wantErr:          false,
		},
		{
			name: "invalid parse",
			args: args{
				format: "yaml",
				config: "remote_host: \nhttps://acme.com/prom/push\n\n",
			},
			wantConfigReturn: map[string]interface{}{RemoteHostURLConfigFeature: "https://acme.com/prom/push"},
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Backend{}
			gotConfigReturn, err := p.ParseConfig(tt.args.format, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(gotConfigReturn, tt.wantConfigReturn) {
				t.Errorf("ParseConfig() gotConfigReturn = %v, want %v", gotConfigReturn, tt.wantConfigReturn)
			}
		})
	}
}

func TestBackend_CreateFeatureConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "valid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Backend{}
			got := p.CreateFeatureConfig()
			remoteHostOk := false
			for _, feature := range got {
				if feature.Name == RemoteHostURLConfigFeature {
					remoteHostOk = true
				}
			}
			if remoteHostOk {
				return
			} else {
				t.Fail()
			}
		})
	}
}
