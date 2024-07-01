package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newResources(t *testing.T) {

	sdkResources := publicCloud.NewResources(
		publicCloud.Cpu{Unit: "cpu"},
		publicCloud.Memory{Unit: "memory"},
		publicCloud.NetworkSpeed{Unit: "publicNetworkSpeed"},
		publicCloud.NetworkSpeed{Unit: "privateNetworkSpeed"},
	)

	resources, _ := newResources(context.TODO(), *sdkResources)

	cpu := Cpu{}
	resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"cpu",
		cpu.Unit.ValueString(),
		"cpu should be set",
	)

	memory := Memory{}
	resources.Memory.As(context.TODO(), &memory, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"memory",
		memory.Unit.ValueString(),
		"memory should be set",
	)

	publicNetworkSpeed := NetworkSpeed{}
	resources.PublicNetworkSpeed.As(
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
	resources.PrivateNetworkSpeed.As(
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
	resources, _ := newResources(context.TODO(), publicCloud.Resources{})

	_, diags := types.ObjectValueFrom(
		context.TODO(),
		resources.AttributeTypes(),
		resources,
	)

	assert.Nil(t, diags, "attributes should be correct")
}
