package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
