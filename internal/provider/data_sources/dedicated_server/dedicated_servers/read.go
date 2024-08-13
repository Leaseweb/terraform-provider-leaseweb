package dedicated_servers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
)

func (d *dedicatedServerDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	tflog.Info(ctx, "Read dedicated servers")
	dedicatedServers, err := d.client.DedicatedServerFacade.GetAllDedicatedServers(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read dedicated servers", err.Error())
		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Unable to read dedicated servers",
			err.Error(),
		)

		return
	}

	state := dedicatedServers

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
