package dedicatedserver

import (
	"context"

	terraformResource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func (d *dedicatedServerResource) Read(ctx context.Context, req terraformResource.ReadRequest, resp *terraformResource.ReadResponse) {
	var data DedicatedServerModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dedicatedServer, err := d.getServer(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Reading dedicated server", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &dedicatedServer)
	resp.Diagnostics.Append(diags...)
}
