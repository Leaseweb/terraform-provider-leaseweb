package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
)

type Image struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Version      types.String `tfsdk:"version"`
	Family       types.String `tfsdk:"family"`
	Flavour      types.String `tfsdk:"flavour"`
	Architecture types.String `tfsdk:"architecture"`
	MarketApps   types.List   `tfsdk:"market_apps"`
	StorageTypes types.List   `tfsdk:"storage_types"`
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
	image domain.Image,
) (*Image, diag.Diagnostics) {
	marketApps, diags := basetypes.NewListValueFrom(
		ctx,
		types.StringType,
		image.MarketApps,
	)
	if diags.HasError() {
		return nil, diags
	}

	storageTypes, diags := basetypes.NewListValueFrom(
		ctx,
		types.StringType,
		image.StorageTypes,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &Image{
		Id:           basetypes.NewStringValue(string(image.Id)),
		Name:         basetypes.NewStringValue(image.Name),
		Version:      basetypes.NewStringValue(image.Version),
		Family:       basetypes.NewStringValue(image.Family),
		Flavour:      basetypes.NewStringValue(image.Flavour),
		Architecture: basetypes.NewStringValue(image.Architecture),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
	}, nil
}
