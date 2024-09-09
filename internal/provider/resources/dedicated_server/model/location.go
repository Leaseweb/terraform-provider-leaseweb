package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Location struct {
	Rack  types.String `tfsdk:"rack"`
	Site  types.String `tfsdk:"site"`
	Suite types.String `tfsdk:"suite"`
	Unit  types.String `tfsdk:"unit"`
}
