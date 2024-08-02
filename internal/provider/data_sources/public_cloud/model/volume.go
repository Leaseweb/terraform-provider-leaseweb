package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Volume struct {
	Size types.Float64 `tfsdk:"size"`
	Unit types.String  `tfsdk:"unit"`
}
