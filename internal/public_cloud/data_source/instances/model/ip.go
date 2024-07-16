package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
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
	Ddos          *ddos        `tfsdk:"ddos"`
}

func newIp(entityIp entity.Ip) ip {
	return ip{
		Ip:            basetypes.NewStringValue(entityIp.Ip),
		PrefixLength:  basetypes.NewStringValue(entityIp.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(entityIp.Version)),
		NullRouted:    basetypes.NewBoolValue(entityIp.NullRouted),
		MainIp:        basetypes.NewBoolValue(entityIp.MainIp),
		NetworkType:   basetypes.NewStringValue(string(entityIp.NetworkType)),
		ReverseLookup: utils.ConvertNullableStringToStringValue(entityIp.ReverseLookup),
		Ddos: utils.ConvertNullableDomainEntityToDatasourceModel(
			entityIp.Ddos,
			newDdos,
		),
	}
}
