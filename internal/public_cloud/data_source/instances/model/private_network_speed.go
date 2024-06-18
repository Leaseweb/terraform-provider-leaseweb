package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type privateNetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newPrivateNetworkSpeed(
	sdkPrivateNetworkSpeed publicCloud.PrivateNetworkSpeed,
) privateNetworkSpeed {
	return privateNetworkSpeed{
		Value: basetypes.NewInt64Value(int64(sdkPrivateNetworkSpeed.GetValue())),
		Unit:  basetypes.NewStringValue(sdkPrivateNetworkSpeed.GetUnit()),
	}
}
