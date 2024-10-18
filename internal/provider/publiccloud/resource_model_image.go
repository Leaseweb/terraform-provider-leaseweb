package publiccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type ResourceModelImage struct {
	Id types.String `tfsdk:"id"`
}

func (i ResourceModelImage) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": types.StringType,
	}
}

func newResourceModelImage(
	ctx context.Context,
	sdkImage publicCloud.Image,
) (*ResourceModelImage, error) {
	return &ResourceModelImage{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}, nil
}
