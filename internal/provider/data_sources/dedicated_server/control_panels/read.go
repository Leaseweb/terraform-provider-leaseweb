package control_panels

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
)

func (c *controlPanelDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	tflog.Info(ctx, "Read dedicated_server_control_panels")
	controlPanels, err := c.client.DedicatedServerFacade.GetAllControlPanels(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read dedicated_server_control_panels", err.Error())
		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Unable to read dedicated_server_control_panels",
			err.Error(),
		)

		return
	}

	state := controlPanels

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
