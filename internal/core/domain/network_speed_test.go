package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNetworkSpeed(t *testing.T) {
	networkSpeed := NewNetworkSpeed(1, "unit")

	assert.Equal(t, 1, networkSpeed.Value)
	assert.Equal(t, "unit", networkSpeed.Unit)
}
