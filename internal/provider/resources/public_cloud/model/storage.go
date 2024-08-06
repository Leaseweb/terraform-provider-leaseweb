package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Storage struct {
	Local   types.Object `tfsdk:"local"`
	Central types.Object `tfsdk:"central"`
}

func (s Storage) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"local": types.ObjectType{
			AttrTypes: Price{}.AttributeTypes(),
		},
		"central": types.ObjectType{
			AttrTypes: Price{}.AttributeTypes(),
		},
	}
}
