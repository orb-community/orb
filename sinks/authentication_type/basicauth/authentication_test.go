package basicauth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
)

func TestAuthConfig_ValidateConfiguration(t *testing.T) {
	type args struct {
		inputFormat string
		input       types.Metadata
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "missing_username",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"password": "test-password",
				},
			},
			wantErr: errors.ErrAuthUsernameNotFound,
		},
		{
			name: "empty_username",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"username": "",
					"password": "test-password",
				},
			},
			wantErr: errors.ErrAuthInvalidUsernameType,
		},
		{
			name: "invalid_username_type",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"username": 1234,
					"password": "test-password",
				},
			},
			wantErr: errors.ErrAuthInvalidUsernameType,
		},
		{
			name: "invalid_username",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"username": " ",
					"password": "test-password",
				},
			},
			wantErr: errors.ErrAuthInvalidUsernameType,
		},
		{
			name: "missing_password",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"username": "test-user",
				},
			},
			wantErr: errors.ErrAuthPasswordNotFound,
		},
		{
			name: "empty_password",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"username": "test-user",
					"password": "",
				},
			},
			wantErr: errors.ErrAuthInvalidPasswordType,
		},
		{
			name: "invalid_password_type",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"username": "test-user",
					"password": 1234,
				},
			},
			wantErr: errors.ErrAuthInvalidPasswordType,
		},
		{
			name: "invalid_password",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"username": "test-user",
					"password": " ",
				},
			},
			wantErr: errors.ErrAuthInvalidPasswordType,
		},
		{
			name: "valid",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"username": "test-user",
					"password": "test-password",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AuthConfig
			err := a.ValidateConfiguration(tt.args.inputFormat, tt.args.input)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
