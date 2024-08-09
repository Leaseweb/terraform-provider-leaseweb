package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func (i *instanceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state model.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Deleting public cloud instance %q",
		state.Id.ValueString(),
	))
	err := i.client.PublicCloudFacade.DeleteInstance(state.Id.ValueString(), ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating Public Cloud Instance",
			fmt.Sprintf(
				"Could not terminate Public Cloud Instance, unexpected error: %q",
				err.Error(),
			),
		)

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Error deleting public cloud instance %q",
				state.Id.ValueString(),
			),
			err.Error(),
		)

		return
	}
}
