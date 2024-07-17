package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StickySession struct {
	Enabled     types.Bool  `tfsdk:"enabled"`
	MaxLifeTime types.Int64 `tfsdk:"max_life_time"`
}

func (s StickySession) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":       types.BoolType,
		"max_life_time": types.Int64Type,
	}
}
