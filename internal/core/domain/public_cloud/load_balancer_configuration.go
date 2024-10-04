package public_cloud

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
)

type LoadBalancerConfiguration struct {
	StickySession *StickySession
	Balance       enum.Balance
	XForwardedFor bool
	IdleTimeout   int
}

type OptionalLoadBalancerConfigurationOptions struct {
	StickySession *StickySession
}

func NewLoadBalancerConfiguration(
	balance enum.Balance,
	xForwardedFor bool,
	idleTimeout int,
	options OptionalLoadBalancerConfigurationOptions,
) LoadBalancerConfiguration {
	loadBalancerConfiguration := LoadBalancerConfiguration{
		Balance:       balance,
		XForwardedFor: xForwardedFor,
		IdleTimeout:   idleTimeout,
	}

	loadBalancerConfiguration.StickySession = options.StickySession

	return loadBalancerConfiguration

}
