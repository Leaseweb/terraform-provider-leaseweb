package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LoadBalancer struct {
	Id                        types.String `tfsdk:"id"`
	Type                      types.Object `tfsdk:"type"`
	Resources                 types.Object `tfsdk:"resources"`
	Region                    types.String `tfsdk:"region"`
	Reference                 types.String `tfsdk:"reference"`
	State                     types.String `tfsdk:"state"`
	Contract                  types.Object `tfsdk:"contract"`
	StartedAt                 types.String `tfsdk:"started_at"`
	Ips                       types.List   `tfsdk:"ips"`
	LoadBalancerConfiguration types.Object `tfsdk:"load_balancer_configuration"`
	PrivateNetwork            types.Object `tfsdk:"private_network"`
}

func (l LoadBalancer) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"type":       types.ObjectType{AttrTypes: InstanceType{}.AttributeTypes()},
		"resources":  types.ObjectType{AttrTypes: Resources{}.AttributeTypes()},
		"region":     types.StringType,
		"reference":  types.StringType,
		"state":      types.StringType,
		"contract":   types.ObjectType{AttrTypes: Contract{}.AttributeTypes()},
		"started_at": types.StringType,
		"ips": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: Ip{}.AttributeTypes(),
			},
		},
		"load_balancer_configuration": types.ObjectType{
			AttrTypes: LoadBalancerConfiguration{}.AttributeTypes(),
		},
		"private_network": types.ObjectType{
			AttrTypes: PrivateNetwork{}.AttributeTypes(),
		},
	}
}
