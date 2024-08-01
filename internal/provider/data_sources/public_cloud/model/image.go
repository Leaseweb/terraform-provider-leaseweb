package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
