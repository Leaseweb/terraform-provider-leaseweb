package model

type NetworkInterfaces struct {
	Public           NetworkInterface `tfsdk:"public"`
	Internal         NetworkInterface `tfsdk:"internal"`
	RemoteManagement NetworkInterface `tfsdk:"remote_management"`
}
