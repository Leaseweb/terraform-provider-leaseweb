package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Image struct {
	Id           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Version      types.String   `tfsdk:"version"`
	Family       types.String   `tfsdk:"family"`
	Flavour      types.String   `tfsdk:"flavour"`
	Architecture types.String   `tfsdk:"architecture"`
	MarketApps   []types.String `tfsdk:"market_apps"`
	StorageTypes []types.String `tfsdk:"storage_types"`
}

func (i Image) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"name":          types.StringType,
		"version":       types.StringType,
		"family":        types.StringType,
		"flavour":       types.StringType,
		"architecture":  types.StringType,
		"market_apps":   types.ListType{ElemType: types.StringType},
		"storage_types": types.ListType{ElemType: types.StringType},
	}
}

func newImage(
	ctx context.Context,
	sdkImage publicCloud.ImageDetails,
) (*Image, diag.Diagnostics) {
	var marketApps []types.String
	var storageTypes []types.String

	for _, marketApp := range sdkImage.MarketApps {
		marketApps = append(marketApps, types.StringValue(marketApp))
	}

	for _, storageType := range sdkImage.StorageTypes {
		storageTypes = append(storageTypes, types.StringValue(storageType))
	}

	return &Image{
		Id:           basetypes.NewStringValue(string(sdkImage.GetId())),
		Name:         basetypes.NewStringValue(sdkImage.GetName()),
		Version:      basetypes.NewStringValue(sdkImage.GetVersion()),
		Family:       basetypes.NewStringValue(sdkImage.GetFamily()),
		Flavour:      basetypes.NewStringValue(sdkImage.GetFlavour()),
		Architecture: basetypes.NewStringValue(sdkImage.GetArchitecture()),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
	}, nil
}
