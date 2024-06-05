package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *instanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}
