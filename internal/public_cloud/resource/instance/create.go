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

	sdkInstance, sdkInstanceResponse, sdkInstanceError := i.client.PublicCloudClient.PublicCloudAPI.
		GetInstance(i.client.AuthContext(ctx), launchedInstance.GetId()).Execute()

	if sdkInstanceError != nil {
		utils.HandleError(
			ctx,
			sdkInstanceResponse,
			&resp.Diagnostics,
			"Error creating Public Cloud Instance",
			sdkInstanceError.Error(),
		)
		return
	}

	var sdkAutoScalingGroupDetails *publicCloud.AutoScalingGroupDetails
	var sdkAutoScalingGroupDetailsResponse *http.Response
	var sdkLoadBalancerDetailsResponse *http.Response
	var sdkLoadBalancerDetails *publicCloud.LoadBalancerDetails
	var err error

	// Get autoScalingGroup details for the Instance as the instanceDetails
	// endpoint is missing data.
	sdkAutoScalingGroup, _ := sdkInstance.GetAutoScalingGroupOk()
	if sdkAutoScalingGroup != nil {
		sdkAutoScalingGroupDetails, sdkAutoScalingGroupDetailsResponse, err = i.client.PublicCloudClient.PublicCloudAPI.GetAutoScalingGroup(
			i.client.AuthContext(ctx),
			sdkAutoScalingGroup.GetId(),
		).Execute()
		if err != nil {
			utils.HandleError(
				ctx,
				sdkAutoScalingGroupDetailsResponse,
				&resp.Diagnostics,
				fmt.Sprintf(
					"Unable to Read Leaseweb AutoScalingGroup %v",
					sdkAutoScalingGroup.GetId(),
				),
				err.Error(),
			)
			return
		}

		// Get loadBalancer details for the AutoScalingGroup as the autoScalingGroupDetails
		// endpoint is missing data.
		if sdkAutoScalingGroupDetails != nil {
			sdkLoadBalancer, _ := sdkAutoScalingGroupDetails.GetLoadBalancerOk()
			if sdkLoadBalancer != nil {
				sdkLoadBalancerDetails, sdkLoadBalancerDetailsResponse, err = i.client.PublicCloudClient.PublicCloudAPI.GetLoadBalancer(
					i.client.AuthContext(ctx),
					sdkLoadBalancer.GetId(),
				).Execute()
				if err != nil {
					utils.HandleError(
						ctx,
						sdkLoadBalancerDetailsResponse,
						&resp.Diagnostics,
						fmt.Sprintf(
							"Unable to Read Leaseweb Configuration %v",
							sdkLoadBalancer.GetId(),
						),
						err.Error(),
					)
					return
				}

			}
		}
	}

	diags = plan.Populate(
		sdkInstance,
		sdkAutoScalingGroupDetails,
		sdkLoadBalancerDetails,
		ctx,
	)
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
