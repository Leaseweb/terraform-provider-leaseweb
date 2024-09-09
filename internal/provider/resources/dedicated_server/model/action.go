package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Action struct {
	LastTriggeredAt types.String `tfsdk:"last_triggered_at"`
	Type            types.String `tfsdk:"type"`
}

func (a Action) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"last_triggered_at": types.StringType,
		"type":              types.StringType,
	}
}
