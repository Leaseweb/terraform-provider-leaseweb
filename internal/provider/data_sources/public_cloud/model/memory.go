package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Memory struct {
	Value types.Float64 `tfsdk:"value"`
	Unit  types.String  `tfsdk:"unit"`
}
