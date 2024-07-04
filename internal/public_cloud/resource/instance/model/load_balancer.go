package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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
	sdkLoadBalancerDetails publicCloud.LoadBalancerDetails,
) (*LoadBalancer, diag.Diagnostics) {
	resourcesObject, diags := utils.ConvertSdkModelToResourceObject(
		sdkLoadBalancerDetails.GetResources(),
		Resources{}.AttributeTypes(),
		ctx,
		newResources,
	)
	if diags.HasError() {
		return nil, diags
	}

	contractObject, diags := utils.ConvertSdkModelToResourceObject(
		sdkLoadBalancerDetails.GetContract(),
		Contract{}.AttributeTypes(),
		ctx,
		newContract,
	)
	if diags.HasError() {
		return nil, diags
	}

	loadBalancerConfigurationObject, diags := utils.ConvertSdkModelToResourceObject(
		sdkLoadBalancerDetails.GetConfiguration(),
		LoadBalancerConfiguration{}.AttributeTypes(),
		ctx,
		newLoadBalancerConfiguration,
	)
	if diags.HasError() {
		return nil, diags
	}

	privateNetworkObject, diags := utils.ConvertSdkModelToResourceObject(
		sdkLoadBalancerDetails.GetPrivateNetwork(),
		PrivateNetwork{}.AttributeTypes(),
		ctx,
		newPrivateNetwork,
	)
	if diags.HasError() {
		return nil, diags
	}

	var ips []Ip
	for _, ip := range sdkLoadBalancerDetails.Ips {
		ipObject, diags := newIp(ctx, &ip)
		if diags != nil {
			return nil, diags
		}
		ips = append(ips, ipObject)
	}
	ipsObject, diags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: Ip{}.AttributeTypes()},
		ips,
	)
	if diags != nil {
		return nil, diags
	}

	return &LoadBalancer{
		Id:                        basetypes.NewStringValue(sdkLoadBalancerDetails.GetId()),
		Type:                      basetypes.NewStringValue(sdkLoadBalancerDetails.GetType()),
		Resources:                 resourcesObject,
		Region:                    basetypes.NewStringValue(sdkLoadBalancerDetails.GetRegion()),
		Reference:                 basetypes.NewStringValue(sdkLoadBalancerDetails.GetReference()),
		State:                     basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetState())),
		Contract:                  contractObject,
		StartedAt:                 basetypes.NewStringValue(sdkLoadBalancerDetails.GetStartedAt().String()),
		LoadBalancerConfiguration: loadBalancerConfigurationObject,
		PrivateNetwork:            privateNetworkObject,
		Ips:                       ipsObject,
	}, nil
}
