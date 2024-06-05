package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
)

type PrivateNetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newPrivateNetworkSpeed(sdkPrivateNetworkSpeed *publicCloud.PrivateNetworkSpeed) PrivateNetworkSpeed {
	return PrivateNetworkSpeed{
		Value: resources.GetIntValue(sdkPrivateNetworkSpeed.HasValue(), sdkPrivateNetworkSpeed.GetValue()),
		Unit:  resources.GetStringValue(sdkPrivateNetworkSpeed.HasUnit(), sdkPrivateNetworkSpeed.GetUnit()),
	}
}
