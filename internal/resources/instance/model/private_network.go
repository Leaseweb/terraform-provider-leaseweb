package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
)

type PrivateNetwork struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}

func (p PrivateNetwork) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":     types.StringType,
		"status": types.StringType,
		"subnet": types.StringType,
	}
}

func newPrivateNetwork(sdkPrivateNetwork publicCloud.PrivateNetwork) PrivateNetwork {
	return PrivateNetwork{
		Id:     resources.GetStringValue(sdkPrivateNetwork.HasPrivateNetworkId(), sdkPrivateNetwork.GetPrivateNetworkId()),
		Status: resources.GetStringValue(sdkPrivateNetwork.HasStatus(), sdkPrivateNetwork.GetStatus()),
		Subnet: resources.GetStringValue(sdkPrivateNetwork.HasSubnet(), sdkPrivateNetwork.GetSubnet()),
	}
}
