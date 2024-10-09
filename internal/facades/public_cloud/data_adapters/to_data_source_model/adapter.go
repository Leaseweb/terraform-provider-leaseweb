// Package to_data_source_model implements adapters to convert domain entities to sdk models.
package to_data_source_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
)

func AdaptInstances(domainInstances public_cloud.Instances) model.Instances {
	var instances model.Instances

	for _, domainInstance := range domainInstances {
		instance := adaptInstance(domainInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func adaptInstance(domainInstance public_cloud.Instance) model.Instance {
	instance := model.Instance{
		Id:     basetypes.NewStringValue(domainInstance.Id),
		Region: basetypes.NewStringValue(domainInstance.Region.String()),
		Reference: shared.AdaptNullableStringToStringValue(
			domainInstance.Reference,
		),
		Resources: adaptResources(
			domainInstance.Resources,
		),
		Image:         adaptImage(domainInstance.Image),
		State:         basetypes.NewStringValue(string(domainInstance.State)),
		ProductType:   basetypes.NewStringValue(domainInstance.ProductType),
		HasPublicIpv4: basetypes.NewBoolValue(domainInstance.HasPublicIpv4),
		HasPrivateNetwork: basetypes.NewBoolValue(
			domainInstance.HasPrivateNetwork,
		),
		Type: basetypes.NewStringValue(domainInstance.Type.String()),
		RootDiskSize: basetypes.NewInt64Value(
			int64(domainInstance.RootDiskSize.Value),
		),
		RootDiskStorageType: basetypes.NewStringValue(
			string(domainInstance.RootDiskStorageType),
		),
		StartedAt: shared.AdaptNullableTimeToStringValue(
			domainInstance.StartedAt,
		),
		Contract: adaptContract(
			domainInstance.Contract,
		),
		MarketAppId: shared.AdaptNullableStringToStringValue(
			domainInstance.MarketAppId,
		),
		AutoScalingGroup: shared.AdaptNullableDomainEntityToDatasourceModel(
			domainInstance.AutoScalingGroup,
			adaptAutoScalingGroup,
		),
		PrivateNetwork: shared.AdaptNullableDomainEntityToDatasourceModel(
			domainInstance.PrivateNetwork,
			adaptPrivateNetwork,
		),
	}

	for _, autoScalingGroupIp := range domainInstance.Ips {
		ip := adaptIp(autoScalingGroupIp)
		instance.Ips = append(instance.Ips, ip)
	}

	return instance
}

func adaptResources(resources public_cloud.Resources) model.Resources {
	return model.Resources{
		Cpu:    adaptCpu(resources.Cpu),
		Memory: adaptMemory(resources.Memory),
		PublicNetworkSpeed: adaptNetworkSpeed(
			resources.PublicNetworkSpeed,
		),
		PrivateNetworkSpeed: adaptNetworkSpeed(
			resources.PrivateNetworkSpeed,
		),
	}
}

func adaptCpu(cpu public_cloud.Cpu) model.Cpu {
	return model.Cpu{
		Value: basetypes.NewInt64Value(int64(cpu.Value)),
		Unit:  basetypes.NewStringValue(cpu.Unit),
	}
}

func adaptMemory(memory public_cloud.Memory) model.Memory {
	return model.Memory{
		Value: basetypes.NewFloat64Value(memory.Value),
		Unit:  basetypes.NewStringValue(memory.Unit),
	}
}

func adaptNetworkSpeed(networkSpeed public_cloud.NetworkSpeed) model.NetworkSpeed {
	return model.NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(networkSpeed.Value)),
		Unit:  basetypes.NewStringValue(networkSpeed.Unit),
	}
}

func adaptImage(domainImage public_cloud.Image) model.Image {
	image := model.Image{
		Id:      basetypes.NewStringValue(domainImage.Id),
		Name:    basetypes.NewStringValue(domainImage.Name),
		Family:  basetypes.NewStringValue(domainImage.Family),
		Flavour: basetypes.NewStringValue(domainImage.Flavour),
		Custom:  shared.AdaptBoolToBoolValue(domainImage.Custom),
	}

	return image
}

func adaptContract(contract public_cloud.Contract) model.Contract {
	return model.Contract{
		BillingFrequency: basetypes.NewInt64Value(
			int64(contract.BillingFrequency),
		),
		Term:       basetypes.NewInt64Value(int64(contract.Term)),
		Type:       basetypes.NewStringValue(string(contract.Type)),
		EndsAt:     shared.AdaptNullableTimeToStringValue(contract.EndsAt),
		RenewalsAt: basetypes.NewStringValue(contract.RenewalsAt.String()),
		CreatedAt:  basetypes.NewStringValue(contract.CreatedAt.String()),
		State:      basetypes.NewStringValue(string(contract.State)),
	}
}

func adaptAutoScalingGroup(autoScalingGroup public_cloud.AutoScalingGroup) *model.AutoScalingGroup {
	return &model.AutoScalingGroup{
		Id:    basetypes.NewStringValue(autoScalingGroup.Id),
		Type:  basetypes.NewStringValue(string(autoScalingGroup.Type)),
		State: basetypes.NewStringValue(string(autoScalingGroup.State)),
		DesiredAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.DesiredAmount,
		),
		Region: basetypes.NewStringValue(autoScalingGroup.Region.String()),
		Reference: basetypes.NewStringValue(
			autoScalingGroup.Reference.String(),
		),
		CreatedAt: basetypes.NewStringValue(
			autoScalingGroup.CreatedAt.String(),
		),
		UpdatedAt: basetypes.NewStringValue(
			autoScalingGroup.UpdatedAt.String(),
		),
		StartsAt: shared.AdaptNullableTimeToStringValue(
			autoScalingGroup.StartsAt,
		),
		EndsAt: shared.AdaptNullableTimeToStringValue(
			autoScalingGroup.EndsAt,
		),
		MinimumAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.MinimumAmount,
		),
		MaximumAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.MaximumAmount,
		),
		CpuThreshold: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.CpuThreshold,
		),
		WarmupTime: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.WarmupTime,
		),
		CooldownTime: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.CooldownTime,
		),
	}
}

func adaptPrivateNetwork(privateNetwork public_cloud.PrivateNetwork) *model.PrivateNetwork {
	return &model.PrivateNetwork{
		Id:     basetypes.NewStringValue(privateNetwork.Id),
		Status: basetypes.NewStringValue(privateNetwork.Status),
		Subnet: basetypes.NewStringValue(privateNetwork.Subnet),
	}
}

func adaptIp(ip public_cloud.Ip) model.Ip {
	return model.Ip{
		Ip:            basetypes.NewStringValue(ip.Ip),
		PrefixLength:  basetypes.NewStringValue(ip.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(ip.Version)),
		NullRouted:    basetypes.NewBoolValue(ip.NullRouted),
		MainIp:        basetypes.NewBoolValue(ip.MainIp),
		NetworkType:   basetypes.NewStringValue(string(ip.NetworkType)),
		ReverseLookup: shared.AdaptNullableStringToStringValue(ip.ReverseLookup),
	}
}
