package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
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

func (i Ip) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip":             types.StringType,
		"prefix_length":  types.StringType,
		"version":        types.Int64Type,
		"null_routed":    types.BoolType,
		"main_ip":        types.BoolType,
		"network_type":   types.StringType,
		"reverse_lookup": types.StringType,
		"ddos":           types.ObjectType{AttrTypes: Ddos{}.attributeTypes()},
	}
}

func newIp(ctx context.Context, sdkIp *publicCloud.Ip) (Ip, diag.Diagnostics) {
	ddos := newDdos(sdkIp.Ddos)

	ddosObject, diags := types.ObjectValueFrom(ctx, ddos.attributeTypes(), ddos)
	if diags != nil {
		return Ip{}, diags
	}

	return Ip{
		Ip:            utils.GenerateString(sdkIp.HasIp(), sdkIp.GetIp()),
		PrefixLength:  utils.GenerateString(sdkIp.HasPrefixLength(), sdkIp.GetPrefixLength()),
		Version:       utils.GenerateInt(sdkIp.HasVersion(), sdkIp.GetVersion()),
		NullRouted:    utils.GenerateBool(sdkIp.HasNullRouted(), sdkIp.GetNullRouted()),
		MainIp:        utils.GenerateBool(sdkIp.HasMainIp(), sdkIp.GetMainIp()),
		NetworkType:   utils.GenerateString(sdkIp.HasNetworkType(), string(sdkIp.GetNetworkType())),
		ReverseLookup: utils.GenerateString(sdkIp.HasReverseLookup(), sdkIp.GetReverseLookup()),
		Ddos:          ddosObject,
	}, nil
}
