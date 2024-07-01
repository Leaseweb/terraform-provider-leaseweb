package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type autoScalingGroup struct {
	Id            types.String  `tfsdk:"id"`
	Type          types.String  `tfsdk:"type"`
	State         types.String  `tfsdk:"state"`
	DesiredAmount types.Int64   `tfsdk:"desired_amount"`
	Region        types.String  `tfsdk:"region"`
	Reference     types.String  `tfsdk:"reference"`
	CreatedAt     types.String  `tfsdk:"created_at"`
	UpdatedAt     types.String  `tfsdk:"updated_at"`
	StartsAt      types.String  `tfsdk:"starts_at"`
	EndsAt        types.String  `tfsdk:"ends_at"`
	MinimumAmount types.Int64   `tfsdk:"minimum_amount"`
	MaximumAmount types.Int64   `tfsdk:"maximum_amount"`
	CpuThreshold  types.Int64   `tfsdk:"cpu_threshold"`
	WarmupTime    types.Int64   `tfsdk:"warmup_time"`
	CooldownTime  types.Int64   `tfsdk:"cooldown_time"`
	LoadBalancer  *loadBalancer `tfsdk:"load_balancer"`
}

func newAutoScalingGroup(sdkAutoScalingGroupDetails publicCloud.AutoScalingGroupDetails) *autoScalingGroup {
	autoScalingLoadBalancer, loadBalancerOk := sdkAutoScalingGroupDetails.GetLoadBalancerOk()

	return &autoScalingGroup{
		Id:            basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetId()),
		Type:          basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetType()),
		State:         basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetState()),
		DesiredAmount: utils.ConvertSdkNullableIntToInt64Value(sdkAutoScalingGroupDetails.GetDesiredAmountOk()),
		Region:        basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetRegion()),
		Reference:     basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetReference()),
		CreatedAt:     basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetCreatedAt().String()),
		UpdatedAt:     basetypes.NewStringValue(sdkAutoScalingGroupDetails.GetUpdatedAt().String()),
		StartsAt:      utils.ConvertSdkNullableTimeToStringValue(sdkAutoScalingGroupDetails.GetStartsAtOk()),
		EndsAt:        utils.ConvertSdkNullableTimeToStringValue(sdkAutoScalingGroupDetails.GetEndsAtOk()),
		MinimumAmount: utils.ConvertSdkNullableIntToInt64Value(sdkAutoScalingGroupDetails.GetMinimumAmountOk()),
		MaximumAmount: utils.ConvertSdkNullableIntToInt64Value(sdkAutoScalingGroupDetails.GetMaximumAmountOk()),
		CpuThreshold:  utils.ConvertSdkNullableIntToInt64Value(sdkAutoScalingGroupDetails.GetCpuThresholdOk()),
		WarmupTime:    utils.ConvertSdkNullableIntToInt64Value(sdkAutoScalingGroupDetails.GetWarmupTimeOk()),
		CooldownTime:  utils.ConvertSdkNullableIntToInt64Value(sdkAutoScalingGroupDetails.GetCooldownTimeOk()),
		LoadBalancer: utils.ConvertNullableSdkModelToDatasourceModel(
			autoScalingLoadBalancer,
			loadBalancerOk,
			newLoadBalancer,
		),
	}
}
