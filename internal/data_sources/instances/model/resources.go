package model

import "github.com/leaseweb/leaseweb-go-sdk/publicCloud"

type resources struct {
	Cpu                 cpu                 `tfsdk:"cpu"`
	Memory              memory              `tfsdk:"memory"`
	PublicNetworkSpeed  publicNetworkSpeed  `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed privateNetworkSpeed `tfsdk:"private_network_speed"`
}

func newResources(sdkResources publicCloud.InstanceResources) resources {
	return resources{
		Cpu:                 newCpu(sdkResources.GetCpu()),
		Memory:              newMemory(sdkResources.GetMemory()),
		PublicNetworkSpeed:  newPublicNetworkSpeed(sdkResources.GetPublicNetworkSpeed()),
		PrivateNetworkSpeed: newPrivateNetworkSpeed(sdkResources.GetPrivateNetworkSpeed()),
	}
}
