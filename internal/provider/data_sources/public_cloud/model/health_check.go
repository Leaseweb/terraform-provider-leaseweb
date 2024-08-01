package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type HealthCheck struct {
	Method types.String `tfsdk:"method"`
	Uri    types.String `tfsdk:"uri"`
	Host   types.String `tfsdk:"host"`
	Port   types.Int64  `tfsdk:"port"`
}
