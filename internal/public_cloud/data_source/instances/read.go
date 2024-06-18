package instances

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"terraform-provider-leaseweb/internal/public_cloud/data_source/instances/model"
	"terraform-provider-leaseweb/internal/utils"
)

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	instances, sdkResponse, err := d.client.PublicCloudClient.PublicCloudAPI.
		GetInstanceList(d.client.AuthContext(ctx)).Execute()
	if err != nil {
		utils.HandleError(
			ctx,
			sdkResponse,
			&resp.Diagnostics,
			"Unable to Read Leaseweb Public Cloud Instances",
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
