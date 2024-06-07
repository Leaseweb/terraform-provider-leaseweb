package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
)

type Memory struct {
	Value types.Float64 `tfsdk:"value"`
	Unit  types.String  `tfsdk:"unit"`
}

func (m Memory) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Float64Type,
		"unit":  types.StringType,
	}
}

func newMemory(sdkMemory *publicCloud.Memory) Memory {
	return Memory{
		Value: resources.GetFloatValue(sdkMemory.HasValue(), sdkMemory.GetValue()),
		Unit:  resources.GetStringValue(sdkMemory.HasUnit(), sdkMemory.GetUnit()),
	}
}
