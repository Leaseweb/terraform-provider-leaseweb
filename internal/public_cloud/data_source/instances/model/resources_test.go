package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

func Test_newResources(t *testing.T) {
	entityResources := entity.NewResources(
		entity.Cpu{Unit: "cpu"},
		entity.Memory{Unit: "memory"},
		entity.NetworkSpeed{Unit: "publicNetworkSpeed"},
		entity.NetworkSpeed{Unit: "NetworkSpeed"},
	)

	got := newResources(entityResources)

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
