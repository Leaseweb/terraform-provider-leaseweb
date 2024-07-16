package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
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

func newImage(entityImage entity.Image) image {
	image := image{
		Id:           basetypes.NewStringValue(string(entityImage.Id)),
		Name:         basetypes.NewStringValue(entityImage.Name),
		Version:      basetypes.NewStringValue(entityImage.Version),
		Family:       basetypes.NewStringValue(entityImage.Family),
		Flavour:      basetypes.NewStringValue(entityImage.Flavour),
		Architecture: basetypes.NewStringValue(entityImage.Architecture),
	}

	for _, marketApp := range entityImage.MarketApps {
		image.MarketApps = append(
			image.MarketApps, types.StringValue(marketApp),
		)
	}

	for _, storageType := range entityImage.StorageTypes {
		image.StorageTypes = append(
			image.StorageTypes, types.StringValue(storageType),
		)
	}

	return image
}
