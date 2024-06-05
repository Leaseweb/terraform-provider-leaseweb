package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Resources struct {
	Cpu                 Cpu                 `tfsdk:"cpu"`
	Memory              Memory              `tfsdk:"memory"`
	PublicNetworkSpeed  PublicNetworkSpeed  `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed PrivateNetworkSpeed `tfsdk:"private_network_speed"`
}

func newResources(sdkResources *publicCloud.InstanceResources) *Resources {
	return &Resources{
		Cpu:                 newCpu(sdkResources.Cpu),
		Memory:              newMemory(sdkResources.Memory),
		PublicNetworkSpeed:  newPublicNetworkSpeed(sdkResources.PublicNetworkSpeed),
		PrivateNetworkSpeed: newPrivateNetworkSpeed(sdkResources.PrivateNetworkSpeed),
	}
}
