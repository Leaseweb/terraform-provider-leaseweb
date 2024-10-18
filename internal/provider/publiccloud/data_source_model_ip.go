package publiccloud

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type DataSourceModelIp struct {
	Ip types.String `tfsdk:"ip"`
}

func newDataSourceModelIp(sdkIp publicCloud.Ip) DataSourceModelIp {
	return DataSourceModelIp{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}
}
