package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StickySession struct {
	Enabled     types.Bool  `tfsdk:"enabled"`
	MaxLifeTime types.Int64 `tfsdk:"max_life_time"`
}
