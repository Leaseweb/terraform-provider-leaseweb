package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type Cpu struct {
	Quantity types.Int32  `tfsdk:"quantity"`
	Type     types.String `tfsdk:"type"`
}
