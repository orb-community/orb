package otlpexporter

import (
	"github.com/orb-community/orb/pkg/types"
	"reflect"
	"testing"
)

func TestBackend_ParseConfig(t *testing.T) {
	type args struct {
		format string
		config string
	}
	tests := []struct {
		name          string
		args          args
		wantRetConfig types.Metadata
		wantErr       bool
	}{
		{
			name: "valid secure but skip verify config",
			args: args{
				format: "yaml",
				config: `endpoint: myserver.local:55690
tls:  
  insecure: false
  insecure_skip_verify: true
auth:
  username: test
  password: test`,
			},
			wantRetConfig: types.Metadata{
				"endpoint": "myserver.local:55690",
				"tls": types.Metadata{
					"insecure":             false,
					"insecure_skip_verify": true,
				},
				"auth": types.Metadata{
					"username": "test",
					"password": "test",
				},
			},
			wantErr: false,
		},
		{
			name: "valid secure but skip verify config",
			args: args{
				format: "yaml",
				config: `endpoint: myserver.local:55690
tls:
  insecure: false
  ca_file: server.crt
  cert_file: client.crt
  key_file: client.key
  min_version: "1.1"
  max_version: "1.2"
auth:
  username: test
  password: test`,
			},
			wantRetConfig: types.Metadata{
				"endpoint": "myserver.local:55690",
				"tls": types.Metadata{
					"insecure":    false,
					"ca_file":     "server.crt",
					"cert_file":   "client.crt",
					"key_file":    "client.key",
					"min_version": "1.1",
					"max_version": "1.2",
				},
				"auth": types.Metadata{
					"username": "test",
					"password": "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Backend{}
			gotRetConfig, err := b.ParseConfig(tt.args.format, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRetConfig, tt.wantRetConfig) {
				t.Errorf("ParseConfig() gotRetConfig = \n%v\n, want \n%v\n", gotRetConfig, tt.wantRetConfig)
			}
		})
	}
}
