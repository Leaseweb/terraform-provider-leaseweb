package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (i *instanceResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_instance"
}
