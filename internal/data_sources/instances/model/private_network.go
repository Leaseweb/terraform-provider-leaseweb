package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type privateNetwork struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}

func newPrivateNetwork(sdkPrivateNetwork publicCloud.PrivateNetwork) privateNetwork {
	return privateNetwork{
		Id: utils.GenerateString(
			sdkPrivateNetwork.HasPrivateNetworkId(),
			sdkPrivateNetwork.GetPrivateNetworkId(),
		),
		Status: utils.GenerateString(
			sdkPrivateNetwork.HasStatus(),
			sdkPrivateNetwork.GetStatus(),
		),
		Subnet: utils.GenerateString(
			sdkPrivateNetwork.HasSubnet(),
			sdkPrivateNetwork.GetSubnet(),
		),
	}
}
