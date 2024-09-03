package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type PciCard struct {
	Description types.String `tfsdk:"description"`
}
