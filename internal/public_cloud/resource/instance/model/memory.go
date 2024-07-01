package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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

func newMemory(sdkMemory *publicCloud.Memory) Memory {
	return Memory{
		Value: basetypes.NewFloat64Value(float64(sdkMemory.GetValue())),
		Unit:  basetypes.NewStringValue(sdkMemory.GetUnit()),
	}
}
