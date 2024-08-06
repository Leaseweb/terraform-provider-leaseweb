package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Price struct {
	HourlyPrice  types.String `tfsdk:"hourly_price"`
	MonthlyPrice types.String `tfsdk:"monthly_price"`
}

func (p Price) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hourly_price":  types.StringType,
		"monthly_price": types.StringType,
	}
}
