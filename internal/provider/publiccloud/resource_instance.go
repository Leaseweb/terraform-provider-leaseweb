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
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/service/public_cloud"
	publicCloud "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/service/public_cloud"
	customValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/validator"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/logging"
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

	tflog.Info(ctx, "Creating public cloud instance")
	instance, err := i.client.PublicCloudService.LaunchInstance(plan, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Instance", err.Error())

		logging.ServiceError(
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
	var state resourceModel.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Deleting public cloud instance %q",
		state.Id.ValueString(),
	))
	err := i.client.PublicCloudService.DeleteInstance(state.Id.ValueString(), ctx)

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
	instance, err := i.client.PublicCloudService.GetInstance(
		state.Id.ValueString(),
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Instance", err.Error())

		logging.ServiceError(
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
	var plan resourceModel.Instance

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Updating public cloud instance %q",
		plan.Id.ValueString(),
	))
	updatedInstance, err := i.client.PublicCloudService.UpdateInstance(plan, ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error updating instance", err.Error())

		logging.ServiceError(
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

func (i *instanceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	publicCloudService := publicCloud.Service{}
	service := public_cloud.Service{}
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
			// TODO Enable SSH key support
			/**
			  "ssh_key": schema.StringAttribute{
			  	Optional:      true,
			  	Sensitive:     true,
			  	PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			  	Description:   "Public SSH key to be installed into the instance. Must be used only on Linux/FreeBSD instances",
			  	Validators: []validator.String{
			  		stringvalidator.RegexMatches(
			  			regexp.MustCompile(service.GetSshKeyRegularExpression()),
			  			"Invalid ssh key",
			  		),
			  	},
			  },
			*/
			"root_disk_size": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances. The maximum size is 1000 GB",
				Validators: []validator.Int64{
					int64validator.Between(
						service.GetMinimumRootDiskSize(),
						service.GetMaximumRootDiskSize(),
					),
				},
			},
			"root_disk_storage_type": schema.StringAttribute{
				Required:    true,
				Description: "The root disk's storage type. Can be *LOCAL* or *CENTRAL*. " + warningError,
				Validators: []validator.String{
					stringvalidator.OneOf(service.GetRootDiskStorageTypes()...),
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
						Description: "The billing frequency (in months). Valid options are " + publicCloudService.GetBillingFrequencies().Markdown(),
						Validators: []validator.Int64{
							int64validator.OneOf(publicCloudService.GetBillingFrequencies().ToInt64()...),
						},
					},
					"term": schema.Int64Attribute{
						Required:    true,
						Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + publicCloudService.GetContractTerms().Markdown(),
						Validators: []validator.Int64{
							int64validator.OneOf(publicCloudService.GetContractTerms().ToInt64()...),
						},
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
						Validators: []validator.String{
							stringvalidator.OneOf(publicCloudService.GetContractTypes()...),
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

	regionValidator := customValidator.NewRegionValidator(
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

	instanceTypeValidator := customValidator.NewInstanceTypeValidator(
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
