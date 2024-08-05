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
	State        types.String   `tfsdk:"state"`
	StateReason  types.String   `tfsdk:"state_reason"`
	Region       types.String   `tfsdk:"region"`
	CreatedAt    types.String   `tfsdk:"created_at"`
	UpdatedAt    types.String   `tfsdk:"updated_at"`
	Custom       types.Bool     `tfsdk:"custom"`
	Architecture types.String   `tfsdk:"architecture"`
	StorageSize  *StorageSize   `tfsdk:"storage_size"`
	MarketApps   []types.String `tfsdk:"market_apps"`
	StorageTypes []types.String `tfsdk:"storage_types"`
}
