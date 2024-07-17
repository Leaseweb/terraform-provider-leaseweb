package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
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

	err := i.client.PublicCloudHandler.DeleteInstance(state.Id.ValueString(), ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating Public Cloud Instance",
			fmt.Sprintf(
				"Could not terminate Public CLoud Instance, unexpected error: %q",
				err.Error(),
			),
		)
		return
	}
}
