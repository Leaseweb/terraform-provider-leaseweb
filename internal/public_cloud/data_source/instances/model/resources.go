package model

import "github.com/leaseweb/leaseweb-go-sdk/publicCloud"

type resources struct {
	Cpu                 cpu          `tfsdk:"cpu"`
	Memory              memory       `tfsdk:"memory"`
	PublicNetworkSpeed  NetworkSpeed `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed NetworkSpeed `tfsdk:"private_network_speed"`
}

func newResources(sdkResources publicCloud.Resources) resources {
	return resources{
		Cpu:                 newCpu(sdkResources.GetCpu()),
		Memory:              newMemory(sdkResources.GetMemory()),
		PublicNetworkSpeed:  newNetworkSpeed(sdkResources.GetPublicNetworkSpeed()),
		PrivateNetworkSpeed: newNetworkSpeed(sdkResources.GetPrivateNetworkSpeed()),
	}
}
