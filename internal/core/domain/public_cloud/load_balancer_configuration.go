package public_cloud

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
)

type LoadBalancerConfiguration struct {
	StickySession *StickySession
	Balance       enum.Balance
	HealthCheck   *HealthCheck
	XForwardedFor bool
	IdleTimeout   int
	TargetPort    int
}

type OptionalLoadBalancerConfigurationOptions struct {
	StickySession *StickySession
	HealthCheck   *HealthCheck
}

func NewLoadBalancerConfiguration(
	balance enum.Balance,
	xForwardedFor bool,
	idleTimeout int,
	targetPort int,
	options OptionalLoadBalancerConfigurationOptions,
) LoadBalancerConfiguration {
	loadBalancerConfiguration := LoadBalancerConfiguration{
		Balance:       balance,
		XForwardedFor: xForwardedFor,
		IdleTimeout:   idleTimeout,
		TargetPort:    targetPort,
	}

	loadBalancerConfiguration.StickySession = options.StickySession
	loadBalancerConfiguration.HealthCheck = options.HealthCheck

	return loadBalancerConfiguration

}
