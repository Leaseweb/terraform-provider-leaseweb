package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Prices struct {
	Currency       types.String `tfsdk:"currency"`
	CurrencySymbol types.String `tfsdk:"currency_symbol"`
	Compute        types.Object `tfsdk:"compute"`
	Storage        types.Object `tfsdk:"storage"`
}

func (p Prices) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"currency":        types.StringType,
		"currency_symbol": types.StringType,
		"compute": types.ObjectType{
			AttrTypes: Price{}.AttributeTypes(),
		},
		"storage": types.ObjectType{
			AttrTypes: Storage{}.AttributeTypes(),
		},
	}
}
