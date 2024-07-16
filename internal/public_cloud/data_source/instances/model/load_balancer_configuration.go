package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/utils"
)

type loadBalancerConfiguration struct {
	Balance       types.String   `tfsdk:"balance"`
	HealthCheck   *healthCheck   `tfsdk:"health_check"`
	StickySession *stickySession `tfsdk:"sticky_session"`
	XForwardedFor types.Bool     `tfsdk:"x_forwarded_for"`
	IdleTimeout   types.Int64    `tfsdk:"idle_timeout"`
	TargetPort    types.Int64    `tfsdk:"target_port"`
}

func newLoadBalancerConfiguration(entityConfiguration domain.LoadBalancerConfiguration) *loadBalancerConfiguration {

	return &loadBalancerConfiguration{
		Balance: basetypes.NewStringValue(entityConfiguration.Balance.String()),
		HealthCheck: utils.ConvertNullableDomainEntityToDatasourceModel(
			entityConfiguration.HealthCheck,
			newHealthCheck,
		),
		StickySession: utils.ConvertNullableDomainEntityToDatasourceModel(
			entityConfiguration.StickySession,
			newStickySession,
		),
		XForwardedFor: basetypes.NewBoolValue(entityConfiguration.XForwardedFor),
		IdleTimeout:   basetypes.NewInt64Value(int64(entityConfiguration.IdleTimeout)),
		TargetPort:    basetypes.NewInt64Value(int64(entityConfiguration.TargetPort)),
	}
}
