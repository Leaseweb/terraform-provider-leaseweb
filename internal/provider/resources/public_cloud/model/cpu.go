package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Cpu struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (c Cpu) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Int64Type,
		"unit":  types.StringType,
	}
}
