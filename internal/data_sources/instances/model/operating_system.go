package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
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
		Id: utils.GenerateString(
			sdkOperatingSystem.HasId(),
			sdkOperatingSystem.GetId(),
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
			sdkOperatingSystem.GetFamily(),
		),
		Flavour: utils.GenerateString(
			sdkOperatingSystem.HasFlavour(),
			sdkOperatingSystem.GetFlavour(),
		),
		Architecture: utils.GenerateString(
			sdkOperatingSystem.HasArchitecture(),
			sdkOperatingSystem.GetArchitecture(),
		),
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
