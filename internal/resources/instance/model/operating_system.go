package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
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
		Id:           resources.GetStringValue(sdkOperatingSystem.HasId(), sdkOperatingSystem.GetId()),
		Name:         resources.GetStringValue(sdkOperatingSystem.HasName(), sdkOperatingSystem.GetName()),
		Version:      resources.GetStringValue(sdkOperatingSystem.HasVersion(), sdkOperatingSystem.GetVersion()),
		Family:       resources.GetStringValue(sdkOperatingSystem.HasFamily(), sdkOperatingSystem.GetFamily()),
		Flavour:      resources.GetStringValue(sdkOperatingSystem.HasFlavour(), sdkOperatingSystem.GetFlavour()),
		Architecture: resources.GetStringValue(sdkOperatingSystem.HasArchitecture(), sdkOperatingSystem.GetArchitecture()),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
	}
}
