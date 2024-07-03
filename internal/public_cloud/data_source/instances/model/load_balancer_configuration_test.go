package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newLoadBalancerConfiguration(t *testing.T) {
	sdkLoadBalancerConfiguration := publicCloud.NewLoadBalancerConfiguration(
		*publicCloud.NewNullableStickySession(&publicCloud.StickySession{MaxLifeTime: 32}),
		"balance",
		*publicCloud.NewNullableHealthCheck(&publicCloud.HealthCheck{Method: "method"}),
		false,
		1,
		2,
	)

	got := newLoadBalancerConfiguration(*sdkLoadBalancerConfiguration)

	assert.Equal(t, int64(32), got.StickySession.MaxLifeTime.ValueInt64())
	assert.Equal(t, "balance", got.Balance.ValueString())
	assert.Equal(t, "method", got.HealthCheck.Method.ValueString())
	assert.False(t, got.XForwardedFor.ValueBool())
	assert.Equal(t, int64(1), got.IdleTimeout.ValueInt64())
	assert.Equal(t, int64(2), got.TargetPort.ValueInt64())
}
