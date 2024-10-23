package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &imageResource{}
	_ resource.ResourceWithImportState = &imageResource{}
	_ resource.ResourceWithModifyPlan  = &imageResource{}
	_ validator.String                 = instanceIdForCustomImageValidator{}
)

func getImage(
	ID string,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.ImageDetails, *utils.SdkError) {
	sdkImages, err := getAllImages(ctx, api)

	if err != nil {
		return nil, err
	}

	for _, sdkImage := range sdkImages {
		if sdkImage.GetId() == ID {
			return &sdkImage, nil
		}
	}

	return nil, nil
}

func updateImage(
	ID string,
	opts publicCloud.UpdateImageOpts,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.ImageDetails, *utils.SdkError) {
	image, response, err := api.UpdateImage(
		ctx,
		ID,
	).UpdateImageOpts(opts).Execute()
	if err != nil {
		return nil, utils.NewSdkError(
			fmt.Sprintf("updateImage %q", ID),
			err,
			response,
		)
	}

	return image, nil
}

func createImage(
	opts publicCloud.CreateImageOpts,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.ImageDetails, *utils.SdkError) {
	image, response, err := api.CreateImage(ctx).CreateImageOpts(opts).Execute()
	if err != nil {
		return nil, utils.NewSdkError("createImage", err, response)
	}

	return image, nil
}

type imageResource struct {
	client client.Client
}

func (i *imageResource) ModifyPlan(
	ctx context.Context,
	request resource.ModifyPlanRequest,
	response *resource.ModifyPlanResponse,
) {
	planImage := resourceModelImage{}
	request.Plan.Get(ctx, &planImage)

	instances, err := getAllInstances(ctx, i.client.PublicCloudAPI)
	if err != nil {
		response.Diagnostics.AddError("Cannot get instances", err.Error())
		return
	}

	idRequest := validator.StringRequest{ConfigValue: planImage.ID}
	idResponse := validator.StringResponse{}

	instanceIdValidator := newInstanceIdForCustomImageValidator(instances)
	instanceIdValidator.ValidateString(ctx, idRequest, &idResponse)

	if idResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(idResponse.Diagnostics.Errors()...)
	}
}

func (i *imageResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		request,
		response,
	)
}

func (i *imageResource) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = request.ProviderTypeName + "_public_cloud_image"
}

func (i *imageResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				Description: `
The id of the instance which the custom image is based on. The following rules apply:
  - instance exists for instanceId
  - instance has state *STOPPED*
  - instance has a maximum rootDiskSize of 100 GB
  - instance OS must not be *windows*`,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Custom image name",
			},
			"custom": schema.BoolAttribute{
				Computed:    true,
				Description: "Standard or Custom image",
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"market_apps": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"storage_types": schema.ListAttribute{
				Computed:    true,
				Description: "The supported storage types for the instance type",
				ElementType: types.StringType,
			},
			"flavour": schema.StringAttribute{
				Computed: true,
			},
			"region": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (i *imageResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan resourceModelImage

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create publiccloud image")

	opts := plan.GetCreateImageOpts()

	sdkImage, sdkErr := createImage(opts, ctx, i.client.PublicCloudAPI)
	if sdkErr != nil {
		response.Diagnostics.AddError(
			"Error creating publiccloud image",
			sdkErr.Error(),
		)

		utils.LogError(
			ctx,
			sdkErr.ErrorResponse,
			&response.Diagnostics,
			"Error creating publiccloud image",
			sdkErr.Error(),
		)

		return
	}

	instance, resourceErr := adaptSdkImageDetailsToResourceImage(ctx, *sdkImage)
	if resourceErr != nil {
		response.Diagnostics.AddError(
			"Error creating publiccloud image resource",
			resourceErr.Error(),
		)

		return
	}

	diags = response.State.Set(ctx, instance)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (i *imageResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state resourceModelImage

	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	sdkImage, err := getImage(state.ID.ValueString(), ctx, i.client.PublicCloudAPI)
	if err != nil {
		response.Diagnostics.AddError("Unable to read images", err.Error())
		utils.LogError(
			ctx,
			err.ErrorResponse,
			&response.Diagnostics,
			"Unable to read images",
			err.Error(),
		)

		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Create publiccloud image resource for %q",
		state.ID.ValueString(),
	))
	instance, resourceErr := adaptSdkImageDetailsToResourceImage(ctx, *sdkImage)
	if resourceErr != nil {
		response.Diagnostics.AddError(
			"Error creating publiccloud image resource",
			resourceErr.Error(),
		)

		return
	}

	diags = response.State.Set(ctx, instance)
	response.Diagnostics.Append(diags...)
}

func (i *imageResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan resourceModelImage

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Update publiccloud image %q",
		plan.ID.ValueString(),
	))
	opts := plan.GetUpdateImageOpts()

	sdkImageDetails, sdkErr := updateImage(
		plan.ID.ValueString(),
		opts,
		ctx,
		i.client.PublicCloudAPI,
	)
	if sdkErr != nil {
		response.Diagnostics.AddError(
			"Error updating publiccloud image",
			sdkErr.Error(),
		)

		utils.LogError(
			ctx,
			sdkErr.ErrorResponse,
			&response.Diagnostics,
			fmt.Sprintf(
				"Unable to update publiccloud image %q",
				plan.ID.ValueString(),
			),
			sdkErr.Error(),
		)

		return
	}

	diags = response.State.Set(ctx, sdkImageDetails)
	response.Diagnostics.Append(diags...)
}

// Delete does nothing as there is no endpoint to delete an Image.
func (i *imageResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}

func (i *imageResource) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
	if request.ProviderData == nil {
		return
	}

	coreClient, ok := request.ProviderData.(client.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				request.ProviderData,
			),
		)

		return
	}

	i.client = coreClient
}

func NewImageResource() resource.Resource {
	return &imageResource{}
}
