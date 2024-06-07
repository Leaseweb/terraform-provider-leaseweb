package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/resources/instance/model"
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

	request := i.client.SdkClient.PublicCloudAPI.TerminateInstance(
		i.client.AuthContext(),
		state.Id.ValueString(),
	)
	_, err := i.client.SdkClient.PublicCloudAPI.TerminateInstanceExecute(request)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating instance",
			"Could not terminate instance, unexpected error: "+err.Error(),
		)
		return
	}
}
