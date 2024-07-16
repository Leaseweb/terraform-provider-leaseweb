package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type Resources struct {
	Cpu                 types.Object `tfsdk:"cpu"`
	Memory              types.Object `tfsdk:"memory"`
	PublicNetworkSpeed  types.Object `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed types.Object `tfsdk:"private_network_speed"`
}

func (r Resources) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cpu":                   types.ObjectType{AttrTypes: Cpu{}.AttributeTypes()},
		"memory":                types.ObjectType{AttrTypes: Memory{}.AttributeTypes()},
		"public_network_speed":  types.ObjectType{AttrTypes: NetworkSpeed{}.AttributeTypes()},
		"private_network_speed": types.ObjectType{AttrTypes: NetworkSpeed{}.AttributeTypes()},
	}
}

func newResources(
	ctx context.Context,
	entityResources entity.Resources,
) (*Resources, diag.Diagnostics) {
	cpu := newCpu(entityResources.Cpu)
	cpuObject, diags := types.ObjectValueFrom(ctx, cpu.AttributeTypes(), cpu)
	if diags != nil {
		return &Resources{}, diags
	}

	memory := newMemory(entityResources.Memory)
	memoryObject, diags := types.ObjectValueFrom(ctx, memory.AttributeTypes(), memory)
	if diags != nil {
		return &Resources{}, diags
	}

	publicNetworkSpeed := newNetworkSpeed(entityResources.PublicNetworkSpeed)
	publicNetworkSpeedObject, diags := types.ObjectValueFrom(
		ctx,
		publicNetworkSpeed.AttributeTypes(),
		publicNetworkSpeed,
	)
	if diags != nil {
		return &Resources{}, diags
	}

	privateNetworkSpeed := newNetworkSpeed(entityResources.PrivateNetworkSpeed)
	privateNetworkSpeedObject, diags := types.ObjectValueFrom(
		ctx,
		privateNetworkSpeed.AttributeTypes(),
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
