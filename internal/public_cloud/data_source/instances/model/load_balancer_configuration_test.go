package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func Test_newLoadBalancerConfiguration(t *testing.T) {
	configuration := domain.NewLoadBalancerConfiguration(
		enum.BalanceSource,
		false,
		1,
		2,
		domain.OptionalLoadBalancerConfigurationOptions{
			StickySession: &domain.StickySession{MaxLifeTime: 32},
			HealthCheck:   &domain.HealthCheck{Method: enum.MethodGet},
		},
	)

	got := newLoadBalancerConfiguration(configuration)

	assert.Equal(t, int64(32), got.StickySession.MaxLifeTime.ValueInt64())
	assert.Equal(t, "source", got.Balance.ValueString())
	assert.Equal(t, "GET", got.HealthCheck.Method.ValueString())
	assert.False(t, got.XForwardedFor.ValueBool())
	assert.Equal(t, int64(1), got.IdleTimeout.ValueInt64())
	assert.Equal(t, int64(2), got.TargetPort.ValueInt64())
}
