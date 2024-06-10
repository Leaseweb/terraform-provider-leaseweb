package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/utils"
)

type PublicNetworkSpeed struct {
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (p PublicNetworkSpeed) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Int64Type,
		"unit":  types.StringType,
	}
}

func newPublicNetworkSpeed(sdkPublicNetworkSpeed *publicCloud.PublicNetworkSpeed) PublicNetworkSpeed {
	return PublicNetworkSpeed{
		Value: utils.GenerateInt(sdkPublicNetworkSpeed.HasValue(), sdkPublicNetworkSpeed.GetValue()),
		Unit:  utils.GenerateString(sdkPublicNetworkSpeed.HasUnit(), sdkPublicNetworkSpeed.GetUnit()),
	}
}
