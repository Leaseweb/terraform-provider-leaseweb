package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newLoadBalancerConfiguration(t *testing.T) {
	sdkLoadBalancerConfiguration := publicCloud.NewLoadBalancerConfiguration(
		*publicCloud.NewNullableStickySession(&publicCloud.StickySession{MaxLifeTime: 5}),
		"balance",
		*publicCloud.NewNullableHealthCheck(&publicCloud.HealthCheck{Method: "method"}),
		false,
		5,
		6,
	)

	got, err := newLoadBalancerConfiguration(
		context.TODO(),
		*sdkLoadBalancerConfiguration,
	)

	assert.Nil(t, err)
	assert.Equal(t, "balance", got.Balance.ValueString())
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
	assert.Equal(t, "method", healthCheck.Method.ValueString())
}

func TestLoadBalancerConfiguration_attributeTypes(t *testing.T) {
	loadBalancerConfiguration, _ := newLoadBalancerConfiguration(
		context.TODO(),
		publicCloud.LoadBalancerConfiguration{},
	)

	_, diags := types.ObjectValueFrom(
		context.TODO(),
		loadBalancerConfiguration.AttributeTypes(),
		loadBalancerConfiguration,
	)

	assert.Nil(t, diags, "attributes should be correct")
}
