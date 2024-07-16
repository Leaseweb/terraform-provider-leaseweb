package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
)

type cpu struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newCpu(entityCpu domain.Cpu) cpu {
	return cpu{
		Value: basetypes.NewInt64Value(int64(entityCpu.Value)),
		Unit:  basetypes.NewStringValue(entityCpu.Unit),
	}
}
