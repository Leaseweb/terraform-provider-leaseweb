package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
)

type memory struct {
	Value types.Float64 `tfsdk:"value"`
	Unit  types.String  `tfsdk:"unit"`
}

func newMemory(entityMemory domain.Memory) memory {
	return memory{
		Value: basetypes.NewFloat64Value(entityMemory.Value),
		Unit:  basetypes.NewStringValue(entityMemory.Unit),
	}
}
