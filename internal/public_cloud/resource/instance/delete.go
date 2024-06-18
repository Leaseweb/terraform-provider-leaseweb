package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
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

	request := i.client.PublicCloudClient.PublicCloudAPI.TerminateInstance(
		i.client.AuthContext(ctx),
		state.Id.ValueString(),
	)
	_, err := i.client.PublicCloudClient.PublicCloudAPI.TerminateInstanceExecute(request)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating Public Cloud Instance",
			"Could not terminate Public CLoud Instance, unexpected error: "+err.Error(),
		)
		return
	}
}
