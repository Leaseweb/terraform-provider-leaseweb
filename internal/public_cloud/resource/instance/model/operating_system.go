package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
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
		Id: utils.GenerateString(
			sdkOperatingSystem.HasId(),
			string(sdkOperatingSystem.GetId()),
		),
		Name: utils.GenerateString(
			sdkOperatingSystem.HasName(),
			sdkOperatingSystem.GetName(),
		),
		Version: utils.GenerateString(
			sdkOperatingSystem.HasVersion(),
			sdkOperatingSystem.GetVersion(),
		),
		Family: utils.GenerateString(
			sdkOperatingSystem.HasFamily(),
			sdkOperatingSystem.GetFamily()),
		Flavour: utils.GenerateString(
			sdkOperatingSystem.HasFlavour(),
			sdkOperatingSystem.GetFlavour(),
		),
		Architecture: utils.GenerateString(
			sdkOperatingSystem.HasArchitecture(),
			sdkOperatingSystem.GetArchitecture()),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
	}
}
