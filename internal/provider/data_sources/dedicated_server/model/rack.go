package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Rack struct {
	Id       types.String `tfsdk:"id"`
	Capacity types.String `tfsdk:"capacity"`
	Type     types.String `tfsdk:"type"`
}
