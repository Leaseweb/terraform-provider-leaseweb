package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type FeatureAvailability struct {
	Automation       types.Bool `tfsdk:"automation"`
	IpmiReboot       types.Bool `tfsdk:"ipmi_reboot"`
	PowerCycle       types.Bool `tfsdk:"power_cycle"`
	PrivateNetwork   types.Bool `tfsdk:"private_network"`
	RemoteManagement types.Bool `tfsdk:"remote_management"`
}
