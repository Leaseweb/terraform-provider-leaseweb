package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type Cpu struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (c Cpu) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Int64Type,
		"unit":  types.StringType,
	}
}

func newCpu(sdkCpu *publicCloud.Cpu) Cpu {
	return Cpu{
		Value: utils.GenerateInt(sdkCpu.HasValue(), sdkCpu.GetValue()),
		Unit:  utils.GenerateString(sdkCpu.HasUnit(), sdkCpu.GetUnit()),
	}
}
