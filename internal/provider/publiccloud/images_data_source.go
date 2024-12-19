package publiccloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
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

func adaptImageToImageDataSource(image publiccloud.Image) imageModelDataSource {
	return imageModelDataSource{
		ID:      basetypes.NewStringValue(image.GetId()),
		Name:    basetypes.NewStringValue(image.GetName()),
		Custom:  basetypes.NewBoolValue(image.GetCustom()),
		Flavour: basetypes.NewStringValue(string(image.GetFlavour())),
	}
}

func adaptImageDetailsToImageDataSource(
	imageDetails publiccloud.ImageDetails,
) imageModelDataSource {
	var marketApps []string
	var storageTypes []string

	for _, marketApp := range imageDetails.GetMarketApps() {
		marketApps = append(marketApps, string(marketApp))
	}

	for _, storageType := range imageDetails.GetStorageTypes() {
		storageTypes = append(storageTypes, string(storageType))
	}

	return imageModelDataSource{
		ID:           basetypes.NewStringValue(imageDetails.GetId()),
		Name:         basetypes.NewStringValue(imageDetails.GetName()),
		Custom:       basetypes.NewBoolValue(imageDetails.GetCustom()),
		State:        basetypes.NewStringValue(string(imageDetails.GetState())),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(imageDetails.GetFlavour())),
		Region:       basetypes.NewStringValue(string(imageDetails.GetRegion())),
	}
}

type imagesDataSourceModel struct {
	Images []imageModelDataSource `tfsdk:"images"`
}

func adaptImagesToImagesDataSource(sdkImages []publiccloud.ImageDetails) imagesDataSourceModel {
	var images imagesDataSourceModel

	for _, imageDetails := range sdkImages {
		image := adaptImageDetailsToImageDataSource(imageDetails)
		images.Images = append(images.Images, image)
	}

	return images
}

func getAllImages(
	ctx context.Context,
	api publiccloud.PubliccloudAPI,
) ([]publiccloud.ImageDetails, *http.Response, error) {
	var images []publiccloud.ImageDetails
	var offset *int32

	request := api.GetImageList(ctx)

	for {
		result, httpResponse, err := request.Execute()
		if err != nil {
			return nil, httpResponse, fmt.Errorf("getAllImages: %w", err)
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

		request = request.Offset(*offset)
	}

	return images, nil, nil
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
	utils.DataSourceAPI
}

func (i *imagesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: utils.BetaDescription,
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
	images, httpResponse, err := getAllImages(ctx, i.PubliccloudAPI)

	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	response.Diagnostics.Append(
		response.State.Set(ctx, adaptImagesToImagesDataSource(images))...,
	)
}

func NewImagesDataSource() datasource.DataSource {
	return &imagesDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "public_cloud_images",
		},
	}
}
