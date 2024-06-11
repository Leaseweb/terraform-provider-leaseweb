package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-leaseweb/internal/public_cloud/opts"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
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

	request := i.client.SdkClient.PublicCloudAPI.
		LaunchInstance(i.client.AuthContext()).
		LaunchInstanceOpts(*instanceOpts.NewLaunchInstanceOpts())
	instance, _, err := i.client.SdkClient.PublicCloudAPI.LaunchInstanceExecute(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Public Cloud Instance",
			"Could not create Public Cloud Instance, unexpected error: "+err.Error(),
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
