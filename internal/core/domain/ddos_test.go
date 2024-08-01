package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDdos(t *testing.T) {
	ddos := NewDdos("detectionProfile", "protectionType")

	assert.Equal(t, "detectionProfile", ddos.DetectionProfile)
	assert.Equal(t, "protectionType", ddos.ProtectionType)
}
