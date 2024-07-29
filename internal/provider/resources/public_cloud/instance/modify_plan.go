package instance

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/modify_plan"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

// ModifyPlan checks that the planned instanceType & region are valid.
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

	// Only validate region if it changes
	if planInstance.Region.ValueString() != "" {
		err := i.validateRegion(ctx, response, planInstance.Region.ValueString())
		if err != nil {
			log.Fatal(err)
		}
	}

	err := i.validateInstanceType(
		typeValidator,
		planInstance.Region,
		stateInstance.Id,
		planInstance.Type,
		response,
		ctx,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *instanceResource) validateInstanceType(
	typeValidator modify_plan.TypeValidator,
	planRegion types.String,
	stateId types.String,
	planInstanceType types.String,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) error {
	if typeValidator.IsBeingCreated() {
		return i.validateInstanceTypeForCreate(
			ctx,
			planRegion.ValueString(),
			planInstanceType.ValueString(),
			response,
		)
	}

	return i.validateInstanceTypeForUpdate(
		ctx,
		typeValidator,
		stateId.ValueString(),
		planInstanceType.ValueString(),
		response,
	)
}

func (i *instanceResource) validateInstanceTypeForCreate(
	ctx context.Context,
	region string,
	instanceType string,
	response *resource.ModifyPlanResponse,
) error {
	isAvailable, availableInstanceTypes, err := i.client.PublicCloudHandler.IsInstanceTypeAvailableForRegion(
		instanceType,
		region,
		ctx,
	)

	if err != nil {
		return fmt.Errorf("validateInstanceTypeForCreate: %w", err)
	}

	if isAvailable {
		return nil
	}

	response.Diagnostics.AddAttributeError(
		path.Root("type"),
		"Invalid Type",
		fmt.Sprintf(
			"Attribute type value must be one of: %q, got: %q",
			availableInstanceTypes,
			instanceType,
		),
	)

	return nil
}

func (i *instanceResource) validateInstanceTypeForUpdate(
	ctx context.Context,
	typeValidator modify_plan.TypeValidator,
	id string,
	instanceType string,
	response *resource.ModifyPlanResponse,
) error {
	if !typeValidator.HasTypeChanged() {
		return nil
	}

	allowedInstanceTypes, err := i.client.PublicCloudHandler.GetAvailableInstanceTypesForUpdate(
		id,
		ctx,
	)

	if err != nil {
		return fmt.Errorf("validateInstanceTypeForUpdate: %w", err)
	}

	if typeValidator.IsTypeValid(allowedInstanceTypes) {
		return nil
	}

	response.Diagnostics.AddAttributeError(
		path.Root("type"),
		"Invalid Type",
		fmt.Sprintf(
			"Attribute type value must be one of: %q, got: %q",
			allowedInstanceTypes,
			instanceType,
		),
	)

	return nil
}

func (i *instanceResource) validateRegion(
	ctx context.Context,
	response *resource.ModifyPlanResponse,
	region string,
) error {
	regionIsValid, validRegions, err := i.client.PublicCloudHandler.IsRegionValid(
		region,
		ctx,
	)

	if err != nil {
		return fmt.Errorf("validateRegion: %w", err)
	}

	if regionIsValid {
		return nil
	}

	response.Diagnostics.AddAttributeError(
		path.Root("region"),
		"Invalid Region",
		fmt.Sprintf(
			"Attribute region value must be one of: %q, got: %q",
			validRegions,
			region,
		),
	)

	return nil
}
