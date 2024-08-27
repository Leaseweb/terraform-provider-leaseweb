package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type Ram struct {
	Size types.Int32  `tfsdk:"size"`
	Unit types.String `tfsdk:"unit"`
}
