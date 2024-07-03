package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type healthCheck struct {
	Method types.String `tfsdk:"method"`
	Uri    types.String `tfsdk:"uri"`
	Host   types.String `tfsdk:"host"`
	Port   types.Int64  `tfsdk:"port"`
}

func newHealthCheck(sdkHealthCheck publicCloud.HealthCheck) *healthCheck {
	host, hostOk := sdkHealthCheck.GetHostOk()

	return &healthCheck{
		Method: basetypes.NewStringValue(sdkHealthCheck.GetMethod()),
		Uri:    basetypes.NewStringValue(sdkHealthCheck.GetUri()),
		Host:   utils.ConvertNullableSdkStringToStringValue(host, hostOk),
		Port:   basetypes.NewInt64Value(int64(sdkHealthCheck.GetPort())),
	}
}
