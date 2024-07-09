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

type AutoScalingGroup struct {
	Id            types.String `tfsdk:"id"`
	Type          types.String `tfsdk:"type"`
	State         types.String `tfsdk:"state"`
	DesiredAmount types.Int64  `tfsdk:"desired_amount"`
	Region        types.String `tfsdk:"region"`
	Reference     types.String `tfsdk:"reference"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	StartsAt      types.String `tfsdk:"starts_at"`
	EndsAt        types.String `tfsdk:"ends_at"`
	MinimumAmount types.Int64  `tfsdk:"minimum_amount"`
	MaximumAmount types.Int64  `tfsdk:"maximum_amount"`
	CpuThreshold  types.Int64  `tfsdk:"cpu_threshold"`
	WarmupTime    types.Int64  `tfsdk:"warmup_time"`
	CooldownTime  types.Int64  `tfsdk:"cooldown_time"`
	LoadBalancer  types.Object `tfsdk:"load_balancer"`
}

func (a AutoScalingGroup) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"type":           types.StringType,
		"state":          types.StringType,
		"desired_amount": types.Int64Type,
		"region":         types.StringType,
		"reference":      types.StringType,
		"created_at":     types.StringType,
		"updated_at":     types.StringType,
		"starts_at":      types.StringType,
		"ends_at":        types.StringType,
		"minimum_amount": types.Int64Type,
		"maximum_amount": types.Int64Type,
		"cpu_threshold":  types.Int64Type,
		"warmup_time":    types.Int64Type,
		"cooldown_time":  types.Int64Type,
		"load_balancer":  types.ObjectType{AttrTypes: LoadBalancer{}.AttributeTypes()},
	}
}

func newAutoScalingGroup(
	ctx context.Context,
	sdkAutoScalingGroupDetails publicCloud.AutoScalingGroupDetails,
	sdkLoadBalancerDetails *publicCloud.LoadBalancerDetails,
) (*AutoScalingGroup, diag.Diagnostics) {

	autoScalingLoadBalancerObject, diags := utils.ConvertNullableSdkModelToResourceObject(
		sdkLoadBalancerDetails,
		true,
		LoadBalancer{}.AttributeTypes(),
		ctx,
		newLoadBalancer,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &AutoScalingGroup{
		Id:            basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetId()),
		Type:          basetypes.NewStringValue(string(sdkAutoScalingGroupDetails.GetType())),
		State:         basetypes.NewStringValue(string(sdkAutoScalingGroupDetails.GetState())),
		DesiredAmount: utils.ConvertNullableSdkIntToInt64Value(sdkAutoScalingGroupDetails.GetDesiredAmountOk()),
		Region:        basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetRegion()),
		Reference:     basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetReference()),
		CreatedAt:     basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetCreatedAt().String()),
		UpdatedAt:     basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetUpdatedAt().String()),
		StartsAt:      utils.ConvertNullableSdkTimeToStringValue(sdkAutoScalingGroupDetails.GetStartsAtOk()),
		EndsAt:        utils.ConvertNullableSdkTimeToStringValue(sdkAutoScalingGroupDetails.GetEndsAtOk()),
		MinimumAmount: utils.ConvertNullableSdkIntToInt64Value(sdkAutoScalingGroupDetails.GetMinimumAmountOk()),
		MaximumAmount: utils.ConvertNullableSdkIntToInt64Value(sdkAutoScalingGroupDetails.GetMaximumAmountOk()),
		CpuThreshold:  utils.ConvertNullableSdkIntToInt64Value(sdkAutoScalingGroupDetails.GetCpuThresholdOk()),
		WarmupTime:    utils.ConvertNullableSdkIntToInt64Value(sdkAutoScalingGroupDetails.GetWarmupTimeOk()),
		CooldownTime:  utils.ConvertNullableSdkIntToInt64Value(sdkAutoScalingGroupDetails.GetCooldownTimeOk()),
		LoadBalancer:  autoScalingLoadBalancerObject,
	}, nil
}
