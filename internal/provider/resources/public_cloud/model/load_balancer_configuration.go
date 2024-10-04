package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LoadBalancerConfiguration struct {
	Balance       types.String `tfsdk:"balance"`
	StickySession types.Object `tfsdk:"sticky_session"`
	XForwardedFor types.Bool   `tfsdk:"x_forwarded_for"`
	IdleTimeout   types.Int64  `tfsdk:"idle_timeout"`
}

func (l LoadBalancerConfiguration) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"balance":         types.StringType,
		"sticky_session":  types.ObjectType{AttrTypes: StickySession{}.AttributeTypes()},
		"x_forwarded_for": types.BoolType,
		"idle_timeout":    types.Int64Type,
	}
}
