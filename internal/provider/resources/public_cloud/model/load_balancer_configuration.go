package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LoadBalancerConfiguration struct {
	Balance       types.String `tfsdk:"balance"`
	HealthCheck   types.Object `tfsdk:"health_check"`
	StickySession types.Object `tfsdk:"sticky_session"`
	XForwardedFor types.Bool   `tfsdk:"x_forwarded_for"`
	IdleTimeout   types.Int64  `tfsdk:"idle_timeout"`
	TargetPort    types.Int64  `tfsdk:"target_port"`
}

func (l LoadBalancerConfiguration) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"balance":         types.StringType,
		"health_check":    types.ObjectType{AttrTypes: HealthCheck{}.AttributeTypes()},
		"sticky_session":  types.ObjectType{AttrTypes: StickySession{}.AttributeTypes()},
		"x_forwarded_for": types.BoolType,
		"idle_timeout":    types.Int64Type,
		"target_port":     types.Int64Type,
	}
}
