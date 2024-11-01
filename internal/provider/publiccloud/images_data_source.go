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
	_ datasource.DataSourceWithConfigure = &imagesDataSource{}
)

type imageModelDataSource struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Custom       types.Bool   `tfsdk:"custom"`
	State        types.String `tfsdk:"state"`
	MarketApps   []string     `tfsdk:"market_apps"`
	StorageTypes []string     `tfsdk:"storage_types"`
	Flavour      types.String `tfsdk:"flavour"`
	Region       types.String `tfsdk:"region"`
}

func adaptImageToImageDataSource(sdkImage publicCloud.Image) imageModelDataSource {
	return imageModelDataSource{
		ID:      basetypes.NewStringValue(sdkImage.GetId()),
		Name:    basetypes.NewStringValue(sdkImage.GetName()),
		Custom:  basetypes.NewBoolValue(sdkImage.GetCustom()),
		Flavour: basetypes.NewStringValue(string(sdkImage.GetFlavour())),
	}
}

func adaptImageDetailsToImageDataSource(
	sdkImageDetails publicCloud.ImageDetails,
) imageModelDataSource {
	var marketApps []string
	var storageTypes []string

	for _, marketApp := range sdkImageDetails.GetMarketApps() {
		marketApps = append(marketApps, string(marketApp))
	}

	for _, storageType := range sdkImageDetails.GetStorageTypes() {
		storageTypes = append(storageTypes, string(storageType))
	}

	return imageModelDataSource{
		ID:           basetypes.NewStringValue(sdkImageDetails.GetId()),
		Name:         basetypes.NewStringValue(sdkImageDetails.GetName()),
		Custom:       basetypes.NewBoolValue(sdkImageDetails.GetCustom()),
		State:        basetypes.NewStringValue(string(sdkImageDetails.GetState())),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(sdkImageDetails.GetFlavour())),
		Region:       basetypes.NewStringValue(string(sdkImageDetails.GetRegion())),
	}
}

type imagesDataSourceModel struct {
	Images []imageModelDataSource `tfsdk:"images"`
}

func adaptImagesToImagesDataSource(sdkImages []publicCloud.ImageDetails) imagesDataSourceModel {
	var images imagesDataSourceModel

	for _, sdkImageDetails := range sdkImages {
		image := adaptImageDetailsToImageDataSource(sdkImageDetails)
		images.Images = append(images.Images, image)
	}

	return images
}

func getAllImages(ctx context.Context, api publicCloud.PublicCloudAPI) (
	[]publicCloud.ImageDetails,
	*utils.SdkError,
) {
	var images []publicCloud.ImageDetails
	var offset *int32

	request := api.GetImageList(ctx)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, utils.NewSdkError("getAllImages", err, response)
		}

		images = append(images, result.GetImages()...)

		metadata := result.GetMetadata()

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		request.Offset(*offset)
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

type imagesDataSource struct {
	name   string
	client publicCloud.PublicCloudAPI
}

func (i *imagesDataSource) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, i.name)
}

func (i *imagesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
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

func (i *imagesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	tflog.Info(ctx, "Read publiccloud images")
	images, err := getAllImages(ctx, i.client)

	if err != nil {
		summary := fmt.Sprintf("Reading data %s", i.name)
		// TODO: for the error details,
		// the implementation of method getAllImages need to be change
		response.Diagnostics.AddError(summary, err.Error())

		utils.LogError(
			ctx,
			err.ErrorResponse,
			&response.Diagnostics,
			summary,
			err.Error(),
		)

		return
	}

	state := adaptImagesToImagesDataSource(images)

	diags := response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (i *imagesDataSource) Configure(
	_ context.Context,
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

	i.client = coreClient.PublicCloudAPI
}

func NewImagesDataSource() datasource.DataSource {
	return &imagesDataSource{
		name: "public_cloud_images",
	}
}
