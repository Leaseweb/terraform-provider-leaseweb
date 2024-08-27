package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Region struct {
	Name     types.String `tfsdk:"name"`
	Location types.String `tfsdk:"location"`
}
