package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type publicNetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newPublicNetworkSpeed(
	sdkPublicNetworkSpeed publicCloud.PublicNetworkSpeed,
) publicNetworkSpeed {
	return publicNetworkSpeed{
		Value: utils.GenerateInt(
			sdkPublicNetworkSpeed.HasValue(),
			sdkPublicNetworkSpeed.GetValue(),
		),
		Unit: utils.GenerateString(
			sdkPublicNetworkSpeed.HasUnit(),
			sdkPublicNetworkSpeed.GetUnit(),
		),
	}
}
