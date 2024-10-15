package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Ip struct {
	Ip types.String `tfsdk:"ip"`
}

func NewIp(sdkIp publicCloud.Ip) Ip {
	return Ip{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}
}
