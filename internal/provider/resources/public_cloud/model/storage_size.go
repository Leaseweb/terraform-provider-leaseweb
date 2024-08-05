package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageSize struct {
	Size types.Float64 `tfsdk:"size"`
	Unit types.String  `tfsdk:"unit"`
}

func (s StorageSize) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"size": types.Float64Type,
		"unit": types.StringType,
	}
}
