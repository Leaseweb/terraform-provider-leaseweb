package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &ImagesDataSource{}
)

type dataSourceModelImage struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Custom       types.Bool   `tfsdk:"custom"`
	State        types.String `tfsdk:"state"`
	MarketApps   []string     `tfsdk:"market_apps"`
	StorageTypes []string     `tfsdk:"storage_types"`
	Flavour      types.String `tfsdk:"flavour"`
	Region       types.String `tfsdk:"region"`
}

func newDataSourceModelImageFromImage(sdkImage publicCloud.Image) dataSourceModelImage {
	return dataSourceModelImage{
		Id:      basetypes.NewStringValue(sdkImage.Id),
		Name:    basetypes.NewStringValue(sdkImage.Name),
		Custom:  basetypes.NewBoolValue(sdkImage.Custom),
		Flavour: basetypes.NewStringValue(string(sdkImage.Flavour)),
	}
}

func newDataSourceModelImageFromImageDetails(
	sdkImageDetails publicCloud.ImageDetails,
) dataSourceModelImage {
	var marketApps []string
	var storageTypes []string

	for _, marketApp := range sdkImageDetails.MarketApps {
		marketApps = append(marketApps, string(marketApp))
	}

	for _, storageType := range sdkImageDetails.StorageTypes {
		storageTypes = append(storageTypes, string(storageType))
	}

	return dataSourceModelImage{
		Id:           basetypes.NewStringValue(sdkImageDetails.Id),
		Name:         basetypes.NewStringValue(sdkImageDetails.Name),
		Custom:       basetypes.NewBoolValue(sdkImageDetails.Custom),
		State:        utils.AdaptNullableStringEnumToStringValue(sdkImageDetails.State.Get()),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(sdkImageDetails.Flavour)),
		Region:       utils.AdaptNullableStringEnumToStringValue(sdkImageDetails.Region.Get()),
	}
}

type dataSourceModelImages struct {
	Images []dataSourceModelImage `tfsdk:"images"`
}

func newDataSourceModelImages(sdkImages []publicCloud.ImageDetails) dataSourceModelImages {
	var images dataSourceModelImages

	for _, sdkImageDetails := range sdkImages {
		image := newDataSourceModelImageFromImageDetails(sdkImageDetails)
		images.Images = append(images.Images, image)
	}

	return images
}

func getAllImages(ctx context.Context, api publicCloud.PublicCloudAPI) (
	[]publicCloud.ImageDetails,
	*utils.SdkError,
) {
	var images []publicCloud.ImageDetails

	request := api.GetImageList(ctx)

	result, response, err := request.Execute()

	if err != nil {
		return nil, utils.NewSdkError("getAllImages", err, response)
	}

	metadata := result.GetMetadata()
	pagination := utils.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, utils.NewSdkError("getAllImages", err, response)
		}

		images = append(images, result.Images...)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, utils.NewSdkError("getAllImages", err, response)
		}
	}

	return images, nil
}

func imageSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Can be either an Operating System or a UUID in case of a Custom Image",
		},
		"name": schema.StringAttribute{
			Computed: true,
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
	}
}

type ImagesDataSource struct {
	client client.Client
}

func (i *ImagesDataSource) Metadata(
	ctx context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = request.ProviderTypeName + "_public_cloud_images"
}

func (i *ImagesDataSource) Schema(
	ctx context.Context,
	request datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: imageSchemaAttributes(),
				},
			},
		},
	}
}

func (i *ImagesDataSource) Read(
	ctx context.Context,
	request datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	tflog.Info(ctx, "Read public cloud images")
	images, err := getAllImages(ctx, i.client.PublicCloudAPI)

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

	state := newDataSourceModelImages(images)

	diags := response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (i *ImagesDataSource) Configure(
	ctx context.Context,
	request datasource.ConfigureRequest,
	response *datasource.ConfigureResponse,
) {
	if request.ProviderData == nil {
		return
	}

	coreClient, ok := request.ProviderData.(client.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected provider.Client, got: %T. Please report this issue to the provider developers.",
				request.ProviderData,
			),
		)

		return
	}

	i.client = coreClient
}

func NewImagesDataSource() datasource.DataSource {
	return &ImagesDataSource{}
}
