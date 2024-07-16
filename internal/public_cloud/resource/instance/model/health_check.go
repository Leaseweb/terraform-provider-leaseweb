package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/utils"
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

func newHealthCheck(
	ctx context.Context,
	entityHealthCheck entity.HealthCheck,
) (*HealthCheck, diag.Diagnostics) {
	return &HealthCheck{
		Method: basetypes.NewStringValue(string(entityHealthCheck.Method)),
		Uri:    basetypes.NewStringValue(entityHealthCheck.Uri),
		Host:   utils.ConvertNullableStringToStringValue(entityHealthCheck.Host),
		Port:   basetypes.NewInt64Value(int64(entityHealthCheck.Port)),
	}, nil
}
