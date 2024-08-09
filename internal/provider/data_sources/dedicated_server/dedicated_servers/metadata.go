package dedicated_servers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *dedicatedServerDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_servers"
}
