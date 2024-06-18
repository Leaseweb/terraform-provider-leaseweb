package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type publicNetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newPublicNetworkSpeed(
	sdkPublicNetworkSpeed publicCloud.PublicNetworkSpeed,
) publicNetworkSpeed {
	return publicNetworkSpeed{
		Value: basetypes.NewInt64Value(int64(sdkPublicNetworkSpeed.GetValue())),
		Unit:  basetypes.NewStringValue(sdkPublicNetworkSpeed.GetUnit()),
	}
}
