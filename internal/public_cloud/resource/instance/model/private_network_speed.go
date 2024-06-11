package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type PrivateNetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (p PrivateNetworkSpeed) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Int64Type,
		"unit":  types.StringType,
	}
}

func newPrivateNetworkSpeed(sdkPrivateNetworkSpeed *publicCloud.PrivateNetworkSpeed) PrivateNetworkSpeed {
	return PrivateNetworkSpeed{
		Value: utils.GenerateInt(sdkPrivateNetworkSpeed.HasValue(), sdkPrivateNetworkSpeed.GetValue()),
		Unit:  utils.GenerateString(sdkPrivateNetworkSpeed.HasUnit(), sdkPrivateNetworkSpeed.GetUnit()),
	}
}
