package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Price struct {
	HourlyPrice  types.String `tfsdk:"hourly_price"`
	MonthlyPrice types.String `tfsdk:"monthly_price"`
}
