package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AutoScalingGroup struct {
	Id            types.String  `tfsdk:"id"`
	Type          types.String  `tfsdk:"type"`
	State         types.String  `tfsdk:"state"`
	DesiredAmount types.Int64   `tfsdk:"desired_amount"`
	Region        Region        `tfsdk:"region"`
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
	LoadBalancer  *LoadBalancer `tfsdk:"load_balancer"`
}
