package instances

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	instances, err := d.client.PublicCloudHandler.GetAllInstances(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read instances", err.Error())
		return
	}

	state := instances

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
