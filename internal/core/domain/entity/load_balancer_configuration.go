package entity

import (
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

type LoadBalancerConfiguration struct {
	StickySession *StickySession
	Balance       enum.Balance
	HealthCheck   *HealthCheck
	XForwardedFor bool
	IdleTimeout   int64
	TargetPort    int64
}

type OptionalLoadBalancerConfigurationOptions struct {
	StickySession *StickySession
	HealthCheck   *HealthCheck
}

func NewLoadBalancerConfiguration(
	balance enum.Balance,
	xForwardedFor bool,
	idleTimeout int64,
	targetPort int64,
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
