package instances

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"terraform-provider-leaseweb/internal/data_sources/instances/model"
)

func (d *instancesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {

	instances, _, err := d.client.SdkClient.PublicCloudAPI.GetInstanceList(d.client.AuthContext()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Leaseweb Instances",
			err.Error(),
		)
		return
	}

	state := model.Instances{}
	state.Populate(instances.Instances)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
