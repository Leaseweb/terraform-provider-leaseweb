package publiccloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure = &imageResource{}
)

type imageResourceModel struct {
	ID           types.String `tfsdk:"id"`
	InstanceID   types.String `tfsdk:"instance_id"`
	Name         types.String `tfsdk:"name"`
	Custom       types.Bool   `tfsdk:"custom"`
	State        types.String `tfsdk:"state"`
	MarketApps   types.List   `tfsdk:"market_apps"`
	StorageTypes types.List   `tfsdk:"storage_types"`
	Flavour      types.String `tfsdk:"flavour"`
	Region       types.String `tfsdk:"region"`
}

func (i imageResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"instance_id":   types.StringType,
		"name":          types.StringType,
		"custom":        types.BoolType,
		"state":         types.StringType,
		"market_apps":   types.ListType{ElemType: types.StringType},
		"storage_types": types.ListType{ElemType: types.StringType},
		"flavour":       types.StringType,
		"region":        types.StringType,
	}
}

func (i imageResourceModel) getUpdateImageOpts() publiccloud.UpdateImageOpts {
	return publiccloud.UpdateImageOpts{
		Name: i.Name.ValueString(),
	}
}

func (i imageResourceModel) getCreateImageOpts() publiccloud.CreateImageOpts {
	return publiccloud.CreateImageOpts{
		Name:       i.Name.ValueString(),
		InstanceId: i.InstanceID.ValueString(),
	}
}

func adaptImageDetailsToImageResource(
	ctx context.Context,
	imageDetails publiccloud.ImageDetails,
) (*imageResourceModel, error) {
	marketApps, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		imageDetails.MarketApps,
	)
	if diags.HasError() {
		return nil, fmt.Errorf(diags.Errors()[0].Summary(), diags.Errors()[0].Detail())
	}

	storageTypes, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		imageDetails.StorageTypes,
	)
	if diags.HasError() {
		return nil, fmt.Errorf(diags.Errors()[0].Summary(), diags.Errors()[0].Detail())
	}

	image := imageResourceModel{
		ID:           basetypes.NewStringValue(imageDetails.GetId()),
		Name:         basetypes.NewStringValue(imageDetails.GetName()),
		Custom:       basetypes.NewBoolValue(imageDetails.GetCustom()),
		State:        basetypes.NewStringValue(string(imageDetails.GetState())),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(imageDetails.GetFlavour())),
		Region:       basetypes.NewStringValue(string(imageDetails.GetRegion())),
	}

	return &image, nil
}

func adaptImageToImageResource(image publiccloud.Image) imageResourceModel {
	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})

	return imageResourceModel{
		ID:           basetypes.NewStringValue(image.GetId()),
		Name:         basetypes.NewStringValue(image.GetName()),
		Custom:       basetypes.NewBoolValue(image.GetCustom()),
		Flavour:      basetypes.NewStringValue(string(image.GetFlavour())),
		MarketApps:   emptyList,
		StorageTypes: emptyList,
	}
}

func getImage(
	ID string,
	ctx context.Context,
	api publiccloud.PubliccloudAPI,
) (*publiccloud.ImageDetails, *http.Response, error) {
	images, httpResponse, err := getAllImages(ctx, api)
	if err != nil {
		return nil, httpResponse, err
	}

	for _, image := range images {
		if image.GetId() == ID {
			return &image, nil, nil
		}
	}

	return nil, nil, nil
}

type imageResource struct {
	name   string
	client publiccloud.PubliccloudAPI
}

func (i *imageResource) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf(
		"%s_%s",
		request.ProviderTypeName,
		i.name,
	)
}

func (i *imageResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Can be either an Operating System or a UUID in case of a Custom Image",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				Description: `
The id of the instance which the custom image is based on. The following rules apply:
  - instance exists for instanceId
  - instance has state *STOPPED*
  - instance has a maximum rootDiskSize of 100 GB
  - instance OS must not be *windows*`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

	utils.AddUnsupportedActionsNotation(
		response,
		[]utils.Action{utils.DeleteAction},
	)
}

func (i *imageResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan imageResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := plan.getCreateImageOpts()
	summary := fmt.Sprintf("Creating resource %s", i.name)

	image, httpResponse, err := i.client.CreateImage(ctx).
		CreateImageOpts(opts).
		Execute()
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	state, resourceErr := adaptImageDetailsToImageResource(ctx, *image)
	if resourceErr != nil {
		response.Diagnostics.AddError(summary, utils.DefaultErrMsg)

		return
	}
	// instanceId has to be set manually as it isn't returned from the API
	state.InstanceID = basetypes.NewStringValue(opts.InstanceId)

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (i *imageResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var currentState imageResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &currentState)...)
	if response.Diagnostics.HasError() {
		return
	}

	summary := fmt.Sprintf("Reading resource %s", i.name)

	image, httpResponse, err := getImage(
		currentState.ID.ValueString(),
		ctx,
		i.client,
	)
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	state, resourceErr := adaptImageDetailsToImageResource(ctx, *image)
	if resourceErr != nil {
		response.Diagnostics.AddError(summary, utils.DefaultErrMsg)

		return
	}

	// instanceId has to be set manually as it isn't returned from the API
	state.InstanceID = currentState.InstanceID

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (i *imageResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan imageResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := plan.getUpdateImageOpts()

	imageDetails, httpResponse, err := i.client.UpdateImage(
		ctx,
		plan.ID.ValueString(),
	).UpdateImageOpts(opts).Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Updating resource %s for id %q",
			i.name,
			plan.ID.ValueString(),
		)
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	state, err := adaptImageDetailsToImageResource(ctx, *imageDetails)
	if err != nil {
		summary := fmt.Sprintf("Reading resource %s", i.name)
		utils.Error(ctx, &response.Diagnostics, summary, err, nil)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
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

	i.client = coreClient.PubliccloudAPI
}

func NewImageResource() resource.Resource {
	return &imageResource{
		name: "public_cloud_image",
	}
}
