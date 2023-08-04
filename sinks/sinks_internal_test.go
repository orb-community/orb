package sinks

import (
	"fmt"
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/authentication_type/basicauth"
	"github.com/orb-community/orb/sinks/backend/otlphttpexporter"
	"github.com/orb-community/orb/sinks/backend/prometheus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func Test_sinkService_validateBackend(t *testing.T) {
	logger := zap.NewNop()
	otlphttpexporter.Register()
	prometheus.Register()
	passwordService := authentication_type.NewPasswordService(logger, "unit-test")
	basicauth.Register(passwordService)
	type fields struct {
		svc sinkService
	}
	type args struct {
		sink *Sink
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantBe  reflect.Type
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "prometheus over yaml",
			fields: fields{
				svc: sinkService{
					logger: logger,
				},
			},
			args: args{
				sink: &Sink{
					Backend:    "prometheus",
					Config:     nil,
					Format:     "yaml",
					ConfigData: "authentication:\n  type: basicauth\n  password: \"password\"\n  username: \"user\"\nexporter:\n  remote_host: \"https://acme.com/api/prom/push\"\nopentelemetry: enabled\n",
				},
			},
			wantBe:  reflect.TypeOf(&prometheus.Backend{}),
			wantErr: nil,
		},
		{
			name: "otlphttp over yaml",
			fields: fields{
				svc: sinkService{
					logger: logger,
				},
			},
			args: args{
				sink: &Sink{
					Backend:    "otlphttp",
					Config:     nil,
					Format:     "yaml",
					ConfigData: "authentication:\n  type: basicauth\n  password: \"password\"\n  username: \"user\"\nexporter:\n  endpoint: \"https://acme.com/api/prom/push\"\nopentelemetry: enabled\n",
				},
			},
			wantBe:  reflect.TypeOf(&otlphttpexporter.OTLPHTTPBackend{}),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBe, err := tt.fields.svc.validateBackend(tt.args.sink)
			if tt.wantErr != nil && !tt.wantErr(t, err, fmt.Sprintf("validateBackend(%v)", tt.args.sink)) {
				return
			}
			assert.Equalf(t, tt.wantBe, reflect.TypeOf(gotBe), "validateBackend(%v)", tt.args.sink.Backend)
		})
	}
}
