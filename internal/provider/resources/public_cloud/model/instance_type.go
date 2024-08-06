package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InstanceType struct {
	Name         types.String `tfsdk:"name"`
	Resources    types.Object `tfsdk:"resources"`
	Prices       types.Object `tfsdk:"prices"`
	StorageTypes types.List   `tfsdk:"storage_types"`
}

func (i InstanceType) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"resources": types.ObjectType{
			AttrTypes: Resources{}.AttributeTypes(),
		},
		"prices": types.ObjectType{
			AttrTypes: Prices{}.AttributeTypes(),
		},
		"storage_types": types.ListType{ElemType: types.StringType},
	}
}
