package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	instanceValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/validator"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared_validators"
)

type immutableString struct {
	stateValue   types.String
	plannedValue types.String
}

type immutableStrings []immutableString

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

	i.validateRegion(planInstance.Region, response, ctx)
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

	immutableStrings := immutableStrings{
		immutableString{
			stateValue:   stateInstance.Region,
			plannedValue: planInstance.Region,
		},
		immutableString{
			stateValue:   stateImage.Id,
			plannedValue: planImage.Id,
		},
		immutableString{
			stateValue:   stateInstance.MarketAppId,
			plannedValue: planInstance.MarketAppId,
		},
		immutableString{
			stateValue:   stateInstance.RootDiskStorageType,
			plannedValue: planInstance.RootDiskStorageType,
		},
		// TODO Enable SSH key support
		/**
		  immutableString{
		  	stateValue:   stateInstance.SshKey,
		  	plannedValue: planInstance.SshKey,
		  },
		*/
	}

	i.validateImmutableString(stateInstance.Id, immutableStrings, response, ctx)
}

func (i *instanceResource) validateImmutableString(
	stateIdValue types.String,
	immutableStrings immutableStrings,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	for _, immutableString := range immutableStrings {
		request := validator.StringRequest{ConfigValue: immutableString.plannedValue}
		validatorResponse := validator.StringResponse{}

		immutableStringValidator := shared_validators.NewImmutableStringValidator(
			stateIdValue,
			immutableString.stateValue,
		)
		immutableStringValidator.ValidateString(ctx, request, &validatorResponse)
		if validatorResponse.Diagnostics.HasError() {
			response.Diagnostics.Append(validatorResponse.Diagnostics.Errors()...)
			return
		}
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
