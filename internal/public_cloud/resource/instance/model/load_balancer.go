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
	Id        types.String `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Resources types.Object `tfsdk:"resources"`
	Region    types.String `tfsdk:"region"`
	Reference types.String `tfsdk:"reference"`
	State     types.String `tfsdk:"state"`
	Contract  types.Object `tfsdk:"contract"`
	StartedAt types.String `tfsdk:"started_at"`
}

func (l LoadBalancer) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"type":       types.StringType,
		"resources":  types.ObjectType{AttrTypes: Resources{}.AttributeTypes()},
		"region":     types.StringType,
		"reference":  types.StringType,
		"state":      types.StringType,
		"contract":   types.ObjectType{AttrTypes: Contract{}.AttributeTypes()},
		"started_at": types.StringType,
	}
}

func newLoadBalancer(
	ctx context.Context,
	sdkLoadBalancer publicCloud.LoadBalancer,
) (*LoadBalancer, diag.Diagnostics) {

	resourcesObject, diags := utils.ConvertSdkModelToResourceObject(
		sdkLoadBalancer.GetResources(),
		Resources{}.AttributeTypes(),
		ctx,
		newResources,
	)
	if diags.HasError() {
		return nil, diags
	}

	contractObject, diags := utils.ConvertSdkModelToResourceObject(
		sdkLoadBalancer.GetContract(),
		Contract{}.AttributeTypes(),
		ctx,
		newContract,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &LoadBalancer{
		Id:        basetypes.NewStringValue(sdkLoadBalancer.GetId()),
		Type:      basetypes.NewStringValue(sdkLoadBalancer.GetType()),
		Resources: resourcesObject,
		Region:    basetypes.NewStringValue(sdkLoadBalancer.GetRegion()),
		Reference: basetypes.NewStringValue(sdkLoadBalancer.GetReference()),
		State:     basetypes.NewStringValue(string(sdkLoadBalancer.GetState())),
		Contract:  contractObject,
		StartedAt: basetypes.NewStringValue(sdkLoadBalancer.GetStartedAt().String()),
	}, nil
}
