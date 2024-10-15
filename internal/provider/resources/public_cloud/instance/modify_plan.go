package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	instanceValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/validator"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

// ModifyPlan calls validators that require access to the handler.
// This needs to be done here as client.Client isn't properly initialized when
// the schema is called.
func (i *instanceResource) ModifyPlan(
	ctx context.Context,
	request resource.ModifyPlanRequest,
	response *resource.ModifyPlanResponse,
) {
	planInstance := model.Instance{}
	request.Plan.Get(ctx, &planInstance)

	planImage := model.Image{}
	planInstance.Image.As(ctx, &planImage, basetypes.ObjectAsOptions{})

	stateInstance := model.Instance{}
	request.State.Get(ctx, &stateInstance)

	stateImage := model.Image{}
	stateInstance.Image.As(ctx, &stateImage, basetypes.ObjectAsOptions{})

	// Before deletion, determine if the instance is allowed to be deleted
	if request.Plan.Raw.IsNull() {
		i.validateInstance(stateInstance, response, ctx)
		if response.Diagnostics.HasError() {
			return
		}
	}

	i.validateRegion(planInstance.Region, response, ctx)
	if response.Diagnostics.HasError() {
		return
	}

	i.validateInstanceType(
		planInstance.Type,
		stateInstance.Type,
		stateInstance.Id,
		planInstance.Region,
		response,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}
}

func (i *instanceResource) validateRegion(
	plannedValue types.String,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: plannedValue}
	regionResponse := validator.StringResponse{}

	regionValidator := instanceValidator.NewRegionValidator(
		i.client.PublicCloudService.DoesRegionExist,
	)
	regionValidator.ValidateString(ctx, request, &regionResponse)
	if regionResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(regionResponse.Diagnostics.Errors()...)
	}
}

func (i *instanceResource) validateInstanceType(
	instanceType types.String,
	currentInstanceType types.String,
	instanceId types.String,
	region types.String,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: instanceType}
	instanceResponse := validator.StringResponse{}

	instanceTypeValidator := instanceValidator.NewInstanceTypeValidator(
		i.client.PublicCloudService.IsInstanceTypeAvailableForRegion,
		i.client.PublicCloudService.CanInstanceTypeBeUsedWithInstance,
		instanceId,
		region,
		currentInstanceType,
	)

	instanceTypeValidator.ValidateString(ctx, request, &instanceResponse)
	if instanceResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(instanceResponse.Diagnostics.Errors()...)
	}
}

// Checks if instance can be deleted.
func (i *instanceResource) validateInstance(
	instance model.Instance,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	instanceObject, diags := basetypes.NewObjectValueFrom(
		ctx,
		model.Instance{}.AttributeTypes(),
		instance,
	)
	if diags.HasError() {
		response.Diagnostics.Append(diags.Errors()...)
		return
	}

	instanceRequest := validator.ObjectRequest{ConfigValue: instanceObject}
	instanceResponse := validator.ObjectResponse{}
	validateInstanceTermination := instanceValidator.ValidateInstanceTermination(
		i.client.PublicCloudService.CanInstanceBeTerminated,
	)
	validateInstanceTermination.ValidateObject(
		ctx,
		instanceRequest,
		&instanceResponse,
	)

	if instanceResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(instanceResponse.Diagnostics.Errors()...)
	}
}
