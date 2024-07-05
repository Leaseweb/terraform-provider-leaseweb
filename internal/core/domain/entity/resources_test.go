package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResources(t *testing.T) {
	cpu := NewCpu(1, "")
	memory := NewMemory(1, "")
	publicNetworkSpeed := NewNetworkSpeed(1, "")
	privateNetworkSpeed := NewNetworkSpeed(1, "")

	resources := NewResources(
		cpu,
		memory,
		publicNetworkSpeed,
		privateNetworkSpeed,
	)

	assert.Equal(t, cpu, resources.Cpu)
	assert.Equal(t, memory, resources.Memory)
	assert.Equal(t, publicNetworkSpeed, resources.PublicNetworkSpeed)
	assert.Equal(t, privateNetworkSpeed, resources.PrivateNetworkSpeed)
}
