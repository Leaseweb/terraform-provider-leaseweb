package public_cloud

import (
	"testing"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/stretchr/testify/assert"
)

func TestNewLoadBalancerConfiguration(t *testing.T) {
	t.Run("required options are set", func(t *testing.T) {
		got := NewLoadBalancerConfiguration(
			enum.BalanceRoundRobin,
			false,
			1,
			OptionalLoadBalancerConfigurationOptions{},
		)

		assert.Equal(t, enum.BalanceRoundRobin, got.Balance)
		assert.False(t, got.XForwardedFor)
		assert.Equal(t, 1, got.IdleTimeout)

		assert.Nil(t, got.StickySession)
	})

	t.Run("optional options are set", func(t *testing.T) {
		got := NewLoadBalancerConfiguration(
			enum.BalanceRoundRobin,
			false,
			1,
			OptionalLoadBalancerConfigurationOptions{
				StickySession: &StickySession{MaxLifeTime: 4},
			},
		)

		assert.Equal(t, 4, got.StickySession.MaxLifeTime)
	})

}
