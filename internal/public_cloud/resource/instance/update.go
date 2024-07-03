package instance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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

	var autoScalingGroupDetails *publicCloud.AutoScalingGroupDetails
	var autoScalingGroupDetailsResponse *http.Response
	var err error

	// Get autoScalingGroup details for each Instance as the instanceDetails
	// endpoint is missing loadBalancer data.
	autoScalingGroup, _ := instance.GetAutoScalingGroupOk()
	if autoScalingGroup != nil {
		autoScalingGroupDetails, autoScalingGroupDetailsResponse, err = i.client.PublicCloudClient.PublicCloudAPI.GetAutoScalingGroup(
			i.client.AuthContext(ctx),
			autoScalingGroup.GetId(),
		).Execute()
		if err != nil {
			utils.HandleError(
				ctx,
				autoScalingGroupDetailsResponse,
				&resp.Diagnostics,
				fmt.Sprintf(
					"Unable to Read Leaseweb Public Cloud Instance %v",
					instance.GetId(),
				),
				err.Error(),
			)
			return
		}
	}

	plan.Populate(instance, autoScalingGroupDetails, ctx)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
