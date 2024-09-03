package dedicated_server

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

func (d *dedicatedServerResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state model.DedicatedServer
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Read dedicated server %q",
		state.Id.ValueString(),
	))
	dedicatedServer, err := d.client.DedicatedServerFacade.GetDedicatedServer(
		ctx,
		state.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading dedicated server", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf("Unable to read dedicated server %q", state.Id.ValueString()),
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, dedicatedServer)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
