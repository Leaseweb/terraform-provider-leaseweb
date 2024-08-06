package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Prices struct {
	Currency       types.String `tfsdk:"currency"`
	CurrencySymbol types.String `tfsdk:"currency_symbol"`
	Compute        Price        `tfsdk:"compute"`
	Storage        Storage      `tfsdk:"storage"`
}
