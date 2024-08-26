package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Port struct {
	Name types.String `tfsdk:"name"`
	Port types.String `tfsdk:"port"`
}
