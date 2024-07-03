package instances

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/public_cloud/data_source/instances/model"
	"terraform-provider-leaseweb/internal/utils"
)

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	sdkInstancesList, sdkResponse, err := d.client.PublicCloudClient.PublicCloudAPI.
		GetInstanceList(d.client.AuthContext(ctx)).Execute()
	if err != nil {
		utils.HandleError(
			ctx,
			sdkResponse,
			&resp.Diagnostics,
			"Unable to Read Leaseweb Public Cloud Instances",
			err.Error(),
		)
		return
	}

	var sdkInstances []publicCloud.InstanceDetails
	var sdkAutoScalingGroups []publicCloud.AutoScalingGroupDetails
	var sdkLoadBalancers []publicCloud.LoadBalancerDetails

	for _, sdkInstance := range sdkInstancesList.Instances {
		sdkInstanceDetails, err := d.getInstanceDetails(&sdkInstance, ctx, resp)
		if err != nil {
			return
		}

		sdkInstances = append(sdkInstances, *sdkInstanceDetails)

		sdkAutoScalingGroupDetails, err := d.getAutoScalingGroupDetails(
			*sdkInstanceDetails,
			ctx,
			resp,
		)
		if err != nil {
			return
		}

		if sdkAutoScalingGroupDetails != nil {
			sdkAutoScalingGroups = append(
				sdkAutoScalingGroups,
				*sdkAutoScalingGroupDetails,
			)

			sdkLoadBalancerDetails, err := d.getLoadBalancerDetails(*sdkAutoScalingGroupDetails, ctx, resp)
			if err != nil {
				return
			}
			if sdkLoadBalancerDetails != nil {
				sdkLoadBalancers = append(sdkLoadBalancers, *sdkLoadBalancerDetails)
			}
		}

	}

	state := model.Instances{}
	err = state.Populate(sdkInstances, sdkAutoScalingGroups, sdkLoadBalancers)
	if err != nil {
		utils.HandleError(
			ctx,
			nil,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Unable to get populate details for Public Cloud Instance %v",
				sdkInstances,
			),
			err.Error(),
		)
		return
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Get instanceDetails for passed instance.
func (d *instancesDataSource) getInstanceDetails(
	instance *publicCloud.Instance,
	ctx context.Context,
	resp *datasource.ReadResponse,
) (
	sdkInstanceDetails *publicCloud.InstanceDetails,
	err error,
) {
	var sdkInstanceDetailsResponse *http.Response

	sdkInstanceDetails, sdkInstanceDetailsResponse, err = d.client.PublicCloudClient.PublicCloudAPI.GetInstance(
		d.client.AuthContext(ctx),
		instance.GetId(),
	).Execute()

	if err != nil {
		utils.HandleError(
			ctx,
			sdkInstanceDetailsResponse,
			&resp.Diagnostics,
			fmt.Sprintf("Unable to Read instance details for %v", instance.GetId()),
			err.Error(),
		)
		return nil, err
	}

	return sdkInstanceDetails, nil
}

// Get autoScalingGroupDetails for passed Instance as the instanceDetails
// endpoint is missing loadBalancer data.
func (d *instancesDataSource) getAutoScalingGroupDetails(
	sdkInstanceDetails publicCloud.InstanceDetails,
	ctx context.Context,
	resp *datasource.ReadResponse,
) (
	sdkAutoScalingGroupDetails *publicCloud.AutoScalingGroupDetails, err error,
) {
	var response *http.Response

	sdkAutoScalingGroup, _ := sdkInstanceDetails.GetAutoScalingGroupOk()
	if sdkAutoScalingGroup == nil {
		return nil, nil
	}

	sdkAutoScalingGroupDetails, response, err = d.client.PublicCloudClient.PublicCloudAPI.GetAutoScalingGroup(
		d.client.AuthContext(ctx),
		sdkAutoScalingGroup.GetId(),
	).Execute()
	if err != nil {
		utils.HandleError(
			ctx,
			response,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Unable to Read Leaseweb Public Cloud Instance %v",
				sdkInstanceDetails.GetId(),
			),
			err.Error(),
		)
		return nil, err
	}

	return sdkAutoScalingGroupDetails, nil
}

// Get loadBalancerDetails for passed autoScalingGroup as the autoScalingGroupDetails
// endpoint is missing loadBalancer data.
func (d *instancesDataSource) getLoadBalancerDetails(
	sdkAutoScalingGroupDetails publicCloud.AutoScalingGroupDetails,
	ctx context.Context,
	resp *datasource.ReadResponse,
) (
	sdkLoadBalancerDetails *publicCloud.LoadBalancerDetails, err error,
) {
	var response *http.Response

	sdkLoadBalancer, _ := sdkAutoScalingGroupDetails.GetLoadBalancerOk()
	if sdkLoadBalancer == nil {
		return nil, nil
	}

	sdkLoadBalancerDetails, response, err = d.client.PublicCloudClient.PublicCloudAPI.GetLoadBalancer(
		d.client.AuthContext(ctx),
		sdkLoadBalancer.GetId(),
	).Execute()
	if err != nil {
		utils.HandleError(
			ctx,
			response,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Unable to Read Leaseweb Public Cloud Load Balancer %v",
				sdkLoadBalancer.GetId(),
			),
			err.Error(),
		)
		return nil, err
	}

	return sdkLoadBalancerDetails, nil
}
