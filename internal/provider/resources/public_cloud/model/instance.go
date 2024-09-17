package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Instance struct {
	Id                  types.String `tfsdk:"id"`
	Region              types.Object `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Resources           types.Object `tfsdk:"resources"`
	Image               types.Object `tfsdk:"image"`
	State               types.String `tfsdk:"state"`
	ProductType         types.String `tfsdk:"product_type"`
	HasPublicIpv4       types.Bool   `tfsdk:"has_public_ipv4"`
	HasPrivateNetwork   types.Bool   `tfsdk:"has_private_network"`
	Type                types.Object `tfsdk:"type"`
	RootDiskSize        types.Int64  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	Ips                 types.List   `tfsdk:"ips"`
	StartedAt           types.String `tfsdk:"started_at"`
	Contract            types.Object `tfsdk:"contract"`
	MarketAppId         types.String `tfsdk:"market_app_id"`
	AutoScalingGroup    types.Object `tfsdk:"auto_scaling_group"`
	Iso                 types.Object `tfsdk:"iso"`
	PrivateNetwork      types.Object `tfsdk:"private_network"`
	// TODO Enable SSH key support
	//SshKey              types.String `tfsdk:"ssh_key"`
}

func (i Instance) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": types.StringType,
		"region": types.ObjectType{
			AttrTypes: Region{}.AttributeTypes(),
		},
		"reference": types.StringType,
		"resources": types.ObjectType{
			AttrTypes: Resources{}.AttributeTypes(),
		},
		"image": types.ObjectType{
			AttrTypes: Image{}.AttributeTypes(),
		},
		"state":               types.StringType,
		"product_type":        types.StringType,
		"has_public_ipv4":     types.BoolType,
		"has_private_network": types.BoolType,
		"type": types.ObjectType{
			AttrTypes: InstanceType{}.AttributeTypes(),
		},
		"root_disk_size":         types.Int64Type,
		"root_disk_storage_type": types.StringType,
		"ips": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: Ip{}.AttributeTypes(),
			},
		},
		"started_at": types.StringType,
		"contract": types.ObjectType{
			AttrTypes: Contract{}.AttributeTypes(),
		},
		"market_app_id": types.StringType,
		"auto_scaling_group": types.ObjectType{
			AttrTypes: AutoScalingGroup{}.AttributeTypes(),
		},
		"iso": types.ObjectType{
			AttrTypes: Iso{}.AttributeTypes(),
		},
		"private_network": types.ObjectType{
			AttrTypes: PrivateNetwork{}.AttributeTypes(),
		},
		// TODO Enable SSH key support
		//"ssh_key": types.StringType,
	}
}
