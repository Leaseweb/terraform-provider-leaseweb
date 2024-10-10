package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Instance struct {
	Id                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Image               types.Object `tfsdk:"image"`
	State               types.String `tfsdk:"state"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int64  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	Ips                 types.List   `tfsdk:"ips"`
	Contract            types.Object `tfsdk:"contract"`
	MarketAppId         types.String `tfsdk:"market_app_id"`
	// TODO Enable SSH key support
	//SshKey              types.String `tfsdk:"ssh_key"`
}

func (i Instance) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"region":    types.StringType,
		"reference": types.StringType,
		"image": types.ObjectType{
			AttrTypes: Image{}.AttributeTypes(),
		},
		"state":                  types.StringType,
		"type":                   types.StringType,
		"root_disk_size":         types.Int64Type,
		"root_disk_storage_type": types.StringType,
		"ips": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: Ip{}.AttributeTypes(),
			},
		},
		"contract": types.ObjectType{
			AttrTypes: Contract{}.AttributeTypes(),
		},
		"market_app_id": types.StringType,
		// TODO Enable SSH key support
		//"ssh_key": types.StringType,
	}
}
