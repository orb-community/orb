package http

import (
	"fmt"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/backend/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_updateSinkReq_validate(t *testing.T) {
	promBe := prometheus.Backend{}
	aDescription := "a description worth reading"
	type fields struct {
		Name        string
		Config      types.Metadata
		Backend     string
		Format      string
		ConfigData  string
		Description *string
		Tags        types.Tags
		id          string
		token       string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "full update no yaml",
			fields: fields{
				Name:        "new-name",
				Config:      map[string]interface{}{"username": "wile.e.coyote", "password": "C@rnivurousVulgar1s", "remote_host": "https://acme.com/prom/push"},
				Backend:     "prometheus",
				Description: &aDescription,
				Tags:        map[string]string{"cloud": "aws", "region": "us-east-1"},
				id:          "1122",
				token:       "valid-token",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "full update yaml",
			fields: fields{
				Name:        "new-name",
				Backend:     "prometheus",
				Format:      "yaml",
				ConfigData:  "remote_host: https://acme.com/prom/push\nusername: wile.e.coyote\npassword: \"@DesertL00kingForMeal\"",
				Description: &aDescription,
				Tags:        map[string]string{"cloud": "aws", "region": "us-east-1"},
				id:          "1122",
				token:       "valid-token",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "only name update",
			fields: fields{
				Name:  "new-name",
				id:    "1122",
				token: "valid-token",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "description update",
			fields: fields{
				Description: &aDescription,
				id:          "1122",
				token:       "valid-token",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "json config update",
			fields: fields{
				Config: map[string]interface{}{"username": "wile.e.coyote", "password": "C@rnivurousVulgar1s", "remote_host": "https://acme.com/prom/push"},
				id:     "1122",
				token:  "valid-token",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "yaml config update",
			fields: fields{
				Format:     "yaml",
				ConfigData: "remote_host: https://acme.com/prom/push\nusername: wile.e.coyote\npassword: \"@DesertL00kingForMeal\"",
				id:         "1122",
				token:      "valid-token",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "tags update",
			fields: fields{
				Tags:  map[string]string{"cloud": "aws", "region": "us-east-1"},
				id:    "1122",
				token: "valid-token",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := updateSinkReq{
				Name:        tt.fields.Name,
				Config:      tt.fields.Config,
				Backend:     tt.fields.Backend,
				Format:      tt.fields.Format,
				ConfigData:  tt.fields.ConfigData,
				Description: tt.fields.Description,
				Tags:        tt.fields.Tags,
				id:          tt.fields.id,
				token:       tt.fields.token,
			}
			tt.wantErr(t, req.validate(), fmt.Sprintf("validate(%v)", promBe))
		})
	}
}
