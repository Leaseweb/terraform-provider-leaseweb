package instances

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"terraform-provider-leaseweb/internal/public_cloud/data_source/instances/model"
)

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	instances, err := d.client.PublicCloud.GetAllInstances(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read instances", err.Error())
		return
	}

	state := model.Instances{}
	state.Populate(instances)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
