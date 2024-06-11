package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

	resources, _ := newResources(context.TODO(), sdkResources)

	cpu := Cpu{}
	resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(t, "cpu", cpu.Unit.ValueString(), "cpu should be set")

	memory := Memory{}
	resources.Memory.As(context.TODO(), &memory, basetypes.ObjectAsOptions{})
	assert.Equal(t, "memory", memory.Unit.ValueString(), "memory should be set")

	publicNetworkSpeed := PublicNetworkSpeed{}
	resources.PublicNetworkSpeed.As(context.TODO(), &publicNetworkSpeed, basetypes.ObjectAsOptions{})
	assert.Equal(t, "publicNetworkSpeed", publicNetworkSpeed.Unit.ValueString(), "publicNetworkSpeed should be set")

	privateNetworkSpeed := PrivateNetworkSpeed{}
	resources.PrivateNetworkSpeed.As(context.TODO(), &privateNetworkSpeed, basetypes.ObjectAsOptions{})
	assert.Equal(t, "privateNetworkSpeed", privateNetworkSpeed.Unit.ValueString(), "privateNetworkSpeed should be set")
}

func TestResources_attributeTypes(t *testing.T) {
	sdkResources := publicCloud.NewInstanceResources()
	resources, _ := newResources(context.TODO(), sdkResources)

	_, diags := types.ObjectValueFrom(
		context.TODO(),
		resources.attributeTypes(),
		resources,
	)

	assert.Nil(t, diags, "attributes should be correct")
}
