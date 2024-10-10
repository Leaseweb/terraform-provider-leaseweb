package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Contract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	State            types.String `tfsdk:"state"`
}

func (c Contract) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_frequency": types.Int64Type,
		"term":              types.Int64Type,
		"type":              types.StringType,
		"state":             types.StringType,
	}
}
