package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetwork struct {
	PrivateNetworkId types.String `tfsdk:"private_network_id"`
	Status           types.String `tfsdk:"status"`
	Subnet           types.String `tfsdk:"subnet"`
}

func (p PrivateNetwork) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"private_network_id": types.StringType,
		"status":             types.StringType,
		"subnet":             types.StringType,
	}
}
