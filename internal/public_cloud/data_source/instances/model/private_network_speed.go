package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type privateNetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newPrivateNetworkSpeed(
	sdkPrivateNetworkSpeed publicCloud.PrivateNetworkSpeed,
) privateNetworkSpeed {
	return privateNetworkSpeed{
		Value: utils.GenerateInt(
			sdkPrivateNetworkSpeed.HasValue(),
			sdkPrivateNetworkSpeed.GetValue(),
		),
		Unit: utils.GenerateString(
			sdkPrivateNetworkSpeed.HasUnit(),
			sdkPrivateNetworkSpeed.GetUnit(),
		),
	}
}
