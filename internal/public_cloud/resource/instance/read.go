package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
)

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state model.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, _, err := i.client.SdkClient.PublicCloudAPI.GetInstance(
		i.client.AuthContext(),
		state.Id.ValueString(),
	).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Public Cloud Instance",
			"Could not read Public Cloud Instance ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Populate(instance, ctx)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
