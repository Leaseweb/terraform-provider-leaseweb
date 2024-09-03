package dedicated_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
)

func (d *dedicatedServerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		req,
		resp,
	)

	dedicatedServer, err := d.client.DedicatedServerFacade.GetDedicatedServer(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing Dedicated server", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Error importing dedicated server",
			err.Error(),
		)

		return
	}

	diags := resp.State.Set(ctx, dedicatedServer)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
