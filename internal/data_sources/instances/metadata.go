package instances

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *instancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instances"
}
