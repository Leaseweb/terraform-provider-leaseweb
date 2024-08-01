package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LoadBalancerConfiguration struct {
	Balance       types.String   `tfsdk:"balance"`
	HealthCheck   *HealthCheck   `tfsdk:"health_check"`
	StickySession *StickySession `tfsdk:"sticky_session"`
	XForwardedFor types.Bool     `tfsdk:"x_forwarded_for"`
	IdleTimeout   types.Int64    `tfsdk:"idle_timeout"`
	TargetPort    types.Int64    `tfsdk:"target_port"`
}
