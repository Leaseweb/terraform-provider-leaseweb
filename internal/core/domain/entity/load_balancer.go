package entity

import (
	"time"

	"github.com/google/uuid"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

type LoadBalancer struct {
	Id             uuid.UUID
	Type           string
	Resources      Resources
	Region         string
	State          enum.State
	Contract       Contract
	Reference      *string
	StartedAt      *time.Time
	Ips            Ips
	Configuration  *LoadBalancerConfiguration
	PrivateNetwork *PrivateNetwork
}

type OptionalLoadBalancerValues struct {
	Reference      *string
	StartedAt      *time.Time
	PrivateNetwork *PrivateNetwork
	Configuration  *LoadBalancerConfiguration
}

func NewLoadBalancer(
	id uuid.UUID,
	loadBalancerType string,
	resources Resources,
	region string,
	state enum.State,
	contract Contract,
	ips Ips,
	options OptionalLoadBalancerValues,
) LoadBalancer {
	loadBalancer := LoadBalancer{
		Id:        id,
		Type:      loadBalancerType,
		Resources: resources,
		Region:    region,
		State:     state,
		Contract:  contract,
		Ips:       ips,
	}

	loadBalancer.Reference = options.Reference
	loadBalancer.StartedAt = options.StartedAt
	loadBalancer.PrivateNetwork = options.PrivateNetwork
	loadBalancer.Configuration = options.Configuration

	return loadBalancer

}
