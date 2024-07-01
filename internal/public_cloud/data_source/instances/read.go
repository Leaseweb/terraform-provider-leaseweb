package instances

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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

	var sdkInstanceDetailsInstances []publicCloud.InstanceDetails
	for _, instance := range instances.Instances {
		sdkInstanceDetails, sdkInstanceDetailsResponse, err := d.client.PublicCloudClient.PublicCloudAPI.GetInstance(
			d.client.AuthContext(ctx),
			instance.GetId(),
		).Execute()
		if err != nil {
			utils.HandleError(
				ctx,
				sdkInstanceDetailsResponse,
				&resp.Diagnostics,
				fmt.Sprintf("Unable to Read Leaseweb Public Cloud Instance %v", instance.GetId()),
				err.Error(),
			)
			return
		}

		sdkInstanceDetailsInstances = append(
			sdkInstanceDetailsInstances,
			*sdkInstanceDetails,
		)

	}

	state := model.Instances{}
	state.Populate(sdkInstanceDetailsInstances)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
