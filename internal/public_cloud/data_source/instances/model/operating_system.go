package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type operatingSystem struct {
	Id           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Version      types.String   `tfsdk:"version"`
	Family       types.String   `tfsdk:"family"`
	Flavour      types.String   `tfsdk:"flavour"`
	Architecture types.String   `tfsdk:"architecture"`
	MarketApps   []types.String `tfsdk:"market_apps"`
	StorageTypes []types.String `tfsdk:"storage_types"`
}

func newOperatingSystem(sdkOperatingSystemDetails publicCloud.OperatingSystemDetails) operatingSystem {
	operatingSystem := operatingSystem{
		Id:           basetypes.NewStringValue(string(sdkOperatingSystemDetails.GetId())),
		Name:         basetypes.NewStringValue(sdkOperatingSystemDetails.GetName()),
		Version:      basetypes.NewStringValue(sdkOperatingSystemDetails.GetVersion()),
		Family:       basetypes.NewStringValue(sdkOperatingSystemDetails.GetFamily()),
		Flavour:      basetypes.NewStringValue(sdkOperatingSystemDetails.GetFlavour()),
		Architecture: basetypes.NewStringValue(sdkOperatingSystemDetails.GetArchitecture()),
	}

	for _, marketApp := range sdkOperatingSystemDetails.MarketApps {
		operatingSystem.MarketApps = append(
			operatingSystem.MarketApps, types.StringValue(marketApp),
		)
	}

	for _, storageType := range sdkOperatingSystemDetails.StorageTypes {
		operatingSystem.StorageTypes = append(
			operatingSystem.StorageTypes, types.StringValue(storageType),
		)
	}

	return operatingSystem
}
