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

func newLoadBalancerConfiguration(
	ctx context.Context,
	entityConfiguration domain.LoadBalancerConfiguration,
) (*LoadBalancerConfiguration, diag.Diagnostics) {
	healthCheckObject, diags := utils.ConvertNullableDomainEntityToResourceObject(
		entityConfiguration.HealthCheck,
		HealthCheck{}.AttributeTypes(),
		ctx,
		newHealthCheck,
	)
	if diags.HasError() {
		return nil, diags
	}

	stickySessionObject, diags := utils.ConvertNullableDomainEntityToResourceObject(
		entityConfiguration.StickySession,
		StickySession{}.AttributeTypes(),
		ctx,
		newStickySession,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &LoadBalancerConfiguration{
		Balance:       basetypes.NewStringValue(entityConfiguration.Balance.String()),
		HealthCheck:   healthCheckObject,
		StickySession: stickySessionObject,
		XForwardedFor: basetypes.NewBoolValue(entityConfiguration.XForwardedFor),
		IdleTimeout:   basetypes.NewInt64Value(int64(entityConfiguration.IdleTimeout)),
		TargetPort:    basetypes.NewInt64Value(int64(entityConfiguration.TargetPort)),
	}, nil
}
