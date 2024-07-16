package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
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
	entityAutoScalingGroup domain.AutoScalingGroup,
) (*AutoScalingGroup, diag.Diagnostics) {
	autoScalingLoadBalancerObject, diags := utils.ConvertNullableDomainEntityToResourceObject(
		entityAutoScalingGroup.LoadBalancer,
		LoadBalancer{}.AttributeTypes(),
		ctx,
		newLoadBalancer,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &AutoScalingGroup{
		Id:            basetypes.NewStringValue(entityAutoScalingGroup.Id.String()),
		Type:          basetypes.NewStringValue(string(entityAutoScalingGroup.Type)),
		State:         basetypes.NewStringValue(string(entityAutoScalingGroup.State)),
		DesiredAmount: utils.ConvertNullableIntToInt64Value(entityAutoScalingGroup.DesiredAmount),
		Region:        basetypes.NewStringValue(entityAutoScalingGroup.Region),
		Reference:     basetypes.NewStringValue(entityAutoScalingGroup.Reference.String()),
		CreatedAt:     basetypes.NewStringValue(entityAutoScalingGroup.CreatedAt.String()),
		UpdatedAt:     basetypes.NewStringValue(entityAutoScalingGroup.UpdatedAt.String()),
		StartsAt:      utils.ConvertNullableTimeToStringValue(entityAutoScalingGroup.StartsAt),
		EndsAt:        utils.ConvertNullableTimeToStringValue(entityAutoScalingGroup.EndsAt),
		MinimumAmount: utils.ConvertNullableIntToInt64Value(entityAutoScalingGroup.MinimumAmount),
		MaximumAmount: utils.ConvertNullableIntToInt64Value(entityAutoScalingGroup.MaximumAmount),
		CpuThreshold:  utils.ConvertNullableIntToInt64Value(entityAutoScalingGroup.CpuThreshold),
		WarmupTime:    utils.ConvertNullableIntToInt64Value(entityAutoScalingGroup.WarmupTime),
		CooldownTime:  utils.ConvertNullableIntToInt64Value(entityAutoScalingGroup.CooldownTime),
		LoadBalancer:  autoScalingLoadBalancerObject,
	}, nil
}
