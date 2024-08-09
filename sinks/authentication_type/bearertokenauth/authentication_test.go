package bearertokenauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type"
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
			name: "missing_schema",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"token": "test_api_key",
				},
			},
			wantErr: errors.ErrAuthSchemeNotFound,
		},
		{
			name: "empty_schema",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"scheme": "",
					"token":  "test_api_key",
				},
			},
			wantErr: errors.ErrAuthInvalidSchemeType,
		},
		{
			name: "invalid_schema_type",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"scheme": 1234,
					"token":  "test_api_key",
				},
			},
			wantErr: errors.ErrAuthInvalidSchemeType,
		},
		{
			name: "invalid_schema",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"scheme": "fake schema",
					"token":  "test_api_key",
				},
			},
			wantErr: errors.ErrAuthInvalidSchemeType,
		},
		{
			name: "missing_token",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"scheme": "Bearer",
				},
			},
			wantErr: errors.ErrAuthTokenNotFound,
		},
		{
			name: "empty_token",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"scheme": "Bearer",
					"token":  "",
				},
			},
			wantErr: errors.ErrAuthInvalidTokenType,
		},
		{
			name: "invalid_token_type",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"scheme": "Bearer",
					"token":  1234,
				},
			},
			wantErr: errors.ErrAuthInvalidTokenType,
		},
		{
			name: "invalid_token_value",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"scheme": "Bearer",
					"token":  "invalid api key",
				},
			},
			wantErr: errors.ErrAuthInvalidTokenType,
		},
		{
			name: "valid",
			args: args{
				inputFormat: "object",
				input: types.Metadata{
					"scheme": "Bearer",
					"token":  "abcdefg",
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

func TestAuthConfig_OmitInformation(t *testing.T) {
	t.Run("invalid output format", func(t *testing.T) {
		input := types.Metadata{
			"authentication": types.Metadata{
				"scheme": "Bearer",
				"token":  "abcdefg",
			},
		}

		var a AuthConfig

		_, err := a.OmitInformation("blah", input)
		assert.Error(t, err)
	})
	t.Run("successfully stripped the secret token", func(t *testing.T) {
		input := types.Metadata{
			"authentication": types.Metadata{
				"scheme": "Bearer",
				"token":  "abcdefg",
			},
		}

		want := types.Metadata{
			"authentication": types.Metadata{
				"scheme": "Bearer",
				"token":  "",
			},
		}

		var a AuthConfig

		got, err := a.OmitInformation("object", input)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestAuthConfig_EncodeInformation(t *testing.T) {
	t.Run("invalid output format", func(t *testing.T) {
		input := types.Metadata{
			"authentication": types.Metadata{
				"scheme": "Bearer",
				"token":  "abcdefg",
			},
		}

		a := AuthConfig{
			encryptionService: authentication_type.NewPasswordService(nil, "test"),
		}

		_, err := a.EncodeInformation("blah", input)
		assert.Error(t, err)
	})

	t.Run("successfully encrypted token", func(t *testing.T) {
		input := types.Metadata{
			"authentication": types.Metadata{
				"scheme": "Bearer",
				"token":  "abcdefg",
			},
		}

		a := AuthConfig{
			encryptionService: authentication_type.NewPasswordService(nil, "test"),
		}

		_, err := a.EncodeInformation("object", input)
		require.NoError(t, err)
	})
}

func TestAuthConfig_DecodeInformation(t *testing.T) {
	t.Run("invalid output format", func(t *testing.T) {
		input := types.Metadata{
			"authentication": types.Metadata{
				"scheme": "Bearer",
				"token":  "dca8757dee5dfcc592c97355396dc2bdb95c6a3f58d4acb4453717c960827602acaf49",
			},
		}

		a := AuthConfig{
			encryptionService: authentication_type.NewPasswordService(nil, "test"),
		}

		_, err := a.DecodeInformation("blah", input)
		assert.Error(t, err)
	})

	t.Run("successfully decrypted token", func(t *testing.T) {
		input := types.Metadata{
			"authentication": types.Metadata{
				"scheme": "Bearer",
				"token":  "dca8757dee5dfcc592c97355396dc2bdb95c6a3f58d4acb4453717c960827602acaf49",
			},
		}

		want := types.Metadata{
			"authentication": types.Metadata{
				"scheme": "Bearer",
				"token":  "abcdefg",
			},
		}

		a := AuthConfig{
			encryptionService: authentication_type.NewPasswordService(nil, "test"),
		}

		got, err := a.DecodeInformation("object", input)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
