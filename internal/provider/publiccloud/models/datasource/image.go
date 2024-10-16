package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Image struct {
	Id types.String `tfsdk:"id"`
}

func newImage(sdkImage publicCloud.Image) Image {
	return Image{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}
}
