package operating_systems

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *operatingSystemsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_operating_systems"
}
