package to_data_source_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/handlers/shared"
	"terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
)

func AdaptInstances(domainInstances domain.Instances) model.Instances {
	var instances model.Instances

	for _, domainInstance := range domainInstances {
		instance := adaptInstance(domainInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func adaptInstance(domainInstance domain.Instance) model.Instance {
	instance := model.Instance{
		Id:     basetypes.NewStringValue(domainInstance.Id.String()),
		Region: basetypes.NewStringValue(domainInstance.Region),
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
		Iso: shared.AdaptNullableDomainEntityToDatasourceModel(
			domainInstance.Iso,
			adaptIso,
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

func adaptResources(resources domain.Resources) model.Resources {
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

func adaptCpu(cpu domain.Cpu) model.Cpu {
	return model.Cpu{
		Value: basetypes.NewInt64Value(int64(cpu.Value)),
		Unit:  basetypes.NewStringValue(cpu.Unit),
	}
}

func adaptMemory(memory domain.Memory) model.Memory {
	return model.Memory{
		Value: basetypes.NewFloat64Value(memory.Value),
		Unit:  basetypes.NewStringValue(memory.Unit),
	}
}

func adaptNetworkSpeed(networkSpeed domain.NetworkSpeed) model.NetworkSpeed {
	return model.NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(networkSpeed.Value)),
		Unit:  basetypes.NewStringValue(networkSpeed.Unit),
	}
}

func adaptImage(domainImage domain.Image) model.Image {
	image := model.Image{
		Id:           basetypes.NewStringValue(string(domainImage.Id)),
		Name:         basetypes.NewStringValue(domainImage.Name),
		Version:      basetypes.NewStringValue(domainImage.Version),
		Family:       basetypes.NewStringValue(domainImage.Family),
		Flavour:      basetypes.NewStringValue(domainImage.Flavour),
		Architecture: basetypes.NewStringValue(domainImage.Architecture),
	}

	for _, marketApp := range domainImage.MarketApps {
		image.MarketApps = append(
			image.MarketApps, types.StringValue(marketApp),
		)
	}

	for _, storageType := range domainImage.StorageTypes {
		image.StorageTypes = append(
			image.StorageTypes, types.StringValue(storageType),
		)
	}

	return image
}

func adaptContract(contract domain.Contract) model.Contract {
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

func adaptAutoScalingGroup(autoScalingGroup domain.AutoScalingGroup) *model.AutoScalingGroup {
	return &model.AutoScalingGroup{
		Id:    basetypes.NewStringValue(autoScalingGroup.Id.String()),
		Type:  basetypes.NewStringValue(string(autoScalingGroup.Type)),
		State: basetypes.NewStringValue(string(autoScalingGroup.State)),
		DesiredAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.DesiredAmount,
		),
		Region: basetypes.NewStringValue(autoScalingGroup.Region),
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
		LoadBalancer: shared.AdaptNullableDomainEntityToDatasourceModel(
			autoScalingGroup.LoadBalancer,
			adaptLoadBalancer,
		),
	}
}

func adaptLoadBalancer(loadBalancer domain.LoadBalancer) *model.LoadBalancer {
	var ips []model.Ip
	for _, ip := range loadBalancer.Ips {
		ips = append(ips, adaptIp(ip))
	}

	return &model.LoadBalancer{
		Id:        basetypes.NewStringValue(loadBalancer.Id.String()),
		Type:      basetypes.NewStringValue(loadBalancer.Type.String()),
		Resources: adaptResources(loadBalancer.Resources),
		Region:    basetypes.NewStringValue(loadBalancer.Region),
		Reference: shared.AdaptNullableStringToStringValue(loadBalancer.Reference),
		State:     basetypes.NewStringValue(string(loadBalancer.State)),
		Contract:  adaptContract(loadBalancer.Contract),
		StartedAt: shared.AdaptNullableTimeToStringValue(loadBalancer.StartedAt),
		Ips:       ips,
		LoadBalancerConfiguration: shared.AdaptNullableDomainEntityToDatasourceModel(
			loadBalancer.Configuration,
			adaptLoadBalancerConfiguration,
		),
		PrivateNetwork: shared.AdaptNullableDomainEntityToDatasourceModel(
			loadBalancer.PrivateNetwork,
			adaptPrivateNetwork,
		),
	}
}

func adaptLoadBalancerConfiguration(configuration domain.LoadBalancerConfiguration) *model.LoadBalancerConfiguration {
	return &model.LoadBalancerConfiguration{
		Balance: basetypes.NewStringValue(configuration.Balance.String()),
		HealthCheck: shared.AdaptNullableDomainEntityToDatasourceModel(
			configuration.HealthCheck,
			adaptHealthCheck,
		),
		StickySession: shared.AdaptNullableDomainEntityToDatasourceModel(
			configuration.StickySession,
			adaptStickySession,
		),
		XForwardedFor: basetypes.NewBoolValue(configuration.XForwardedFor),
		IdleTimeout:   basetypes.NewInt64Value(int64(configuration.IdleTimeout)),
		TargetPort:    basetypes.NewInt64Value(int64(configuration.TargetPort)),
	}
}

func adaptHealthCheck(healthCheck domain.HealthCheck) *model.HealthCheck {
	return &model.HealthCheck{
		Method: basetypes.NewStringValue(healthCheck.Method.String()),
		Uri:    basetypes.NewStringValue(healthCheck.Uri),
		Host:   shared.AdaptNullableStringToStringValue(healthCheck.Host),
		Port:   basetypes.NewInt64Value(int64(healthCheck.Port)),
	}
}

func adaptStickySession(stickySession domain.StickySession) *model.StickySession {
	return &model.StickySession{
		Enabled:     basetypes.NewBoolValue(stickySession.Enabled),
		MaxLifeTime: basetypes.NewInt64Value(int64(stickySession.MaxLifeTime)),
	}
}

func adaptPrivateNetwork(privateNetwork domain.PrivateNetwork) *model.PrivateNetwork {
	return &model.PrivateNetwork{
		Id:     basetypes.NewStringValue(privateNetwork.Id),
		Status: basetypes.NewStringValue(privateNetwork.Status),
		Subnet: basetypes.NewStringValue(privateNetwork.Subnet),
	}
}

func adaptIso(iso domain.Iso) *model.Iso {
	return &model.Iso{
		Id:   basetypes.NewStringValue(iso.Id),
		Name: basetypes.NewStringValue(iso.Name),
	}
}

func adaptIp(ip domain.Ip) model.Ip {
	return model.Ip{
		Ip:            basetypes.NewStringValue(ip.Ip),
		PrefixLength:  basetypes.NewStringValue(ip.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(ip.Version)),
		NullRouted:    basetypes.NewBoolValue(ip.NullRouted),
		MainIp:        basetypes.NewBoolValue(ip.MainIp),
		NetworkType:   basetypes.NewStringValue(string(ip.NetworkType)),
		ReverseLookup: shared.AdaptNullableStringToStringValue(ip.ReverseLookup),
		Ddos: shared.AdaptNullableDomainEntityToDatasourceModel(
			ip.Ddos,
			adaptDdos,
		),
	}
}

func adaptDdos(ddos domain.Ddos) *model.Ddos {
	return &model.Ddos{
		DetectionProfile: basetypes.NewStringValue(ddos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(ddos.ProtectionType),
	}
}
