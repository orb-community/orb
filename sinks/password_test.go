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
			encodedString: "589e5c8bc26165a515697ec4231ddcb2bb93162e5ab3eac6e60685914b03a26b",
		},
		{
			name:          "with smaller key",
			key:           "testing",
			plainText:     "test",
			encodedString: "daa753d663f98ab0825ad0a3fd61c8956a05bbb573e7e7e46091bd043c4161e8",
		},
		{
			name:          "with uuid-key",
			key:           "eb1bc7f4-2031-41c4-85fa-2ddce3abfc3b",
			plainText:     "test",
			encodedString: "8c6ac82169d7228c9880f7276a9b06a2bfa0288a8586ad47c2807896c8e63d42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := NewPasswordService(logger, tt.key)
			got := ps.EncodePassword(tt.plainText)
			t.Logf("storing %s", got)
			password := ps.GetPassword(got)
			t.Logf("retrieving %s", password)
			assert.Equalf(t, tt.plainText, password, "Got Decoded Password %s", password)
			getPassword := ps.GetPassword(tt.encodedString)
			t.Logf("retrieving %s", getPassword)
			assert.Equalf(t, getPassword, password, "Stored coded password is %s", getPassword)
		})
	}
}
