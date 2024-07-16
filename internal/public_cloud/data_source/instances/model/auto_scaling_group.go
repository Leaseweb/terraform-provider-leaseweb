package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
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

func newAutoScalingGroup(
	entity domain.AutoScalingGroup,
) *autoScalingGroup {
	return &autoScalingGroup{
		Id:            basetypes.NewStringValue(entity.Id.String()),
		Type:          basetypes.NewStringValue(string(entity.Type)),
		State:         basetypes.NewStringValue(string(entity.State)),
		DesiredAmount: utils.ConvertNullableIntToInt64Value(entity.DesiredAmount),
		Region:        basetypes.NewStringValue(entity.Region),
		Reference:     basetypes.NewStringValue(entity.Reference.String()),
		CreatedAt:     basetypes.NewStringValue(entity.CreatedAt.String()),
		UpdatedAt:     basetypes.NewStringValue(entity.UpdatedAt.String()),
		StartsAt:      utils.ConvertNullableTimeToStringValue(entity.StartsAt),
		EndsAt:        utils.ConvertNullableTimeToStringValue(entity.EndsAt),
		MinimumAmount: utils.ConvertNullableIntToInt64Value(entity.MinimumAmount),
		MaximumAmount: utils.ConvertNullableIntToInt64Value(entity.MaximumAmount),
		CpuThreshold:  utils.ConvertNullableIntToInt64Value(entity.CpuThreshold),
		WarmupTime:    utils.ConvertNullableIntToInt64Value(entity.WarmupTime),
		CooldownTime:  utils.ConvertNullableIntToInt64Value(entity.CooldownTime),
		LoadBalancer: utils.ConvertNullableDomainEntityToDatasourceModel(
			entity.LoadBalancer,
			newLoadBalancer,
		),
	}
}
