package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type NetworkInterface struct {
	Mac        types.String `tfsdk:"mac"`
	Ip         types.String `tfsdk:"ip"`
	Gateway    types.String `tfsdk:"gateway"`
	Ports      Ports        `tfsdk:"ports"`
	NullRouted types.Bool   `tfsdk:"null_routed"`
	LocationId types.String `tfsdk:"location_id"`
}
