package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/utils"
)

type LoadBalancer struct {
	Id                        types.String `tfsdk:"id"`
	Type                      types.String `tfsdk:"type"`
	Resources                 types.Object `tfsdk:"resources"`
	Region                    types.String `tfsdk:"region"`
	Reference                 types.String `tfsdk:"reference"`
	State                     types.String `tfsdk:"state"`
	Contract                  types.Object `tfsdk:"contract"`
	StartedAt                 types.String `tfsdk:"started_at"`
	Ips                       types.List   `tfsdk:"ips"`
	LoadBalancerConfiguration types.Object `tfsdk:"load_balancer_configuration"`
	PrivateNetwork            types.Object `tfsdk:"private_network"`
}

func (l LoadBalancer) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                          types.StringType,
		"type":                        types.StringType,
		"resources":                   types.ObjectType{AttrTypes: Resources{}.AttributeTypes()},
		"region":                      types.StringType,
		"reference":                   types.StringType,
		"state":                       types.StringType,
		"contract":                    types.ObjectType{AttrTypes: Contract{}.AttributeTypes()},
		"started_at":                  types.StringType,
		"ips":                         types.ListType{ElemType: types.ObjectType{AttrTypes: Ip{}.AttributeTypes()}},
		"load_balancer_configuration": types.ObjectType{AttrTypes: LoadBalancerConfiguration{}.AttributeTypes()},
		"private_network":             types.ObjectType{AttrTypes: PrivateNetwork{}.AttributeTypes()},
	}
}

func newLoadBalancer(
	ctx context.Context,
	entityLoadBalancer entity.LoadBalancer,
) (*LoadBalancer, diag.Diagnostics) {
	resources, diags := utils.ConvertDomainEntityToResourceObject(
		entityLoadBalancer.Resources,
		Resources{}.AttributeTypes(),
		ctx,
		newResources,
	)
	if diags.HasError() {
		return nil, diags
	}

	contract, diags := utils.ConvertDomainEntityToResourceObject(
		entityLoadBalancer.Contract,
		Contract{}.AttributeTypes(),
		ctx,
		newContract,
	)
	if diags.HasError() {
		return nil, diags
	}

	configuration, diags := utils.ConvertNullableDomainEntityToResourceObject(
		entityLoadBalancer.Configuration,
		LoadBalancerConfiguration{}.AttributeTypes(),
		ctx,
		newLoadBalancerConfiguration,
	)
	if diags.HasError() {
		return nil, diags
	}

	privateNetwork, diags := utils.ConvertNullableDomainEntityToResourceObject(
		entityLoadBalancer.PrivateNetwork,
		PrivateNetwork{}.AttributeTypes(),
		ctx,
		newPrivateNetwork,
	)
	if diags.HasError() {
		return nil, diags
	}

	ips, diags := utils.ConvertEntitiesToListValue(
		entityLoadBalancer.Ips,
		Ip{}.AttributeTypes(),
		ctx,
		newIp,
	)
	if diags != nil {
		return nil, diags
	}

	return &LoadBalancer{
		Id:                        basetypes.NewStringValue(entityLoadBalancer.Id.String()),
		Type:                      basetypes.NewStringValue(entityLoadBalancer.Type),
		Resources:                 resources,
		Region:                    basetypes.NewStringValue(entityLoadBalancer.Region),
		Reference:                 utils.ConvertNullableStringToStringValue(entityLoadBalancer.Reference),
		State:                     basetypes.NewStringValue(string(entityLoadBalancer.State)),
		Contract:                  contract,
		StartedAt:                 basetypes.NewStringValue(entityLoadBalancer.StartedAt.String()),
		LoadBalancerConfiguration: configuration,
		PrivateNetwork:            privateNetwork,
		Ips:                       ips,
	}, nil
}
