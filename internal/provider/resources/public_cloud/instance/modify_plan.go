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

	planInstanceType := model.InstanceType{}
	planInstance.Type.As(ctx, &planInstanceType, basetypes.ObjectAsOptions{})

	planImage := model.Image{}
	planInstance.Image.As(ctx, &planImage, basetypes.ObjectAsOptions{})

	stateInstance := model.Instance{}
	request.State.Get(ctx, &stateInstance)

	stateImage := model.Image{}
	stateInstance.Image.As(ctx, &stateImage, basetypes.ObjectAsOptions{})

	i.validateRegion(stateInstance.Region, planInstance.Region, response, ctx)
	if response.Diagnostics.HasError() {
		return
	}

	i.validateInstanceType(
		planInstanceType.Name,
		stateInstance.Id,
		planInstance.Region,
		response,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	i.validateImageId(stateImage.Id, planImage.Id, response, ctx)
}

func (i *instanceResource) validateImageId(
	stateValue types.String,
	plannedValue types.String,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: plannedValue}
	imageIdResponse := validator.StringResponse{}

	nonUpdatableStringValidator := instanceValidator.NewNonUpdatableStringValidator(stateValue)
	nonUpdatableStringValidator.ValidateString(ctx, request, &imageIdResponse)
	if imageIdResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(imageIdResponse.Diagnostics.Errors()...)
		return
	}
}

func (i *instanceResource) validateRegion(
	stateValue types.String,
	plannedValue types.String,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: plannedValue}
	regionResponse := validator.StringResponse{}

	nonUpdatableStringValidator := instanceValidator.NewNonUpdatableStringValidator(stateValue)
	nonUpdatableStringValidator.ValidateString(ctx, request, &regionResponse)
	if regionResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(regionResponse.Diagnostics.Errors()...)
		return
	}

	regionValidator := instanceValidator.NewRegionValidator(
		i.client.PublicCloudFacade.DoesRegionExist,
	)
	regionValidator.ValidateString(ctx, request, &regionResponse)
	if regionResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(regionResponse.Diagnostics.Errors()...)
	}
}

func (i *instanceResource) validateInstanceType(
	instanceType types.String,
	instanceId types.String,
	region types.String,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: instanceType}
	instanceResponse := validator.StringResponse{}

	instanceTypeValidator := instanceValidator.NewInstanceTypeValidator(
		i.client.PublicCloudFacade.IsInstanceTypeAvailableForRegion,
		i.client.PublicCloudFacade.CanInstanceTypeBeUsedWithInstance,
		instanceId,
		region,
	)

	instanceTypeValidator.ValidateString(ctx, request, &instanceResponse)
	if instanceResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(instanceResponse.Diagnostics.Errors()...)
	}
}
