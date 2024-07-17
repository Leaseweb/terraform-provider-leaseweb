package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LoadBalancer struct {
	Id                        types.String               `tfsdk:"id"`
	Type                      types.String               `tfsdk:"type"`
	Resources                 Resources                  `tfsdk:"resources"`
	Region                    types.String               `tfsdk:"region"`
	Reference                 types.String               `tfsdk:"reference"`
	State                     types.String               `tfsdk:"state"`
	Contract                  Contract                   `tfsdk:"contract"`
	StartedAt                 types.String               `tfsdk:"started_at"`
	Ips                       []Ip                       `tfsdk:"ips"`
	LoadBalancerConfiguration *LoadBalancerConfiguration `tfsdk:"load_balancer_configuration"`
	PrivateNetwork            *PrivateNetwork            `tfsdk:"private_network"`
}
