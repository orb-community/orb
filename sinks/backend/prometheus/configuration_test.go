package prometheus

import (
	"github.com/orb-community/orb/pkg/types"
	"reflect"
	"testing"
)

var (
	validConfiguration = map[string]interface{}{RemoteHostURLConfigFeature: "https://acme.com/prom/push", UsernameConfigFeature: "wile.e.coyote", PasswordConfigFeature: "@secr3t-passw0rd"}
	validYaml          = "remote_host: https://acme.com/prom/push\nusername: wile.e.coyote\npassword: \"@secr3t-passw0rd\""
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
				config: map[string]interface{}{RemoteHostURLConfigFeature: "acme.com/prom/push", UsernameConfigFeature: "wile.e.coyote", PasswordConfigFeature: "@secr3t-passw0rd"},
			},
			wantErr: true,
		},
		{
			name: "missing host configuration",
			args: args{
				config: map[string]interface{}{UsernameConfigFeature: "wile.e.coyote", PasswordConfigFeature: "@secr3t-passw0rd"},
			},
			wantErr: true,
		},
		{
			name: "missing username configuration",
			args: args{
				config: map[string]interface{}{RemoteHostURLConfigFeature: "acme.com/prom/push", PasswordConfigFeature: "@secr3t-passw0rd"},
			},
			wantErr: true,
		},
		{
			name: "missing password configuration",
			args: args{
				config: map[string]interface{}{RemoteHostURLConfigFeature: "acme.com/prom/push", UsernameConfigFeature: "wile.e.coyote"},
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
	pass := "@secr3t-passw0rd"
	user := "wile.e.coyote"
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
			wantConfigReturn: map[string]interface{}{RemoteHostURLConfigFeature: "https://acme.com/prom/push", UsernameConfigFeature: &user, PasswordConfigFeature: &pass},
			wantErr:          false,
		},
		{
			name: "invalid parse",
			args: args{
				format: "yaml",
				config: "remote_host: https://acme.com/prom/push\nusername: wile.e.coyote\npassword \"@secr3t-passw0rd\"",
			},
			wantConfigReturn: map[string]interface{}{RemoteHostURLConfigFeature: "https://acme.com/prom/push", UsernameConfigFeature: &user, PasswordConfigFeature: &pass},
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
