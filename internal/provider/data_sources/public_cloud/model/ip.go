package model

import (
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
}
