package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Ip struct {
	Ip            types.String `tfsdk:"ip"`
	PrefixLength  types.String `tfsdk:"prefix_length"`
	Version       types.Int64  `tfsdk:"version"`
	NullRouted    types.Bool   `tfsdk:"null_routed"`
	MainIp        types.Bool   `tfsdk:"main_ip"`
	NetworkType   types.String `tfsdk:"network_type"`
	ReverseLookup types.String `tfsdk:"reverse_lookup"`
	Ddos          types.Object `tfsdk:"ddos"`
}

func (i Ip) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip":             types.StringType,
		"prefix_length":  types.StringType,
		"version":        types.Int64Type,
		"null_routed":    types.BoolType,
		"main_ip":        types.BoolType,
		"network_type":   types.StringType,
		"reverse_lookup": types.StringType,
		"ddos":           types.ObjectType{AttrTypes: Ddos{}.AttributeTypes()},
	}
}
