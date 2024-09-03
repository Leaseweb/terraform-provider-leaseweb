package dedicated_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (d *dedicatedServerResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_servers"
}
