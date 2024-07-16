package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type NetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newNetworkSpeed(
	entityNetworkSpeed entity.NetworkSpeed,
) NetworkSpeed {
	return NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(entityNetworkSpeed.Value)),
		Unit:  basetypes.NewStringValue(entityNetworkSpeed.Unit),
	}
}
