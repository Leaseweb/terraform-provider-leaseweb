package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/public_cloud/opts"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"terraform-provider-leaseweb/internal/utils"
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
	updateInstanceOpts, optsErr := instanceOpts.NewUpdateInstanceOpts()
	if optsErr != nil {
		utils.HandleError(
			ctx,
			nil,
			&resp.Diagnostics,
			"Error updating Public Cloud Instance",
			optsErr.Error(),
		)
		return
	}

	request := i.client.PublicCloudClient.PublicCloudAPI.UpdateInstance(
		i.client.AuthContext(ctx),
		plan.Id.ValueString(),
	).UpdateInstanceOpts(*updateInstanceOpts)
	instance, sdkResponse, sdkErr := i.client.PublicCloudClient.PublicCloudAPI.UpdateInstanceExecute(request)
	if sdkErr != nil {
		utils.HandleError(
			ctx,
			sdkResponse,
			&resp.Diagnostics,
			"Error updating Public Cloud Instance",
			sdkErr.Error(),
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
