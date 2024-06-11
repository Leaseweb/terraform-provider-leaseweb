package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/public_cloud/opts"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
)

func (i *instanceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan model.Instance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceOpts := opts.NewInstanceOpts(plan, ctx)
	updateInstanceOpts := instanceOpts.NewUpdateInstanceOpts()

	request := i.client.SdkClient.PublicCloudAPI.UpdateInstance(
		i.client.AuthContext(),
		plan.Id.ValueString(),
	).UpdateInstanceOpts(*updateInstanceOpts)
	instance, _, err := i.client.SdkClient.PublicCloudAPI.UpdateInstanceExecute(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Public Cloud Instance",
			"Could not update Public Cloud Instance, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Populate(instance, ctx)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
