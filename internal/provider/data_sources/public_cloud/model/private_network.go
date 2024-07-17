package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetwork struct {
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Subnet types.String `tfsdk:"subnet"`
}
