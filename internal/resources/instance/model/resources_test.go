package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newResources(t *testing.T) {
	cpu := publicCloud.NewCpu()
	cpu.SetUnit("cpu")

	memory := publicCloud.NewMemory()
	memory.SetUnit("memory")

	publicNetworkSpeed := publicCloud.NewPublicNetworkSpeed()
	publicNetworkSpeed.SetUnit("publicNetworkSpeed")

	privateNetworkSpeed := publicCloud.NewPrivateNetworkSpeed()
	privateNetworkSpeed.SetUnit("privateNetworkSpeed")

	sdkResources := publicCloud.NewInstanceResources()
	sdkResources.SetCpu(*cpu)
	sdkResources.SetMemory(*memory)
	sdkResources.SetPublicNetworkSpeed(*publicNetworkSpeed)
	sdkResources.SetPrivateNetworkSpeed(*privateNetworkSpeed)

	resources := newResources(sdkResources)

	assert.Equal(t, "cpu", resources.Cpu.Unit.ValueString(), "cpu should be set")
	assert.Equal(t, "memory", resources.Memory.Unit.ValueString(), "memory should be set")
	assert.Equal(t, "publicNetworkSpeed", resources.PublicNetworkSpeed.Unit.ValueString(), "publicNetworkSpeed should be set")
	assert.Equal(t, "privateNetworkSpeed", resources.PrivateNetworkSpeed.Unit.ValueString(), "privateNetworkSpeed should be set")
}
