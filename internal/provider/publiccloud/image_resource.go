package publiccloud

import (
	"context"
	"fmt"

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
	_ validator.String                 = instanceIdForCustomImageValidator{}
)

type resourceModelImage struct {
	ID           types.String `tfsdk:"id"`
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
		InstanceId: i.ID.ValueString(),
	}
}

func mapSdkImageDetailsToResourceImage(
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
		ID:           basetypes.NewStringValue(sdkImageDetails.GetId()),
		Name:         basetypes.NewStringValue(sdkImageDetails.GetName()),
		Custom:       basetypes.NewBoolValue(sdkImageDetails.GetCustom()),
		State:        basetypes.NewStringValue(string(sdkImageDetails.GetState())),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(sdkImageDetails.Flavour)),
		Region:       basetypes.NewStringValue(string(sdkImageDetails.GetRegion())),
	}

	return &image, nil
}

func mapSdkImageToResourceImage(
	_ context.Context,
	sdkImage publicCloud.Image,
) (*resourceModelImage, error) {
	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})

	return &resourceModelImage{
		ID:           basetypes.NewStringValue(sdkImage.GetId()),
		Name:         basetypes.NewStringValue(sdkImage.GetName()),
		Custom:       basetypes.NewBoolValue(sdkImage.GetCustom()),
		Flavour:      basetypes.NewStringValue(string(sdkImage.GetFlavour())),
		MarketApps:   emptyList,
		StorageTypes: emptyList,
	}, nil
}

// - Does not yet test
// that the customer has an object storage in the given entity,
// as there's currently no public endpoint for this.
// - Does not check that an image configuration in place.
type instanceIdForCustomImageValidator struct {
	validIds  []string
	instances []publicCloud.Instance
}

func (i instanceIdForCustomImageValidator) Description(_ context.Context) string {
	return `Checks the following:
  - instance exists for instanceId
  - instance has state "STOPPED"
  - instance has a maximum rootDiskSize of 100 GB
  - instance OS must not be Windows
`
}

func (i instanceIdForCustomImageValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i instanceIdForCustomImageValidator) ValidateString(
	_ context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	const maxRootDiskSize = 100

	// If the instanceId is unknown or null, there is nothing to validate.
	if request.ConfigValue.IsUnknown() || request.ConfigValue.IsNull() {
		return
	}

	foundInstance := i.findInstance(request.ConfigValue.ValueString())
	if foundInstance == nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid id",
			fmt.Sprintf(
				"Attribute id value must be one of: %q, got: %q",
				i.validIds,
				request.ConfigValue.ValueString(),
			),
		)

		return
	}

	if foundInstance.GetState() != publicCloud.STATE_STOPPED {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid instance state",
			fmt.Sprintf(
				"Instance linked to attribute ID %q does not have state %q, has state %q",
				request.ConfigValue.ValueString(),
				publicCloud.STATE_STOPPED,
				foundInstance.GetState(),
			),
		)

		return
	}

	if foundInstance.GetRootDiskSize() >= maxRootDiskSize {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid instance rootDiskSize",
			fmt.Sprintf(
				"Instance linked to attribute ID %q has rootDiskSize of %d GB, maximum allowed size is %d GB",
				request.ConfigValue.ValueString(),
				foundInstance.GetRootDiskSize(),
				maxRootDiskSize,
			),
		)

		return
	}

	if foundInstance.Image.GetFlavour() == publicCloud.FLAVOUR_WINDOWS {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid instance OS",
			fmt.Sprintf(
				"Instance linked to attribute ID %q has OS %q, only Linux & BSD are allowed",
				request.ConfigValue.ValueString(),
				foundInstance.Image.GetFlavour(),
			),
		)

		return
	}
}

func (i instanceIdForCustomImageValidator) findInstance(id string) *publicCloud.Instance {
	for _, instance := range i.instances {
		if instance.Id == id {
			return &instance
		}
	}

	return nil
}

func newInstanceIdForCustomImageValidator(instances []publicCloud.Instance) instanceIdForCustomImageValidator {
	var validIds []string

	for _, instance := range instances {
		if instance.GetState() == publicCloud.STATE_STOPPED {
			validIds = append(validIds, instance.Id)
		}
	}

	return instanceIdForCustomImageValidator{
		instances: instances,
		validIds:  validIds,
	}
}

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

	instance, resourceErr := mapSdkImageDetailsToResourceImage(ctx, *sdkImage)
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
	instance, resourceErr := mapSdkImageDetailsToResourceImage(ctx, *sdkImage)
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
