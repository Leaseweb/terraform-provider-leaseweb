package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Image struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Version      types.String `tfsdk:"version"`
	Family       types.String `tfsdk:"family"`
	Flavour      types.String `tfsdk:"flavour"`
	State        types.String `tfsdk:"state"`
	StateReason  types.String `tfsdk:"state_reason"`
	Region       types.Object `tfsdk:"region"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
	Custom       types.Bool   `tfsdk:"custom"`
	Architecture types.String `tfsdk:"architecture"`
	MarketApps   types.List   `tfsdk:"market_apps"`
	StorageTypes types.List   `tfsdk:"storage_types"`
	StorageSize  types.Object `tfsdk:"storage_size"`
}

func (i Image) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.StringType,
		"name":         types.StringType,
		"version":      types.StringType,
		"family":       types.StringType,
		"flavour":      types.StringType,
		"architecture": types.StringType,
		"state":        types.StringType,
		"state_reason": types.StringType,
		"region": types.ObjectType{
			AttrTypes: Region{}.AttributeTypes(),
		},
		"created_at": types.StringType,
		"updated_at": types.StringType,
		"custom":     types.BoolType,
		"storage_size": types.ObjectType{
			AttrTypes: StorageSize{}.AttributeTypes(),
		},
		"market_apps":   types.ListType{ElemType: types.StringType},
		"storage_types": types.ListType{ElemType: types.StringType},
	}
}
