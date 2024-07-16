package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newResources(t *testing.T) {
	resources := domain.NewResources(
		domain.Cpu{Unit: "cpu"},
		domain.Memory{Unit: "memory"},
		domain.NetworkSpeed{Unit: "publicNetworkSpeed"},
		domain.NetworkSpeed{Unit: "NetworkSpeed"},
	)

	got := newResources(resources)

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
