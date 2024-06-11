package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type cpu struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func newCpu(sdkCpu publicCloud.Cpu) cpu {
	return cpu{
		Value: utils.GenerateInt(sdkCpu.HasValue(), sdkCpu.GetValue()),
		Unit:  utils.GenerateString(sdkCpu.HasUnit(), sdkCpu.GetUnit()),
	}
}
