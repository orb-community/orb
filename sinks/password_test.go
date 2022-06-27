package sinks

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func Test_passwordService_EncodePassword(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	type fields struct {
		key string
	}
	type args struct {
		plainText string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "encoding with 16 char key",
			fields: fields{
				key: "aaaaaaaaaaaaaaaa",
			},
			args: args{
				"test",
			},
			want: "aVz6gbJ6NhKhzfQ96YW54ys9LEe/ZqDDtAe9Kua8ixU=",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "encoding with 24 char key",
			fields: fields{
				key: "aaaaaaaaaaaaaaaaaaaaaaaa",
			},
			args: args{
				"test",
			},
			want: "ghX4QeVtlLmUD99ZlFKX3nQEpEv/NygQPNpo9eEQeZk=",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "encoding with 32 char key",
			fields: fields{
				key: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
			args: args{
				"test",
			},
			want: "Ze8ovcvKNoT/K+YfXInHQdw9WiZ/NMvBGQnSb1mczAk=",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
		{
			name: "encoding with invalid key",
			fields: fields{
				key: "aaa",
			},
			args: args{
				"test",
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, "invalid key size")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := NewPasswordService(logger, tt.fields.key)
			got, err := ps.EncodePassword(tt.args.plainText)
			if !tt.wantErr(t, err, fmt.Sprintf("EncodePassword(%v)", tt.args.plainText)) {
				return
			}
			assert.Equalf(t, tt.want, got, "EncodePassword(%v)", tt.args.plainText)
		})
	}
}

func Test_passwordService_GetPassword(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	type fields struct {
		key string
	}
	type args struct {
		cipheredText string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "encoding with 16 char key",
			fields: fields{
				key: "aaaaaaaaaaaaaaaa",
			},
			args: args{
				"aVz6gbJ6NhKhzfQ96YW54ys9LEe/ZqDDtAe9Kua8ixU=",
			},
			want: "test",
		},
		{
			name: "encoding with 24 char key",
			fields: fields{
				key: "aaaaaaaaaaaaaaaaaaaaaaaa",
			},
			args: args{
				"ghX4QeVtlLmUD99ZlFKX3nQEpEv/NygQPNpo9eEQeZk=",
			},
			want: "test",
		},
		{
			name: "encoding with 32 char key",
			fields: fields{
				key: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
			args: args{
				"Ze8ovcvKNoT/K+YfXInHQdw9WiZ/NMvBGQnSb1mczAk=",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := NewPasswordService(logger, tt.fields.key)
			assert.Equalf(t, tt.want, ps.GetPassword(tt.args.cipheredText), "GetPassword(%v)", tt.args.cipheredText)
		})
	}
}
