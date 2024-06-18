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

func newOperatingSystem(sdkOperatingSystem publicCloud.OperatingSystem) operatingSystem {
	operatingSystem := operatingSystem{
		Id:           basetypes.NewStringValue(string(sdkOperatingSystem.GetId())),
		Name:         basetypes.NewStringValue(sdkOperatingSystem.GetName()),
		Version:      basetypes.NewStringValue(sdkOperatingSystem.GetVersion()),
		Family:       basetypes.NewStringValue(sdkOperatingSystem.GetFamily()),
		Flavour:      basetypes.NewStringValue(sdkOperatingSystem.GetFlavour()),
		Architecture: basetypes.NewStringValue(sdkOperatingSystem.GetArchitecture()),
	}

	for _, marketApp := range sdkOperatingSystem.MarketApps {
		operatingSystem.MarketApps = append(
			operatingSystem.MarketApps, types.StringValue(marketApp),
		)
	}

	for _, storageType := range sdkOperatingSystem.StorageTypes {
		operatingSystem.StorageTypes = append(
			operatingSystem.StorageTypes, types.StringValue(storageType),
		)
	}

	return operatingSystem
}
