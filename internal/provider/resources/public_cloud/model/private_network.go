package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetwork struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}

func (p PrivateNetwork) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":     types.StringType,
		"status": types.StringType,
		"subnet": types.StringType,
	}
}
