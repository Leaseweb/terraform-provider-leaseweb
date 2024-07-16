package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
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

func newIp(
	ctx context.Context,
	entityIp entity.Ip,
) (*Ip, diag.Diagnostics) {
	ddosObject, diags := utils.ConvertNullableDomainEntityToResourceObject(
		entityIp.Ddos,
		Ddos{}.AttributeTypes(),
		ctx,
		newDdos,
	)

	if diags != nil {
		return nil, diags
	}

	return &Ip{
		Ip:            basetypes.NewStringValue(entityIp.Ip),
		PrefixLength:  basetypes.NewStringValue(entityIp.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(entityIp.Version)),
		NullRouted:    basetypes.NewBoolValue(entityIp.NullRouted),
		MainIp:        basetypes.NewBoolValue(entityIp.MainIp),
		NetworkType:   basetypes.NewStringValue(string(entityIp.NetworkType)),
		ReverseLookup: utils.ConvertNullableStringToStringValue(entityIp.ReverseLookup),
		Ddos:          ddosObject,
	}, nil
}
