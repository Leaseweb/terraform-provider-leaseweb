package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/utils"
)

type loadBalancer struct {
	Id                        types.String               `tfsdk:"id"`
	Type                      types.String               `tfsdk:"type"`
	Resources                 resources                  `tfsdk:"resources"`
	Region                    types.String               `tfsdk:"region"`
	Reference                 types.String               `tfsdk:"reference"`
	State                     types.String               `tfsdk:"state"`
	Contract                  contract                   `tfsdk:"contract"`
	StartedAt                 types.String               `tfsdk:"started_at"`
	Ips                       []ip                       `tfsdk:"ips"`
	LoadBalancerConfiguration *loadBalancerConfiguration `tfsdk:"load_balancer_configuration"`
	PrivateNetwork            *privateNetwork            `tfsdk:"private_network"`
}

func newLoadBalancer(entityLoadBalancer entity.LoadBalancer) *loadBalancer {

	var ips []ip
	for _, ip := range entityLoadBalancer.Ips {
		ips = append(ips, newIp(ip))
	}

	return &loadBalancer{
		Id:        basetypes.NewStringValue(entityLoadBalancer.Id.String()),
		Type:      basetypes.NewStringValue(entityLoadBalancer.Type),
		Resources: newResources(entityLoadBalancer.Resources),
		Region:    basetypes.NewStringValue(entityLoadBalancer.Region),
		Reference: utils.ConvertNullableStringToStringValue(entityLoadBalancer.Reference),
		State:     basetypes.NewStringValue(string(entityLoadBalancer.State)),
		Contract:  newContract(entityLoadBalancer.Contract),
		StartedAt: utils.ConvertNullableTimeToStringValue(entityLoadBalancer.StartedAt),
		Ips:       ips,
		LoadBalancerConfiguration: utils.ConvertNullableDomainEntityToDatasourceModel(
			entityLoadBalancer.Configuration,
			newLoadBalancerConfiguration,
		),
		PrivateNetwork: utils.ConvertNullableDomainEntityToDatasourceModel(
			entityLoadBalancer.PrivateNetwork,
			newPrivateNetwork,
		),
	}
}
