package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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
	sdkHealthCheck publicCloud.HealthCheck,
) (*HealthCheck, diag.Diagnostics) {
	return &HealthCheck{
		Method: basetypes.NewStringValue(sdkHealthCheck.GetMethod()),
		Uri:    basetypes.NewStringValue(sdkHealthCheck.GetUri()),
		Host:   basetypes.NewStringValue(sdkHealthCheck.GetHost()),
		Port:   basetypes.NewInt64Value(int64(sdkHealthCheck.GetPort())),
	}, nil
}
