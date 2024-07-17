package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Iso struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (i Iso) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}
