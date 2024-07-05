package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSshKey(t *testing.T) {
	t.Run("valid key is set properly", func(t *testing.T) {
		sshKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

		got, err := NewSshKey(sshKey)

		assert.Nil(t, err)
		assert.Equal(t, sshKey, got.value)
	})

	t.Run("invalid key returns error", func(t *testing.T) {
		_, err := NewSshKey("tralala")

		assert.NotNil(t, err)
	})

	t.Run("struct string value is correct", func(t *testing.T) {
		sshKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

		got, _ := NewSshKey(sshKey)

		assert.Equal(t, sshKey, got.String())
	})
}
