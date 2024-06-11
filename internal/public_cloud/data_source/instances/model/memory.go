package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type memory struct {
	Value types.Float64 `tfsdk:"value"`
	Unit  types.String  `tfsdk:"unit"`
}

func newMemory(sdkMemory publicCloud.Memory) memory {
	return memory{
		Value: utils.GenerateFloat(sdkMemory.HasValue(), sdkMemory.GetValue()),
		Unit:  utils.GenerateString(sdkMemory.HasUnit(), sdkMemory.GetUnit()),
	}
}
