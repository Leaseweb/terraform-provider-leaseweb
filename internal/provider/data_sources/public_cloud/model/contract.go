package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Contract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	State            types.String `tfsdk:"state"`
}
