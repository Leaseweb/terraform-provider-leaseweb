package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}
