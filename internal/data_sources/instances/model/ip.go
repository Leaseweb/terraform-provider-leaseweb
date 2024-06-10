package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
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
		Ip:            utils.GenerateString(sdkIp.HasIp(), sdkIp.GetIp()),
		PrefixLength:  utils.GenerateString(sdkIp.HasPrefixLength(), sdkIp.GetPrefixLength()),
		Version:       utils.GenerateInt(sdkIp.HasVersion(), sdkIp.GetVersion()),
		NullRouted:    utils.GenerateBool(sdkIp.HasNullRouted(), sdkIp.GetNullRouted()),
		MainIp:        utils.GenerateBool(sdkIp.HasMainIp(), sdkIp.GetMainIp()),
		NetworkType:   utils.GenerateString(sdkIp.HasNetworkType(), string(sdkIp.GetNetworkType())),
		ReverseLookup: utils.GenerateString(sdkIp.HasReverseLookup(), sdkIp.GetReverseLookup()),
		Ddos:          newDdos(sdkIp.GetDdos()),
	}
}
