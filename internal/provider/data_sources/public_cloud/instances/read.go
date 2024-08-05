package instances

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
)

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	tflog.Info(ctx, "Read public cloud instances")
	instances, err := d.client.PublicCloudFacade.GetAllInstances(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read instances", err.Error())
		logging.HandleError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Unable to read instances",
			err.Error(),
		)

		return
	}

	state := instances

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
