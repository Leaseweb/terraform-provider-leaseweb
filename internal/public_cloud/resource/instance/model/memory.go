package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
)

type Memory struct {
	Value types.Float64 `tfsdk:"value"`
	Unit  types.String  `tfsdk:"unit"`
}

func (m Memory) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Float64Type,
		"unit":  types.StringType,
	}
}

func newMemory(entityMemory domain.Memory) Memory {
	return Memory{
		Value: basetypes.NewFloat64Value(entityMemory.Value),
		Unit:  basetypes.NewStringValue(entityMemory.Unit),
	}
}
