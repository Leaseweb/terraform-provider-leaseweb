package publiccloud

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	_ resource.ResourceWithModifyPlan  = &imageResource{}
	_ validator.String                 = instanceIdValidator{}
)

type resourceModelImage struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Custom       types.Bool   `tfsdk:"custom"`
	State        types.String `tfsdk:"state"`
	MarketApps   types.List   `tfsdk:"market_apps"`
	StorageTypes types.List   `tfsdk:"storage_types"`
	Flavour      types.String `tfsdk:"flavour"`
	Region       types.String `tfsdk:"region"`
}

func (i resourceModelImage) AttributeTypes() map[string]attr.Type {
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

func (i resourceModelImage) GetUpdateImageOpts() publicCloud.UpdateImageOpts {
	return publicCloud.UpdateImageOpts{
		Name: i.Name.ValueString(),
	}
}

func (i resourceModelImage) GetCreateImageOpts() publicCloud.CreateImageOpts {
	return publicCloud.CreateImageOpts{
		Name:       i.Name.ValueString(),
		InstanceId: i.Id.ValueString(),
	}
}

func newResourceModelImageFromImageDetails(
	ctx context.Context,
	sdkImageDetails publicCloud.ImageDetails,
) (*resourceModelImage, error) {
	marketApps, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		sdkImageDetails.MarketApps,
	)
	if diags.HasError() {
		return nil, fmt.Errorf(
			diags.Errors()[0].Summary(),
			diags.Errors()[0].Detail(),
		)
	}

	storageTypes, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		sdkImageDetails.StorageTypes,
	)
	if diags.HasError() {
		return nil, fmt.Errorf(
			diags.Errors()[0].Summary(),
			diags.Errors()[0].Detail(),
		)
	}

	image := resourceModelImage{
		Id:           basetypes.NewStringValue(sdkImageDetails.Id),
		Name:         basetypes.NewStringValue(sdkImageDetails.Name),
		Custom:       basetypes.NewBoolValue(sdkImageDetails.Custom),
		State:        utils.AdaptNullableStringToStringValue(sdkImageDetails.State.Get()),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(sdkImageDetails.Flavour),
		Region:       utils.AdaptNullableStringEnumToStringValue(sdkImageDetails.Region.Get()),
	}

	return &image, nil
}

func newResourceModelImageFromImage(
	_ context.Context,
	sdkImage publicCloud.Image,
) (*resourceModelImage, error) {
	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})

	return &resourceModelImage{
		Id:           basetypes.NewStringValue(sdkImage.Id),
		Name:         basetypes.NewStringValue(sdkImage.Name),
		Custom:       basetypes.NewBoolValue(sdkImage.Custom),
		Flavour:      basetypes.NewStringValue(sdkImage.Flavour),
		MarketApps:   emptyList,
		StorageTypes: emptyList,
	}, nil
}

type instanceIdValidator struct {
	instanceIds []string
}

func (i instanceIdValidator) Description(_ context.Context) string {
	return `instanceId needs to exist.`
}

func (i instanceIdValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i instanceIdValidator) ValidateString(
	_ context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	// If the instanceId is unknown or null, there is nothing to validate.
	if request.ConfigValue.IsUnknown() || request.ConfigValue.IsNull() {
		return
	}

	instanceIdExists := slices.Contains(i.instanceIds, request.ConfigValue.ValueString())

	if !instanceIdExists {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid id",
			fmt.Sprintf(
				"Attribute id value must be one of: %q, got: %q",
				i.instanceIds,
				request.ConfigValue.ValueString(),
			),
		)
	}
}

func newInstanceIdValidator(sdkInstances []publicCloud.Instance) instanceIdValidator {
	var instanceIds []string

	for _, sdkInstance := range sdkInstances {
		instanceIds = append(instanceIds, sdkInstance.Id)
	}

	return instanceIdValidator{instanceIds: instanceIds}
}

func getImage(
	id string,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.ImageDetails, *utils.SdkError) {
	sdkImages, err := getAllImages(ctx, api)

	if err != nil {
		return nil, err
	}

	for _, sdkImage := range sdkImages {
		if sdkImage.Id == id {
			return &sdkImage, nil
		}
	}

	return nil, nil
}

func updateImage(
	id string,
	opts publicCloud.UpdateImageOpts,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.ImageDetails, *utils.SdkError) {
	image, response, err := api.UpdateImage(
		ctx,
		id,
	).UpdateImageOpts(opts).Execute()
	if err != nil {
		return nil, utils.NewSdkError(
			fmt.Sprintf("updateImage %q", id),
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

// ModifyPlan check that passed id exists.
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

	idRequest := validator.StringRequest{ConfigValue: planImage.Id}
	idResponse := validator.StringResponse{}

	instanceIdValidator := newInstanceIdValidator(instances)
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
				Required:    true,
				Description: "The id of the instance which the custom image is based on",
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

	instance, resourceErr := newResourceModelImageFromImageDetails(ctx, *sdkImage)
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

	sdkImage, err := getImage(state.Id.ValueString(), ctx, i.client.PublicCloudAPI)
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
		state.Id.ValueString(),
	))
	instance, resourceErr := newResourceModelImageFromImageDetails(ctx, *sdkImage)
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
		plan.Id.ValueString(),
	))
	opts := plan.GetUpdateImageOpts()

	sdkImageDetails, sdkErr := updateImage(
		plan.Id.ValueString(),
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
				plan.Id.ValueString(),
			),
			sdkErr.Error(),
		)

		return
	}

	diags = response.State.Set(ctx, sdkImageDetails)
	response.Diagnostics.Append(diags...)
}

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
