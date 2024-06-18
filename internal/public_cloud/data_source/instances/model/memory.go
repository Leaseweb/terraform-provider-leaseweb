package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type memory struct {
	Value types.Float64 `tfsdk:"value"`
	Unit  types.String  `tfsdk:"unit"`
}

func newMemory(sdkMemory publicCloud.Memory) memory {
	return memory{
		Value: basetypes.NewFloat64Value(float64(sdkMemory.GetValue())),
		Unit:  basetypes.NewStringValue(sdkMemory.GetUnit()),
	}
}
