package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/dataadapters/shared"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
	customValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/validator"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service"
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
	var plan resourceModel.Instance

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Launch public cloud instance on API")

	opts, err := plan.GetLaunchInstanceOpts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating launch instance opts",
			err.Error(),
		)

		return
	}

	sdkInstance, repositoryErr := i.client.PublicCloudRepository.LaunchInstance(
		*opts,
		ctx,
	)
	if repositoryErr != nil {
		resp.Diagnostics.AddError(
			"Error creating Instance",
			repositoryErr.Error(),
		)

		logging.ServiceError(
			ctx,
			repositoryErr.ErrorResponse,
			&resp.Diagnostics,
			"Error launching public cloud instance",
			repositoryErr.Error(),
		)

		return
	}

	instance, resourceErr := resourceModel.NewFromInstance(*sdkInstance, ctx)
	if resourceErr != nil {
		resp.Diagnostics.AddError(
			"Error creating public cloud instance resource",
			resourceErr.Error(),
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
	var state resourceModel.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Terminate public cloud instance %q",
		state.Id.ValueString(),
	))
	err := i.client.PublicCloudRepository.DeleteInstance(state.Id.ValueString(), ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating Public Cloud Instance",
			fmt.Sprintf(
				"Could not terminate Public Cloud Instance, unexpected error: %q",
				err.Error(),
			),
		)

		logging.ServiceError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Error terminating public cloud instance %q",
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

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state resourceModel.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Read public cloud instance %q",
		state.Id.ValueString(),
	))
	sdkInstance, err := i.client.PublicCloudRepository.GetInstance(
		state.Id.ValueString(),
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Instance", err.Error())

		logging.ServiceError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf("Unable to read Instance %q", state.Id.ValueString()),
			err.Error(),
		)

		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Create public cloud instance resource for %q",
		state.Id.ValueString(),
	))
	instance, resourceErr := resourceModel.NewFromInstanceDetails(
		*sdkInstance,
		ctx,
	)
	if resourceErr != nil {
		resp.Diagnostics.AddError(
			"Error creating public cloud instance resource",
			resourceErr.Error(),
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
	var plan resourceModel.Instance

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Update public cloud instance %q",
		plan.Id.ValueString(),
	))
	opts, err := plan.GetUpdateInstanceOpts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating UpdateInstanceOpts",
			err.Error(),
		)
		return
	}

	sdkInstance, repositoryErr := i.client.PublicCloudRepository.UpdateInstance(
		plan.Id.ValueString(),
		*opts,
		ctx,
	)
	if repositoryErr != nil {
		resp.Diagnostics.AddError(
			"Error updating instance",
			repositoryErr.Error(),
		)

		logging.ServiceError(
			ctx,
			repositoryErr.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Unable to update public cloud instance %q",
				plan.Id.ValueString(),
			),
			repositoryErr.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, sdkInstance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (i *instanceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := service.NewIntMarkdownList(
		append(
			[]publicCloud.BillingFrequency{0},
			publicCloud.AllowedBillingFrequencyEnumValues...,
		),
	)
	contractTerms := service.NewIntMarkdownList(publicCloud.AllowedContractTermEnumValues)
	warningError := "**WARNING!** Changing this value once running will cause this instance to be destroyed and a new one to be created."

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The instance unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Our current regions can be found in the [developer documentation](https://developer.leaseweb.com/api-docs/publiccloud_v1.html#tag/Instances/operation/launchInstance)" + warningError,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reference": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The identifying name set to the instance",
			},
			"image": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Can be either an Operating System or a UUID in case of a Custom Image ID." + warningError,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The instance's current state",
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.Expressions{path.MatchRoot("region")}...,
					),
				},
			},
			"root_disk_size": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances. The maximum size is 1000 GB",
				Validators: []validator.Int64{
					int64validator.Between(5, 1000),
				},
			},
			"root_disk_storage_type": schema.StringAttribute{
				Required:    true,
				Description: "The root disk's storage type. Can be *LOCAL* or *CENTRAL*. " + warningError,
				Validators: []validator.String{
					stringvalidator.OneOf(shared.AdaptStringTypeArrayToStringArray(publicCloud.AllowedStorageTypeEnumValues)...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{Computed: true},
					},
				},
			},
			"contract": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"billing_frequency": schema.Int64Attribute{
						Required:    true,
						Description: "The billing frequency (in months). Valid options are " + billingFrequencies.Markdown(),
						Validators: []validator.Int64{
							int64validator.OneOf(billingFrequencies.ToInt64()...),
						},
					},
					"term": schema.Int64Attribute{
						Required:    true,
						Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + contractTerms.Markdown(),
						Validators: []validator.Int64{
							int64validator.OneOf(contractTerms.ToInt64()...),
						},
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
						Validators: []validator.String{
							stringvalidator.OneOf(shared.AdaptStringTypeArrayToStringArray(publicCloud.AllowedContractTypeEnumValues)...),
						},
					},
					"ends_at": schema.StringAttribute{Computed: true},
					"state": schema.StringAttribute{
						Computed: true,
					},
				},
				Validators: []validator.Object{customValidator.ContractTermIsValid()},
			},
			"market_app_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Market App ID that must be installed into the instance." + warningError,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

// ModifyPlan calls validators that require access to the handler.
// This needs to be done here as client.Client isn't properly initialized when
// the schema is called.
func (i *instanceResource) ModifyPlan(
	ctx context.Context,
	request resource.ModifyPlanRequest,
	response *resource.ModifyPlanResponse,
) {
	planInstance := resourceModel.Instance{}
	request.Plan.Get(ctx, &planInstance)

	planImage := resourceModel.Image{}
	planInstance.Image.As(ctx, &planImage, basetypes.ObjectAsOptions{})

	stateInstance := resourceModel.Instance{}
	request.State.Get(ctx, &stateInstance)

	stateImage := resourceModel.Image{}
	stateInstance.Image.As(ctx, &stateImage, basetypes.ObjectAsOptions{})

	// Before deletion, determine if the instance is allowed to be deleted
	if request.Plan.Raw.IsNull() {
		i.validateInstance(stateInstance, response, ctx)
		if response.Diagnostics.HasError() {
			return
		}
	}

	regions, err := i.client.PublicCloudRepository.GetRegions(ctx)
	if err != nil {
		response.Diagnostics.AddError("Cannot get regions", err.Error())
		return
	}

	// The Region has
	//to be validated first or getAvailableInstanceTypes will throw an error on creation,
	//as the region could be invalid.
	i.validateRegion(planInstance.Region, response, regions, ctx)
	if response.Diagnostics.HasError() {
		return
	}

	availableInstanceTypes := i.getAvailableInstanceTypes(
		response,
		stateInstance.Id,
		planInstance.Region.ValueString(),
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	i.validateInstanceType(
		planInstance.Type,
		stateInstance.Type,
		response,
		availableInstanceTypes,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}
}

// When creating a new Instance,
// any instanceType available in the region is good.
// On update, the criteria is more limited.
func (i *instanceResource) getAvailableInstanceTypes(
	response *resource.ModifyPlanResponse,
	id basetypes.StringValue,
	region string,
	ctx context.Context,
) []string {
	// Instance is being created.
	if id.IsNull() {
		availableInstanceTypes, err := i.client.PublicCloudRepository.GetInstanceTypesForRegion(region, ctx)
		if err != nil {
			response.Diagnostics.AddError("Cannot get available instanceTypes for region", err.Error())
			return nil
		}

		return availableInstanceTypes
	}

	availableInstanceTypes, err := i.client.PublicCloudRepository.GetAvailableInstanceTypesForUpdate(id.ValueString(), ctx)
	if err != nil {
		response.Diagnostics.AddError("Cannot get available instanceTypes for update", err.Error())
		return nil
	}

	return availableInstanceTypes
}

func (i *instanceResource) validateRegion(
	plannedValue types.String,
	response *resource.ModifyPlanResponse,
	regions []string,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: plannedValue}
	regionResponse := validator.StringResponse{}

	regionValidator := customValidator.NewRegionValidator(regions)
	regionValidator.ValidateString(ctx, request, &regionResponse)
	if regionResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(regionResponse.Diagnostics.Errors()...)
	}
}

func (i *instanceResource) validateInstanceType(
	instanceType types.String,
	currentInstanceType types.String,
	response *resource.ModifyPlanResponse,
	availableInstanceTypes []string,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: instanceType}
	instanceResponse := validator.StringResponse{}

	instanceTypeValidator := customValidator.NewInstanceTypeValidator(
		currentInstanceType,
		availableInstanceTypes,
	)

	instanceTypeValidator.ValidateString(ctx, request, &instanceResponse)
	if instanceResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(instanceResponse.Diagnostics.Errors()...)
	}
}

// Checks if instance can be deleted.
func (i *instanceResource) validateInstance(
	instance resourceModel.Instance,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	instanceObject, diags := basetypes.NewObjectValueFrom(
		ctx,
		resourceModel.Instance{}.AttributeTypes(),
		instance,
	)
	if diags.HasError() {
		response.Diagnostics.Append(diags.Errors()...)
		return
	}

	instanceRequest := validator.ObjectRequest{ConfigValue: instanceObject}
	instanceResponse := validator.ObjectResponse{}
	validateInstanceTermination := customValidator.ValidateInstanceTermination(
		instance.CanBeTerminated,
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
