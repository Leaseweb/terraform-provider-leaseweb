package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
)

type PrivateNetwork struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}

func newPrivateNetwork(sdkPrivateNetwork publicCloud.PrivateNetwork) PrivateNetwork {
	return PrivateNetwork{
		Id:     resources.GetStringValue(sdkPrivateNetwork.HasPrivateNetworkId(), sdkPrivateNetwork.GetPrivateNetworkId()),
		Status: resources.GetStringValue(sdkPrivateNetwork.HasStatus(), sdkPrivateNetwork.GetStatus()),
		Subnet: resources.GetStringValue(sdkPrivateNetwork.HasSubnet(), sdkPrivateNetwork.GetSubnet()),
	}
}
