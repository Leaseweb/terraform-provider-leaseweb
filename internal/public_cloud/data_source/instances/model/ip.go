package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type ip struct {
	Ip            types.String `tfsdk:"ip"`
	PrefixLength  types.String `tfsdk:"prefix_length"`
	Version       types.Int64  `tfsdk:"version"`
	NullRouted    types.Bool   `tfsdk:"null_routed"`
	MainIp        types.Bool   `tfsdk:"main_ip"`
	NetworkType   types.String `tfsdk:"network_type"`
	ReverseLookup types.String `tfsdk:"reverse_lookup"`
	Ddos          ddos         `tfsdk:"ddos"`
}

func newIp(sdkIp publicCloud.Ip) ip {
	return ip{
		Ip:            basetypes.NewStringValue(sdkIp.GetIp()),
		PrefixLength:  basetypes.NewStringValue(sdkIp.GetPrefixLength()),
		Version:       basetypes.NewInt64Value(int64(sdkIp.GetVersion())),
		NullRouted:    basetypes.NewBoolValue(sdkIp.GetNullRouted()),
		MainIp:        basetypes.NewBoolValue(sdkIp.GetMainIp()),
		NetworkType:   basetypes.NewStringValue(string(sdkIp.GetNetworkType())),
		ReverseLookup: basetypes.NewStringValue(sdkIp.GetReverseLookup()),
		Ddos:          newDdos(sdkIp.GetDdos()),
	}
}
