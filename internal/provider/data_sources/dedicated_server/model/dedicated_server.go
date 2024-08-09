package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DedicatedServer struct {
	Id types.String `tfsdk:"id"`
}
