package instance

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	validator2 "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/modify_plan"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/validator"
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

	i.validateRegion(planInstance.Region, response, ctx)
	if response.Diagnostics.HasError() {
		return
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

// Region validator has to be called here instead of in the schema as the client
// isn't initialized yet when the schema is generated.
func (i *instanceResource) validateRegion(
	region types.String,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	regionRequest := validator2.StringRequest{ConfigValue: region}
	regionResponse := validator2.StringResponse{}

	regionValidator := validator.NewRegionValidator(
		i.client.PublicCloudHandler.DoesRegionExist,
	)
	regionValidator.ValidateString(ctx, regionRequest, &regionResponse)

	if regionResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(regionResponse.Diagnostics.Errors()...)
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

// On creation, we need to check that the instance is available for the region.
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

// On update check that the passed instanceType can be used wit the instance.
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

	canInstanceTypeBeUsed, allowedInstanceTypes, err := i.client.PublicCloudHandler.CanInstanceTypeBeUsedWithInstance(
		id,
		instanceType,
		ctx,
	)

	if err != nil {
		return fmt.Errorf("validateInstanceTypeForUpdate: %w", err)
	}

	if canInstanceTypeBeUsed {
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
