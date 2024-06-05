package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/resources"
)

type Ip struct {
	Ip            types.String `tfsdk:"ip"`
	PrefixLength  types.String `tfsdk:"prefix_length"`
	Version       types.Int64  `tfsdk:"version"`
	NullRouted    types.Bool   `tfsdk:"null_routed"`
	MainIp        types.Bool   `tfsdk:"main_ip"`
	NetworkType   types.String `tfsdk:"network_type"`
	ReverseLookup types.String `tfsdk:"reverse_lookup"`
	Ddos          Ddos         `tfsdk:"ddos"`
}

func newIp(sdkIp *publicCloud.Ip) Ip {
	return Ip{
		Ip:            resources.GetStringValue(sdkIp.HasIp(), sdkIp.GetIp()),
		PrefixLength:  resources.GetStringValue(sdkIp.HasPrefixLength(), sdkIp.GetPrefixLength()),
		Version:       resources.GetIntValue(sdkIp.HasVersion(), sdkIp.GetVersion()),
		NullRouted:    resources.GetBoolValue(sdkIp.HasNullRouted(), sdkIp.GetNullRouted()),
		MainIp:        resources.GetBoolValue(sdkIp.HasMainIp(), sdkIp.GetMainIp()),
		NetworkType:   resources.GetStringValue(sdkIp.HasNetworkType(), string(sdkIp.GetNetworkType())),
		ReverseLookup: resources.GetStringValue(sdkIp.HasReverseLookup(), sdkIp.GetReverseLookup()),
		Ddos:          newDdos(sdkIp.Ddos),
	}
}
