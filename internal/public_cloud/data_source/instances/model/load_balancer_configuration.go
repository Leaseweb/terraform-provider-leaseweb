package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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

func newLoadBalancerConfiguration(sdkLoadBalancerConfiguration publicCloud.LoadBalancerConfiguration) *loadBalancerConfiguration {
	healthCheck, healthCheckOk := sdkLoadBalancerConfiguration.GetHealthCheckOk()
	stickySession, stickySessionOk := sdkLoadBalancerConfiguration.GetStickySessionOk()

	return &loadBalancerConfiguration{
		Balance: basetypes.NewStringValue(sdkLoadBalancerConfiguration.GetBalance()),
		HealthCheck: utils.ConvertNullableSdkModelToDatasourceModel(
			healthCheck,
			healthCheckOk,
			newHealthCheck,
		),
		StickySession: utils.ConvertNullableSdkModelToDatasourceModel(
			stickySession,
			stickySessionOk,
			newStickySession,
		),
		XForwardedFor: basetypes.NewBoolValue(sdkLoadBalancerConfiguration.GetXForwardedFor()),
		IdleTimeout:   basetypes.NewInt64Value(int64(sdkLoadBalancerConfiguration.GetIdleTimeOut())),
		TargetPort:    basetypes.NewInt64Value(int64(sdkLoadBalancerConfiguration.GetTargetPort())),
	}
}
