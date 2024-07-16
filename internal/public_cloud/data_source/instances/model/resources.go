package model

import (
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type resources struct {
	Cpu                 cpu          `tfsdk:"cpu"`
	Memory              memory       `tfsdk:"memory"`
	PublicNetworkSpeed  NetworkSpeed `tfsdk:"public_network_speed"`
	PrivateNetworkSpeed NetworkSpeed `tfsdk:"private_network_speed"`
}

func newResources(entityResources entity.Resources) resources {
	return resources{
		Cpu:                 newCpu(entityResources.Cpu),
		Memory:              newMemory(entityResources.Memory),
		PublicNetworkSpeed:  newNetworkSpeed(entityResources.PublicNetworkSpeed),
		PrivateNetworkSpeed: newNetworkSpeed(entityResources.PrivateNetworkSpeed),
	}
}
