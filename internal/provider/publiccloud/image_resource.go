package publiccloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &imageResource{}
	_ resource.ResourceWithImportState = &imageResource{}
)

type imageResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Custom       types.Bool   `tfsdk:"custom"`
	State        types.String `tfsdk:"state"`
	MarketApps   types.List   `tfsdk:"market_apps"`
	StorageTypes types.List   `tfsdk:"storage_types"`
	Flavour      types.String `tfsdk:"flavour"`
	Region       types.String `tfsdk:"region"`
}

func (i imageResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"name":          types.StringType,
		"custom":        types.BoolType,
		"state":         types.StringType,
		"market_apps":   types.ListType{ElemType: types.StringType},
		"storage_types": types.ListType{ElemType: types.StringType},
		"flavour":       types.StringType,
		"region":        types.StringType,
	}
}

func (i imageResourceModel) GetUpdateImageOpts() publicCloud.UpdateImageOpts {
	return publicCloud.UpdateImageOpts{
		Name: i.Name.ValueString(),
	}
}

func (i imageResourceModel) GetCreateImageOpts() publicCloud.CreateImageOpts {
	return publicCloud.CreateImageOpts{
		Name:       i.Name.ValueString(),
		InstanceId: i.ID.ValueString(),
	}
}

func adaptImageDetailsToImageResource(
	ctx context.Context,
	sdkImageDetails publicCloud.ImageDetails,
) (*imageResourceModel, error) {
	marketApps, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		sdkImageDetails.MarketApps,
	)
	if diags.HasError() {
		return nil, fmt.Errorf(diags.Errors()[0].Summary(), diags.Errors()[0].Detail())
	}

	storageTypes, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		sdkImageDetails.StorageTypes,
	)
	if diags.HasError() {
		return nil, fmt.Errorf(diags.Errors()[0].Summary(), diags.Errors()[0].Detail())
	}

	image := imageResourceModel{
		ID:           basetypes.NewStringValue(sdkImageDetails.GetId()),
		Name:         basetypes.NewStringValue(sdkImageDetails.GetName()),
		Custom:       basetypes.NewBoolValue(sdkImageDetails.GetCustom()),
		State:        basetypes.NewStringValue(string(sdkImageDetails.GetState())),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(sdkImageDetails.GetFlavour())),
		Region:       basetypes.NewStringValue(string(sdkImageDetails.GetRegion())),
	}

	return &image, nil
}

func adaptImageToImageResource(sdkImage publicCloud.Image) imageResourceModel {
	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})

	return imageResourceModel{
		ID:           basetypes.NewStringValue(sdkImage.GetId()),
		Name:         basetypes.NewStringValue(sdkImage.GetName()),
		Custom:       basetypes.NewBoolValue(sdkImage.GetCustom()),
		Flavour:      basetypes.NewStringValue(string(sdkImage.GetFlavour())),
		MarketApps:   emptyList,
		StorageTypes: emptyList,
	}
}

func getImage(
	ID string,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.ImageDetails, *http.Response, error) {
	sdkImages, httpResponse, err := getAllImages(ctx, api)

	if err != nil {
		return nil, httpResponse, err
	}

	for _, sdkImage := range sdkImages {
		if sdkImage.GetId() == ID {
			return &sdkImage, nil, nil
		}
	}

	return nil, nil, nil
}

type imageResource struct {
	client publicCloud.PublicCloudAPI
}

func (i *imageResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
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
		Description: "Once created, an image resource cannot be deleted via Terraform",
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
	var plan imageResourceModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Create publiccloud image")

	opts := plan.GetCreateImageOpts()

	sdkImage, httpResponse, err := i.client.CreateImage(ctx).
		CreateImageOpts(opts).
		Execute()
	if err != nil {
		utils.HandleSdkError(
			"Error creating Public Cloud image",
			httpResponse,
			err,
			&response.Diagnostics,
			ctx,
		)

		return
	}

	image, resourceErr := adaptImageDetailsToImageResource(ctx, *sdkImage)
	if resourceErr != nil {
		response.Diagnostics.AddError("Error creating publiccloud image resource", resourceErr.Error())

		return
	}

	diags = response.State.Set(ctx, image)
	response.Diagnostics.Append(diags...)
}

func (i *imageResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var state imageResourceModel

	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	sdkImage, httpResponse, err := getImage(state.ID.ValueString(), ctx, i.client)
	if err != nil {
		utils.HandleSdkError(
			"Error reading Public Cloud images",
			httpResponse,
			err,
			&response.Diagnostics,
			ctx,
		)

		return
	}

	tflog.Info(ctx, fmt.Sprintf("Create publiccloud image resource for %q", state.ID.ValueString()))
	instance, resourceErr := adaptImageDetailsToImageResource(ctx, *sdkImage)
	if resourceErr != nil {
		response.Diagnostics.AddError("Error creating publiccloud image resource", resourceErr.Error())

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
	var plan imageResourceModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Update publiccloud image %q", plan.ID.ValueString()))
	opts := plan.GetUpdateImageOpts()

	sdkImageDetails, httpResponse, err := i.client.UpdateImage(
		ctx,
		plan.ID.ValueString(),
	).UpdateImageOpts(opts).Execute()
	if err != nil {
		utils.HandleSdkError(
			"Error updating Public Cloud image",
			httpResponse,
			err,
			&response.Diagnostics,
			ctx,
		)

		return
	}

	diags = response.State.Set(ctx, sdkImageDetails)
	response.Diagnostics.Append(diags...)
}

// Delete does nothing as there is no endpoint to delete an Image.
func (i *imageResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
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

	i.client = coreClient.PublicCloudAPI
}

func NewImageResource() resource.Resource {
	return &imageResource{}
}
