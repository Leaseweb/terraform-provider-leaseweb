package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type cpu struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newCpu(sdkCpu publicCloud.Cpu) cpu {
	return cpu{
		Value: basetypes.NewInt64Value(int64(sdkCpu.GetValue())),
		Unit:  basetypes.NewStringValue(sdkCpu.GetUnit()),
	}
}
