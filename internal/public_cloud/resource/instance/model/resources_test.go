package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

func Test_newResources(t *testing.T) {
	resources := entity.NewResources(
		entity.Cpu{Unit: "cpu"},
		entity.Memory{Unit: "memory"},
		entity.NetworkSpeed{Unit: "publicNetworkSpeed"},
		entity.NetworkSpeed{Unit: "privateNetworkSpeed"},
	)

	got, _ := newResources(context.TODO(), resources)

	cpu := Cpu{}
	got.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"cpu",
		cpu.Unit.ValueString(),
		"cpu should be set",
	)

	memory := Memory{}
	got.Memory.As(context.TODO(), &memory, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"memory",
		memory.Unit.ValueString(),
		"memory should be set",
	)

	publicNetworkSpeed := NetworkSpeed{}
	got.PublicNetworkSpeed.As(
		context.TODO(),
		&publicNetworkSpeed,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"publicNetworkSpeed",
		publicNetworkSpeed.Unit.ValueString(),
		"publicNetworkSpeed should be set",
	)

	privateNetworkSpeed := NetworkSpeed{}
	got.PrivateNetworkSpeed.As(
		context.TODO(),
		&privateNetworkSpeed,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"privateNetworkSpeed",
		privateNetworkSpeed.Unit.ValueString(),
		"privateNetworkSpeed should be set",
	)
}

func TestResources_attributeTypes(t *testing.T) {
	resources, _ := newResources(context.TODO(), entity.Resources{})

	_, diags := types.ObjectValueFrom(
		context.TODO(),
		resources.AttributeTypes(),
		resources,
	)

	assert.Nil(t, diags, "attributes should be correct")
}
