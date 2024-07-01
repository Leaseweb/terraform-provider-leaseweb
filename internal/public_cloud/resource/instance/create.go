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

	launchInstanceRequest := i.client.PublicCloudClient.PublicCloudAPI.
		LaunchInstance(i.client.AuthContext(ctx)).
		LaunchInstanceOpts(*launchInstanceOpts)

	launchedInstance, launchInstanceSdkResponse, launchInstanceSdkError := i.client.PublicCloudClient.PublicCloudAPI.
		LaunchInstanceExecute(launchInstanceRequest)

	if launchInstanceSdkError != nil {
		utils.HandleError(
			ctx,
			launchInstanceSdkResponse,
			&resp.Diagnostics,
			"Error creating Public Cloud Instance",
			launchInstanceSdkError.Error(),
		)
		return
	}

	instanceRequest, instanceSdkResponse, instanceSdkError := i.client.PublicCloudClient.PublicCloudAPI.
		GetInstance(i.client.AuthContext(ctx), launchedInstance.GetId()).Execute()

	if instanceSdkError != nil {
		utils.HandleError(
			ctx,
			instanceSdkResponse,
			&resp.Diagnostics,
			"Error creating Public Cloud Instance",
			instanceSdkError.Error(),
		)
		return
	}

	diags = plan.Populate(instanceRequest, ctx)
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
