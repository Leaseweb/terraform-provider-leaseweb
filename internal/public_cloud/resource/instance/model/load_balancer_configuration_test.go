package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
)

func Test_newLoadBalancerConfiguration(t *testing.T) {
	loadBalancerConfiguration := domain.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		5,
		6,
		domain.OptionalLoadBalancerConfigurationOptions{
			StickySession: &domain.StickySession{MaxLifeTime: 5},
			HealthCheck:   &domain.HealthCheck{Method: enum.MethodHead},
		},
	)

	got, err := newLoadBalancerConfiguration(
		context.TODO(),
		loadBalancerConfiguration,
	)

	assert.Nil(t, err)
	assert.Equal(t, "source", got.Balance.ValueString())
	assert.False(t, got.XForwardedFor.ValueBool())
	assert.Equal(t, int64(5), got.IdleTimeout.ValueInt64())
	assert.Equal(t, int64(6), got.TargetPort.ValueInt64())

	stickySession := StickySession{}
	got.StickySession.As(
		context.TODO(),
		&stickySession,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, int64(5), stickySession.MaxLifeTime.ValueInt64())

	healthCheck := HealthCheck{}
	got.HealthCheck.As(
		context.TODO(),
		&healthCheck,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, "HEAD", healthCheck.Method.ValueString())
}

func TestLoadBalancerConfiguration_attributeTypes(t *testing.T) {
	loadBalancerConfiguration, _ := newLoadBalancerConfiguration(
		context.TODO(),
		domain.LoadBalancerConfiguration{},
	)

	_, diags := types.ObjectValueFrom(
		context.TODO(),
		loadBalancerConfiguration.AttributeTypes(),
		loadBalancerConfiguration,
	)

	assert.Nil(t, diags, "attributes should be correct")
}
