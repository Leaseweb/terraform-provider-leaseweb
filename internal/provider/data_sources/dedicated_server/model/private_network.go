package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetwork struct {
	Id        types.String `tfsdk:"id"`
	LinkSpeed types.Int32  `tfsdk:"link_speed"`
	Status    types.String `tfsdk:"status"`
	Subnet    types.String `tfsdk:"subnet"`
	VlanId    types.String `tfsdk:"vlan_id"`
}
