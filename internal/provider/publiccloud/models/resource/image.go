package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Image struct {
	Id types.String `tfsdk:"id"`
}

func (i Image) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": types.StringType,
	}
}

func newImage(
	ctx context.Context,
	sdkImage publicCloud.Image,
) (*Image, error) {
	return &Image{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}, nil
}
