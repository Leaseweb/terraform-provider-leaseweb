package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Resources struct {
	Cpu                 types.Object `tfsdk:"cpu"`
	Memory              types.Object `tfsdk:"memory"`
	PublicNetworkSpeed  types.Object `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed types.Object `tfsdk:"private_network_speed"`
}

func (r Resources) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cpu":                   types.ObjectType{AttrTypes: Cpu{}.AttributeTypes()},
		"memory":                types.ObjectType{AttrTypes: Memory{}.AttributeTypes()},
		"public_network_speed":  types.ObjectType{AttrTypes: NetworkSpeed{}.AttributeTypes()},
		"private_network_speed": types.ObjectType{AttrTypes: NetworkSpeed{}.AttributeTypes()},
	}
}
