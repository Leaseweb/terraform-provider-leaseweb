package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/modify_plan"
	"terraform-provider-leaseweb/internal/utils"
)

func (i *instanceResource) ModifyPlan(
	ctx context.Context,
	request resource.ModifyPlanRequest,
	response *resource.ModifyPlanResponse,
) {
	planInstance := model.Instance{}
	request.Plan.Get(ctx, &planInstance)

	stateInstance := model.Instance{}
	request.State.Get(ctx, &stateInstance)

	i.validateInstanceType(
		ctx,
		stateInstance.Id,
		stateInstance.Type,
		planInstance.Type,
		response,
	)

	i.validateRegion(ctx, response, planInstance.Region.ValueString())
}

func (i *instanceResource) validateInstanceType(
	ctx context.Context,
	stateId types.String,
	stateType types.String,
	planType types.String,
	response *resource.ModifyPlanResponse,
) {
	typeValidator := modify_plan.NewTypeValidator(
		stateId,
		stateType,
		planType,
	)

	instanceTypes := modify_plan.NewInstanceTypes(*i.client, ctx)

	hasTypeChanged := typeValidator.HashTypeChanged()

	if !hasTypeChanged {
		return
	}

	allowedInstanceTypes, sdkResponse, err := instanceTypes.
		GetAllowedInstanceTypes(stateId.ValueString())

	if err != nil {
		utils.HandleError(
			ctx,
			sdkResponse,
			&response.Diagnostics,
			fmt.Sprintf(
				"Error getting updateInstanceType list for %q",
				stateId.ValueString(),
			),
			err.Error(),
		)
		return
	}

	if typeValidator.IsTypeValid(allowedInstanceTypes) {
		return
	}

	response.Diagnostics.AddAttributeError(
		path.Root("type"),
		"Invalid Instance Type",
		fmt.Sprintf(
			"Allowed types are %v",
			allowedInstanceTypes,
		),
	)
}

func (i *instanceResource) validateRegion(
	ctx context.Context,
	response *resource.ModifyPlanResponse,
	region string,
) {
	// Region has not changed here.
	if region == "" {
		return
	}

	request := i.client.PublicCloudClient.PublicCloudAPI.GetRegionList(i.client.AuthContext(ctx))
	sdkRegions, sdkResponse, err := i.client.PublicCloudClient.PublicCloudAPI.GetRegionListExecute(request)

	if err != nil {
		utils.HandleError(
			ctx,
			sdkResponse,
			&response.Diagnostics,
			"Error getting region list",
			err.Error(),
		)
		return
	}

	regions := modify_plan.NewRegions(sdkRegions.GetRegions())

	if regions.Contains(region) {
		return
	}

	response.Diagnostics.AddAttributeError(
		path.Root("region"),
		"Invalid Region",
		fmt.Sprintf(
			"Allowed regions are %v",
			regions,
		),
	)
}
