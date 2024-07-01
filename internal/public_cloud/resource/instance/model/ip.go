package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
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

func newIp(ctx context.Context, sdkIp *publicCloud.Ip) (Ip, diag.Diagnostics) {
	ddosObject, diags := generateDdos(*sdkIp, ctx)

	if diags != nil {
		return Ip{}, diags
	}

	return Ip{
		Ip:            basetypes.NewStringValue(sdkIp.GetIp()),
		PrefixLength:  basetypes.NewStringValue(sdkIp.GetPrefixLength()),
		Version:       basetypes.NewInt64Value(int64(sdkIp.GetVersion())),
		NullRouted:    basetypes.NewBoolValue(sdkIp.GetNullRouted()),
		MainIp:        basetypes.NewBoolValue(sdkIp.GetMainIp()),
		NetworkType:   basetypes.NewStringValue(string(sdkIp.GetNetworkType())),
		ReverseLookup: basetypes.NewStringValue(sdkIp.GetReverseLookup()),
		Ddos:          ddosObject,
	}, nil
}

func generateDdos(
	ip publicCloud.Ip,
	ctx context.Context,
) (basetypes.ObjectValue, diag.Diagnostics) {
	ddos := newDdos(ip.Ddos.Get())
	ddosObject, diags := types.ObjectValueFrom(ctx, ddos.AttributeTypes(), ddos)

	if diags != nil {
		return ddosObject, diags
	}

	return ddosObject, nil
}
