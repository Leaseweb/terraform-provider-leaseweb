package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Image struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Family  types.String `tfsdk:"family"`
	Flavour types.String `tfsdk:"flavour"`
	Custom  types.Bool   `tfsdk:"custom"`
}

func (i Image) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":      types.StringType,
		"name":    types.StringType,
		"family":  types.StringType,
		"flavour": types.StringType,
		"custom":  types.BoolType,
	}
}
