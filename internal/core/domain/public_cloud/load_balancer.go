package public_cloud

import (
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
)

type LoadBalancer struct {
	Id             string
	Type           InstanceType
	Resources      Resources
	Region         Region
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
	id string,
	loadBalancerType InstanceType,
	resources Resources,
	region Region,
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
