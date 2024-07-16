package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
)

type NetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (p NetworkSpeed) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Int64Type,
		"unit":  types.StringType,
	}
}

func newNetworkSpeed(entityNetworkSpeed domain.NetworkSpeed) NetworkSpeed {
	return NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(entityNetworkSpeed.Value)),
		Unit:  basetypes.NewStringValue(entityNetworkSpeed.Unit),
	}
}
