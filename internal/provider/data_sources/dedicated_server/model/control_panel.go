package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ControlPanel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
