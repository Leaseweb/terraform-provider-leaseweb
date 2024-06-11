package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newResources(t *testing.T) {
	sdkCpu := publicCloud.NewCpu()
	sdkCpu.SetUnit("cpu")

	sdkMemory := publicCloud.NewMemory()
	sdkMemory.SetUnit("memory")

	sdkPublicNetworkSpeed := publicCloud.NewPublicNetworkSpeed()
	sdkPublicNetworkSpeed.SetUnit("publicNetworkSpeed")

	sdkPrivateNetworkSpeed := publicCloud.NewPrivateNetworkSpeed()
	sdkPrivateNetworkSpeed.SetUnit("privateNetworkSpeed")

	sdkResources := publicCloud.NewInstanceResources()
	sdkResources.SetCpu(*sdkCpu)
	sdkResources.SetMemory(*sdkMemory)
	sdkResources.SetPublicNetworkSpeed(*sdkPublicNetworkSpeed)
	sdkResources.SetPrivateNetworkSpeed(*sdkPrivateNetworkSpeed)

	resources := newResources(*sdkResources)

	assert.Equal(
		t,
		"cpu",
		resources.Cpu.Unit.ValueString(),
		"cpu should be set",
	)
	assert.Equal(
		t,
		"memory",
		resources.Memory.Unit.ValueString(),
		"memory should be set",
	)
	assert.Equal(
		t,
		"publicNetworkSpeed",
		resources.PublicNetworkSpeed.Unit.ValueString(),
		"publicNetworkSpeed should be set",
	)
	assert.Equal(
		t,
		"privateNetworkSpeed",
		resources.PrivateNetworkSpeed.Unit.ValueString(),
		"privateNetworkSpeed should be set",
	)
}
