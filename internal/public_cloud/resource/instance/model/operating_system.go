package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type OperatingSystem struct {
	Id           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Version      types.String   `tfsdk:"version"`
	Family       types.String   `tfsdk:"family"`
	Flavour      types.String   `tfsdk:"flavour"`
	Architecture types.String   `tfsdk:"architecture"`
	MarketApps   []types.String `tfsdk:"market_apps"`
	StorageTypes []types.String `tfsdk:"storage_types"`
}

func (o OperatingSystem) attributeTypes() map[string]attr.Type {
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

func newOperatingSystem(sdkOperatingSystem *publicCloud.OperatingSystem) *OperatingSystem {
	var marketApps []types.String
	var storageTypes []types.String

	for _, marketApp := range sdkOperatingSystem.MarketApps {
		marketApps = append(marketApps, types.StringValue(marketApp))
	}

	for _, storageType := range sdkOperatingSystem.StorageTypes {
		storageTypes = append(storageTypes, types.StringValue(storageType))
	}

	return &OperatingSystem{
		Id:           basetypes.NewStringValue(string(sdkOperatingSystem.GetId())),
		Name:         basetypes.NewStringValue(sdkOperatingSystem.GetName()),
		Version:      basetypes.NewStringValue(sdkOperatingSystem.GetVersion()),
		Family:       basetypes.NewStringValue(sdkOperatingSystem.GetFamily()),
		Flavour:      basetypes.NewStringValue(sdkOperatingSystem.GetFlavour()),
		Architecture: basetypes.NewStringValue(sdkOperatingSystem.GetArchitecture()),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
	}
}
