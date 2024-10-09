package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Image struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Family  types.String `tfsdk:"family"`
	Flavour types.String `tfsdk:"flavour"`
	Custom  types.Bool   `tfsdk:"custom"`
}
