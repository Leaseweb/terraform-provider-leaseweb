package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/utils"
)

type healthCheck struct {
	Method types.String `tfsdk:"method"`
	Uri    types.String `tfsdk:"uri"`
	Host   types.String `tfsdk:"host"`
	Port   types.Int64  `tfsdk:"port"`
}

func newHealthCheck(entityHealthCheck domain.HealthCheck) *healthCheck {

	return &healthCheck{
		Method: basetypes.NewStringValue(entityHealthCheck.Method.String()),
		Uri:    basetypes.NewStringValue(entityHealthCheck.Uri),
		Host:   utils.ConvertNullableStringToStringValue(entityHealthCheck.Host),
		Port:   basetypes.NewInt64Value(int64(entityHealthCheck.Port)),
	}
}
