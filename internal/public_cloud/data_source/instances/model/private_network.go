package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type privateNetwork struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}

func newPrivateNetwork(sdkPrivateNetwork publicCloud.PrivateNetwork) *privateNetwork {
	return &privateNetwork{
		Id:     basetypes.NewStringValue(sdkPrivateNetwork.GetPrivateNetworkId()),
		Status: basetypes.NewStringValue(sdkPrivateNetwork.GetStatus()),
		Subnet: basetypes.NewStringValue(sdkPrivateNetwork.GetSubnet()),
	}
}
