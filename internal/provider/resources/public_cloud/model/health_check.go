package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type HealthCheck struct {
	Method types.String `tfsdk:"method"`
	Uri    types.String `tfsdk:"uri"`
	Host   types.String `tfsdk:"host"`
	Port   types.Int64  `tfsdk:"port"`
}

func (h HealthCheck) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"method": types.StringType,
		"uri":    types.StringType,
		"host":   types.StringType,
		"port":   types.Int64Type,
	}
}
