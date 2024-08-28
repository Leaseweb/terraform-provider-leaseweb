package control_panels

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (c *controlPanelDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_control_panels"
}
