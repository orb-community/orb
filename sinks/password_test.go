package sinks

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func Test_passwordService_EncodePassword(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name          string
		key           string
		plainText     string
		encodedString string
	}{
		{
			name:          "with 32 char key",
			key:           "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			plainText:     "test",
			encodedString: "bbf4b204e5daea6e7cb4cb8dec2011c91de502db08c1fc37f4e1ba8b8da60cf0",
		},
		{
			name:          "with smaller key",
			key:           "testing",
			plainText:     "test",
			encodedString: "c8dd6f7f76d1b988574559959c68615ae72487b13bef2f7c4afbce204cc11864",
		},
		{
			name:          "with uuid-key",
			key:           "eb1bc7f4-2031-41c4-85fa-2ddce3abfc3b",
			plainText:     "test",
			encodedString: "1f1114dd9e7953585a768d280a3d0f8592647e0761d085bfa83b9b57c2110a5c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := NewPasswordService(logger, tt.key)
			got, err := ps.EncodePassword(tt.plainText)
			if err != nil {
				t.Fatalf("received error on encoding password: %e", err)
			}
			t.Logf("storing %s", got)
			password, err := ps.DecodePassword(got)
			if err != nil {
				t.Fatalf("received error on decoding password: %e", err)
			}
			t.Logf("retrieving %s", password)
			assert.Equalf(t, tt.plainText, password, "Got Decoded Password %s", password)
			getPassword, err := ps.DecodePassword(tt.encodedString)
			if err != nil {
				t.Fatalf("received error on decoding stored password: %e", err)
			}
			t.Logf("retrieving %s", getPassword)
			assert.Equalf(t, getPassword, password, "Stored coded password is %s", getPassword)
		})
	}
}
