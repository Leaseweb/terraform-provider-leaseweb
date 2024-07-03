package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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

func newLoadBalancer(sdkLoadBalancer publicCloud.LoadBalancerDetails) *loadBalancer {
	loadBalancerConfiguration, loadBalancerConfigurationOk := sdkLoadBalancer.GetConfigurationOk()
	loadBalancerPrivateNetwork, loadBalancerPrivateNetworkOk := sdkLoadBalancer.GetPrivateNetworkOk()

	var ips []ip
	for _, ip := range sdkLoadBalancer.Ips {
		ips = append(ips, newIp(ip))
	}

	return &loadBalancer{
		Id:        basetypes.NewStringValue(sdkLoadBalancer.GetId()),
		Type:      basetypes.NewStringValue(sdkLoadBalancer.GetType()),
		Resources: newResources(sdkLoadBalancer.GetResources()),
		Region:    basetypes.NewStringValue(sdkLoadBalancer.GetRegion()),
		Reference: basetypes.NewStringValue(sdkLoadBalancer.GetReference()),
		State:     basetypes.NewStringValue(string(sdkLoadBalancer.GetState())),
		Contract:  newContract(sdkLoadBalancer.GetContract()),
		StartedAt: basetypes.NewStringValue(sdkLoadBalancer.GetStartedAt().String()),
		Ips:       ips,
		LoadBalancerConfiguration: utils.ConvertNullableSdkModelToDatasourceModel(
			loadBalancerConfiguration,
			loadBalancerConfigurationOk,
			newLoadBalancerConfiguration,
		),
		PrivateNetwork: utils.ConvertNullableSdkModelToDatasourceModel(
			loadBalancerPrivateNetwork,
			loadBalancerPrivateNetworkOk,
			newPrivateNetwork,
		),
	}
}
