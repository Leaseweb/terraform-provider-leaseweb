package model

type Resources struct {
	Cpu                 Cpu          `tfsdk:"cpu"`
	Memory              Memory       `tfsdk:"memory"`
	PublicNetworkSpeed  NetworkSpeed `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed NetworkSpeed `tfsdk:"private_network_speed"`
}
