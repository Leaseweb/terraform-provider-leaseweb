package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

var (
	_ resource.Resource                = &instanceResource{}
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
	_ resource.ResourceWithModifyPlan  = &instanceResource{}
)

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

type instanceResource struct {
	client client.Client
}

func (i *instanceResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	coreClient, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	i.client = coreClient
}

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

	tflog.Info(ctx, "Creating public cloud instance")
	instance, err := i.client.PublicCloudFacade.CreateInstance(plan, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Instance", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Error creating public cloud instance",
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (i *instanceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state model.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Deleting public cloud instance %q",
		state.Id.ValueString(),
	))
	err := i.client.PublicCloudFacade.DeleteInstance(state.Id.ValueString(), ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating Public Cloud Instance",
			fmt.Sprintf(
				"Could not terminate Public Cloud Instance, unexpected error: %q",
				err.Error(),
			),
		)

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Error deleting public cloud instance %q",
				state.Id.ValueString(),
			),
			err.Error(),
		)

		return
	}
}

func (i *instanceResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {

	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		req,
		resp,
	)
}

func (i *instanceResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_instance"
}

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

	planRegion := model.Region{}
	planInstance.Region.As(ctx, &planRegion, basetypes.ObjectAsOptions{})

	stateInstance := model.Instance{}
	request.State.Get(ctx, &stateInstance)

	stateInstanceType := model.InstanceType{}
	stateInstance.Type.As(ctx, &stateInstanceType, basetypes.ObjectAsOptions{})

	stateImage := model.Image{}
	stateInstance.Image.As(ctx, &stateImage, basetypes.ObjectAsOptions{})

	// Before deletion, determine if the instance is allowed to be deleted
	if request.Plan.Raw.IsNull() {
		i.validateInstance(stateInstance, response, ctx)
		if response.Diagnostics.HasError() {
			return
		}
	}

	i.validateRegion(planRegion.Name, response, ctx)
	if response.Diagnostics.HasError() {
		return
	}

	i.validateInstanceType(
		planInstanceType.Name,
		stateInstanceType.Name,
		stateInstance.Id,
		planRegion.Name,
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
		i.client.PublicCloudFacade.DoesRegionExist,
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
		i.client.PublicCloudFacade.IsInstanceTypeAvailableForRegion,
		i.client.PublicCloudFacade.CanInstanceTypeBeUsedWithInstance,
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
		i.client.PublicCloudFacade.CanInstanceBeTerminated,
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

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state model.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Read public cloud instance %q",
		state.Id.ValueString(),
	))
	instance, err := i.client.PublicCloudFacade.GetInstance(
		state.Id.ValueString(),
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Instance", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf("Unable to read instance %q", state.Id.ValueString()),
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

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

	tflog.Info(ctx, fmt.Sprintf(
		"Updating public cloud instance %q",
		plan.Id.ValueString(),
	))
	updatedInstance, err := i.client.PublicCloudFacade.UpdateInstance(plan, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error updating instance", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Unable to update public cloud instance %q",
				plan.Id.ValueString(),
			),
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, updatedInstance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
