package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageSize struct {
	Size types.Float64 `tfsdk:"size"`
	Unit types.String  `tfsdk:"unit"`
}
