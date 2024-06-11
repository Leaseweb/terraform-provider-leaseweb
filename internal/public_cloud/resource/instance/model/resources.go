package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Resources struct {
	Cpu                 types.Object `tfsdk:"cpu"`
	Memory              types.Object `tfsdk:"memory"`
	PublicNetworkSpeed  types.Object `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed types.Object `tfsdk:"private_network_speed"`
}

func (r Resources) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cpu":                   types.ObjectType{AttrTypes: Cpu{}.attributeTypes()},
		"memory":                types.ObjectType{AttrTypes: Memory{}.attributeTypes()},
		"public_network_speed":  types.ObjectType{AttrTypes: PublicNetworkSpeed{}.attributeTypes()},
		"private_network_speed": types.ObjectType{AttrTypes: PrivateNetworkSpeed{}.attributeTypes()},
	}
}

func newResources(ctx context.Context, sdkResources *publicCloud.InstanceResources) (*Resources, diag.Diagnostics) {
	cpu := newCpu(sdkResources.Cpu)
	cpuObject, diags := types.ObjectValueFrom(ctx, cpu.attributeTypes(), cpu)
	if diags != nil {
		return &Resources{}, diags
	}

	memory := newMemory(sdkResources.Memory)
	memoryObject, diags := types.ObjectValueFrom(ctx, memory.attributeTypes(), memory)
	if diags != nil {
		return &Resources{}, diags
	}

	publicNetworkSpeed := newPublicNetworkSpeed(sdkResources.PublicNetworkSpeed)
	publicNetworkSpeedObject, diags := types.ObjectValueFrom(
		ctx,
		publicNetworkSpeed.attributeTypes(),
		publicNetworkSpeed,
	)
	if diags != nil {
		return &Resources{}, diags
	}

	privateNetworkSpeed := newPrivateNetworkSpeed(sdkResources.PrivateNetworkSpeed)
	privateNetworkSpeedObject, diags := types.ObjectValueFrom(
		ctx,
		privateNetworkSpeed.attributeTypes(),
		privateNetworkSpeed,
	)
	if diags != nil {
		return &Resources{}, diags
	}

	return &Resources{
		Cpu:                 cpuObject,
		Memory:              memoryObject,
		PublicNetworkSpeed:  publicNetworkSpeedObject,
		PrivateNetworkSpeed: privateNetworkSpeedObject,
	}, nil
}
