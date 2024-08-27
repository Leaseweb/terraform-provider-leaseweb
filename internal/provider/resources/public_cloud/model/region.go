package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Region struct {
	Name     types.String `tfsdk:"name"`
	Location types.String `tfsdk:"location"`
}

func (r Region) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":     types.StringType,
		"location": types.StringType,
	}
}
