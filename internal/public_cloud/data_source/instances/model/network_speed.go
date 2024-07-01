package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type NetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newNetworkSpeed(
	sdkNetworkSpeed publicCloud.NetworkSpeed,
) NetworkSpeed {
	return NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(sdkNetworkSpeed.GetValue())),
		Unit:  basetypes.NewStringValue(sdkNetworkSpeed.GetUnit()),
	}
}
