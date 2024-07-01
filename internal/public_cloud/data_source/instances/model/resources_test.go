package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newResources(t *testing.T) {
	sdkResources := publicCloud.NewResources(
		publicCloud.Cpu{Unit: "cpu"},
		publicCloud.Memory{Unit: "memory"},
		publicCloud.NetworkSpeed{Unit: "publicNetworkSpeed"},
		publicCloud.NetworkSpeed{Unit: "NetworkSpeed"},
	)

	got := newResources(*sdkResources)

	assert.Equal(
		t,
		"cpu",
		got.Cpu.Unit.ValueString(),
		"cpu should be set",
	)
	assert.Equal(
		t,
		"memory",
		got.Memory.Unit.ValueString(),
		"memory should be set",
	)
	assert.Equal(
		t,
		"publicNetworkSpeed",
		got.PublicNetworkSpeed.Unit.ValueString(),
		"publicNetworkSpeed should be set",
	)
	assert.Equal(
		t,
		"NetworkSpeed",
		got.PrivateNetworkSpeed.Unit.ValueString(),
		"NetworkSpeed should be set",
	)
}
