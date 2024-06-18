package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

	typeValidator := modify_plan.NewTypeValidator(
		stateInstance.Id,
		stateInstance.Type,
		planInstance.Type,
	)

	instanceTypes := modify_plan.NewInstanceTypes(*i.client, ctx)

	hasTypeChanged := typeValidator.HashTypeChanged()

	if !hasTypeChanged {
		return
	}

	allowedInstanceTypes, sdkResponse, err := instanceTypes.
		GetAllowedInstanceTypes(stateInstance.Id.ValueString())

	if err != nil {
		utils.HandleError(
			ctx,
			sdkResponse,
			&response.Diagnostics,
			fmt.Sprintf(
				"Error getting updateInstanceType list for %q",
				stateInstance.Id.ValueString(),
			),
			err.Error(),
		)
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
