package publiccloud

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type DataSourceModelImage struct {
	Id types.String `tfsdk:"id"`
}

func newDataSourceModelImage(sdkImage publicCloud.Image) DataSourceModelImage {
	return DataSourceModelImage{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}
}
