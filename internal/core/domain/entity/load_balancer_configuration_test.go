package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func TestNewLoadBalancerConfiguration(t *testing.T) {
	t.Run("required options are set", func(t *testing.T) {
		got := NewLoadBalancerConfiguration(
			enum.BalanceRoundRobin,
			false,
			1,
			2,
			OptionalLoadBalancerConfigurationOptions{},
		)

		assert.Equal(t, enum.BalanceRoundRobin, got.Balance)
		assert.False(t, got.XForwardedFor)
		assert.Equal(t, 1, got.IdleTimeout)
		assert.Equal(t, 2, got.TargetPort)

		assert.Nil(t, got.StickySession)
		assert.Nil(t, got.HealthCheck)
	})

	t.Run("optional options are set", func(t *testing.T) {
		got := NewLoadBalancerConfiguration(
			enum.BalanceRoundRobin,
			false,
			1,
			2,
			OptionalLoadBalancerConfigurationOptions{
				StickySession: &StickySession{MaxLifeTime: 4},
				HealthCheck:   &HealthCheck{Uri: "uri"},
			},
		)

		assert.Equal(t, 4, got.StickySession.MaxLifeTime)
		assert.Equal(t, "uri", got.HealthCheck.Uri)
	})

}
