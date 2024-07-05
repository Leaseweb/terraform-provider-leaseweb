package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type image struct {
	Id           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Version      types.String   `tfsdk:"version"`
	Family       types.String   `tfsdk:"family"`
	Flavour      types.String   `tfsdk:"flavour"`
	Architecture types.String   `tfsdk:"architecture"`
	MarketApps   []types.String `tfsdk:"market_apps"`
	StorageTypes []types.String `tfsdk:"storage_types"`
}

func newImage(sdkImageDetails publicCloud.ImageDetails) image {
	image := image{
		Id:           basetypes.NewStringValue(string(sdkImageDetails.GetId())),
		Name:         basetypes.NewStringValue(sdkImageDetails.GetName()),
		Version:      basetypes.NewStringValue(sdkImageDetails.GetVersion()),
		Family:       basetypes.NewStringValue(sdkImageDetails.GetFamily()),
		Flavour:      basetypes.NewStringValue(sdkImageDetails.GetFlavour()),
		Architecture: basetypes.NewStringValue(sdkImageDetails.GetArchitecture()),
	}

	for _, marketApp := range sdkImageDetails.MarketApps {
		image.MarketApps = append(
			image.MarketApps, types.StringValue(marketApp),
		)
	}

	for _, storageType := range sdkImageDetails.StorageTypes {
		image.StorageTypes = append(
			image.StorageTypes, types.StringValue(storageType),
		)
	}

	return image
}
