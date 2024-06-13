package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/public_cloud/opts"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"terraform-provider-leaseweb/internal/utils"
)

func (i *instanceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan model.Instance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceOpts := opts.NewInstanceOpts(plan, ctx)
	launchInstanceOpts, optsErr := instanceOpts.NewLaunchInstanceOpts()

	if optsErr != nil {
		utils.HandleError(
			ctx,
			nil,
			&resp.Diagnostics,
			"Error creating Public Cloud Instance",
			optsErr.Error(),
		)
		return
	}

	request := i.client.SdkClient.PublicCloudAPI.
		LaunchInstance(i.client.AuthContext(ctx)).
		LaunchInstanceOpts(*launchInstanceOpts)

	instance, sdkResponse, sdkErr := i.client.SdkClient.PublicCloudAPI.
		LaunchInstanceExecute(request)

	if sdkErr != nil {
		utils.HandleError(
			ctx,
			sdkResponse,
			&resp.Diagnostics,
			"Error creating Public Cloud Instance",
			sdkErr.Error(),
		)
		return
	}

	diags = plan.Populate(instance, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
