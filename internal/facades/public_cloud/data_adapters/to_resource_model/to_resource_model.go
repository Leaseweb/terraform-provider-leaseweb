package to_resource_model

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func AdaptInstance(
	instance domain.Instance,
	ctx context.Context,
) (*model.Instance, error) {
	plan := model.Instance{}

	plan.Id = basetypes.NewStringValue(instance.Id.String())
	plan.Region = basetypes.NewStringValue(instance.Region)
	plan.Reference = shared.AdaptNullableStringToStringValue(instance.Reference)
	plan.State = basetypes.NewStringValue(string(instance.State))
	plan.ProductType = basetypes.NewStringValue(instance.ProductType)
	plan.HasPublicIpv4 = basetypes.NewBoolValue(instance.HasPublicIpv4)
	plan.HasPrivateNetwork = basetypes.NewBoolValue(instance.HasPrivateNetwork)
	plan.Type = basetypes.NewStringValue(instance.Type.String())
	plan.RootDiskSize = basetypes.NewInt64Value(
		int64(instance.RootDiskSize.Value),
	)
	plan.RootDiskStorageType = basetypes.NewStringValue(
		string(instance.RootDiskStorageType),
	)
	plan.StartedAt = shared.AdaptNullableTimeToStringValue(instance.StartedAt)
	plan.MarketAppId = shared.AdaptNullableStringToStringValue(
		instance.MarketAppId,
	)

	if instance.SshKey != nil {
		plan.SshKey = basetypes.NewStringValue(instance.SshKey.String())
	}

	image, err := shared.AdaptDomainEntityToResourceObject(
		instance.Image,
		model.Image{}.AttributeTypes(),
		ctx,
		adaptImage,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Image = image

	contract, err := shared.AdaptDomainEntityToResourceObject(
		instance.Contract,
		model.Contract{}.AttributeTypes(),
		ctx,
		adaptContract,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Contract = contract

	iso, err := shared.AdaptNullableDomainEntityToResourceObject(
		instance.Iso,
		model.Iso{}.AttributeTypes(),
		ctx,
		adaptIso,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Iso = iso

	privateNetwork, err := shared.AdaptNullableDomainEntityToResourceObject(
		instance.PrivateNetwork,
		model.PrivateNetwork{}.AttributeTypes(),
		ctx,
		adaptPrivateNetwork,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.PrivateNetwork = privateNetwork

	resources, err := shared.AdaptDomainEntityToResourceObject(
		instance.Resources,
		model.Resources{}.AttributeTypes(),
		ctx,
		adaptResources,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Resources = resources

	autoScalingGroup, err := shared.AdaptNullableDomainEntityToResourceObject(
		instance.AutoScalingGroup,
		model.AutoScalingGroup{}.AttributeTypes(),
		ctx,
		adaptAutoScalingGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.AutoScalingGroup = autoScalingGroup

	ips, err := shared.AdaptEntitiesToListValue(
		instance.Ips,
		model.Ip{}.AttributeTypes(),
		ctx,
		adaptIp,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Ips = ips

	volume, err := shared.AdaptNullableDomainEntityToResourceObject(
		instance.Volume,
		model.Volume{}.AttributeTypes(),
		ctx,
		adaptVolume,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Volume = volume

	return &plan, nil
}

func adaptImage(
	ctx context.Context,
	image domain.Image,
) (*model.Image, error) {
	plan := &model.Image{}

	marketApps, marketAppsDiags := basetypes.NewListValueFrom(
		ctx,
		types.StringType,
		image.MarketApps,
	)
	if marketAppsDiags != nil {
		return nil, shared.ReturnError(
			"adaptImage",
			marketAppsDiags,
		)
	}

	storageTypes, storageTypesDiags := basetypes.NewListValueFrom(
		ctx,
		types.StringType,
		image.StorageTypes,
	)
	if storageTypesDiags != nil {
		return nil, shared.ReturnError(
			"adaptImage",
			storageTypesDiags,
		)
	}

	storageSize, err := shared.AdaptNullableDomainEntityToResourceObject(
		image.StorageSize,
		model.StorageSize{}.AttributeTypes(),
		ctx,
		adaptStorageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.StorageSize = storageSize

	plan.Id = basetypes.NewStringValue(image.Id)
	plan.Name = basetypes.NewStringValue(image.Name)
	plan.Version = basetypes.NewStringValue(image.Version)
	plan.Family = basetypes.NewStringValue(image.Family)
	plan.Flavour = basetypes.NewStringValue(image.Flavour)
	plan.Architecture = basetypes.NewStringValue(image.Architecture)
	plan.State = shared.AdaptNullableStringToStringValue(image.State)
	plan.StateReason = shared.AdaptNullableStringToStringValue(image.StateReason)
	plan.Region = shared.AdaptNullableStringToStringValue(image.Region)
	plan.CreatedAt = shared.AdaptNullableTimeToStringValue(image.CreatedAt)
	plan.UpdatedAt = shared.AdaptNullableTimeToStringValue(image.UpdatedAt)
	plan.Custom = shared.AdaptNullableBoolToBoolValue(image.Custom)

	plan.MarketApps = marketApps
	plan.StorageTypes = storageTypes

	return plan, nil
}

func adaptContract(
	ctx context.Context,
	contract domain.Contract,
) (*model.Contract, error) {

	return &model.Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(contract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(contract.Term)),
		Type:             basetypes.NewStringValue(string(contract.Type)),
		EndsAt:           shared.AdaptNullableTimeToStringValue(contract.EndsAt),
		RenewalsAt:       basetypes.NewStringValue(contract.RenewalsAt.String()),
		CreatedAt:        basetypes.NewStringValue(contract.CreatedAt.String()),
		State:            basetypes.NewStringValue(string(contract.State)),
	}, nil
}

func adaptIso(
	ctx context.Context,
	iso domain.Iso,
) (*model.Iso, error) {
	return &model.Iso{
		Id:   basetypes.NewStringValue(iso.Id),
		Name: basetypes.NewStringValue(iso.Name),
	}, nil
}

func adaptPrivateNetwork(
	ctx context.Context,
	privateNetwork domain.PrivateNetwork,
) (*model.PrivateNetwork, error) {
	return &model.PrivateNetwork{
		Id:     basetypes.NewStringValue(privateNetwork.Id),
		Status: basetypes.NewStringValue(privateNetwork.Status),
		Subnet: basetypes.NewStringValue(privateNetwork.Subnet),
	}, nil
}

func adaptResources(
	ctx context.Context,
	domainResources domain.Resources,
) (*model.Resources, error) {
	var resources model.Resources

	cpu, cpuDiags := shared.AdaptDomainEntityToResourceObject(
		domainResources.Cpu,
		model.Cpu{}.AttributeTypes(),
		ctx,
		adaptCpu,
	)
	if cpuDiags != nil {
		return &resources, cpuDiags
	}
	resources.Cpu = cpu

	memory, memoryDiags := shared.AdaptDomainEntityToResourceObject(
		domainResources.Memory,
		model.Memory{}.AttributeTypes(),
		ctx,
		adaptMemory,
	)
	if memoryDiags != nil {
		return &resources, memoryDiags
	}
	resources.Memory = memory

	publicNetworkSpeed, publicNetworkSpeedDiags := shared.AdaptDomainEntityToResourceObject(
		domainResources.PublicNetworkSpeed,
		model.NetworkSpeed{}.AttributeTypes(),
		ctx,
		adaptNetworkSpeed,
	)
	if publicNetworkSpeedDiags != nil {
		return &resources, publicNetworkSpeedDiags
	}
	resources.PublicNetworkSpeed = publicNetworkSpeed

	privateNetworkSpeed, privateNetworkSpeedDiags := shared.AdaptDomainEntityToResourceObject(
		domainResources.PrivateNetworkSpeed,
		model.NetworkSpeed{}.AttributeTypes(),
		ctx,
		adaptNetworkSpeed,
	)
	if privateNetworkSpeedDiags != nil {
		return &resources, privateNetworkSpeedDiags
	}
	resources.PrivateNetworkSpeed = privateNetworkSpeed

	return &resources, nil
}

func adaptCpu(
	ctx context.Context,
	cpu domain.Cpu,
) (*model.Cpu, error) {
	return &model.Cpu{
		Value: basetypes.NewInt64Value(int64(cpu.Value)),
		Unit:  basetypes.NewStringValue(cpu.Unit),
	}, nil
}

func adaptMemory(
	ctx context.Context,
	memory domain.Memory,
) (*model.Memory, error) {
	return &model.Memory{
		Value: basetypes.NewFloat64Value(memory.Value),
		Unit:  basetypes.NewStringValue(memory.Unit),
	}, nil
}

func adaptNetworkSpeed(
	ctx context.Context,
	networkSpeed domain.NetworkSpeed,
) (*model.NetworkSpeed, error) {

	return &model.NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(networkSpeed.Value)),
		Unit:  basetypes.NewStringValue(networkSpeed.Unit),
	}, nil
}

func adaptAutoScalingGroup(
	ctx context.Context,
	autoScalingGroup domain.AutoScalingGroup,
) (*model.AutoScalingGroup, error) {

	loadBalancer, loadBalancerDiags := shared.AdaptNullableDomainEntityToResourceObject(
		autoScalingGroup.LoadBalancer,
		model.LoadBalancer{}.AttributeTypes(),
		ctx,
		adaptLoadBalancer,
	)
	if loadBalancerDiags != nil {
		return nil, loadBalancerDiags
	}
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
		LoadBalancer: loadBalancer,
	}, nil
}

func adaptLoadBalancer(
	ctx context.Context,
	loadBalancer domain.LoadBalancer,
) (*model.LoadBalancer, error) {

	resources, diags := shared.AdaptDomainEntityToResourceObject(
		loadBalancer.Resources,
		model.Resources{}.AttributeTypes(),
		ctx,
		adaptResources,
	)
	if diags != nil {
		return nil, diags
	}

	contract, diags := shared.AdaptDomainEntityToResourceObject(
		loadBalancer.Contract,
		model.Contract{}.AttributeTypes(),
		ctx,
		adaptContract,
	)
	if diags != nil {
		return nil, diags
	}

	configuration, diags := shared.AdaptNullableDomainEntityToResourceObject(
		loadBalancer.Configuration,
		model.LoadBalancerConfiguration{}.AttributeTypes(),
		ctx,
		adaptLoadBalancerConfiguration,
	)
	if diags != nil {
		return nil, diags
	}

	privateNetwork, diags := shared.AdaptNullableDomainEntityToResourceObject(
		loadBalancer.PrivateNetwork,
		model.PrivateNetwork{}.AttributeTypes(),
		ctx,
		adaptPrivateNetwork,
	)
	if diags != nil {
		return nil, diags
	}

	ips, diags := shared.AdaptEntitiesToListValue(
		loadBalancer.Ips,
		model.Ip{}.AttributeTypes(),
		ctx,
		adaptIp,
	)
	if diags != nil {
		return nil, diags
	}

	return &model.LoadBalancer{
		Id: basetypes.NewStringValue(loadBalancer.Id.String()),
		Type: basetypes.NewStringValue(
			loadBalancer.Type.String(),
		),
		Resources: resources,
		Region:    basetypes.NewStringValue(loadBalancer.Region),
		Reference: shared.AdaptNullableStringToStringValue(
			loadBalancer.Reference,
		),
		State: basetypes.NewStringValue(
			loadBalancer.State.String(),
		),
		Contract: contract,
		StartedAt: basetypes.NewStringValue(
			loadBalancer.StartedAt.String(),
		),
		LoadBalancerConfiguration: configuration,
		PrivateNetwork:            privateNetwork,
		Ips:                       ips,
	}, nil
}

func adaptLoadBalancerConfiguration(
	ctx context.Context,
	configuration domain.LoadBalancerConfiguration,
) (*model.LoadBalancerConfiguration, error) {

	healthCheckObject, diags := shared.AdaptNullableDomainEntityToResourceObject(
		configuration.HealthCheck,
		model.HealthCheck{}.AttributeTypes(),
		ctx,
		adaptHealthCheck,
	)
	if diags != nil {
		return nil, diags
	}

	stickySessionObject, diags := shared.AdaptNullableDomainEntityToResourceObject(
		configuration.StickySession,
		model.StickySession{}.AttributeTypes(),
		ctx,
		adaptStickySession,
	)
	if diags != nil {
		return nil, diags
	}

	return &model.LoadBalancerConfiguration{
		Balance:       basetypes.NewStringValue(configuration.Balance.String()),
		HealthCheck:   healthCheckObject,
		StickySession: stickySessionObject,
		XForwardedFor: basetypes.NewBoolValue(configuration.XForwardedFor),
		IdleTimeout:   basetypes.NewInt64Value(int64(configuration.IdleTimeout)),
		TargetPort:    basetypes.NewInt64Value(int64(configuration.TargetPort)),
	}, nil
}

func adaptHealthCheck(
	ctx context.Context,
	healthCheck domain.HealthCheck,
) (*model.HealthCheck, error) {

	return &model.HealthCheck{
		Method: basetypes.NewStringValue(string(healthCheck.Method)),
		Uri:    basetypes.NewStringValue(healthCheck.Uri),
		Host:   shared.AdaptNullableStringToStringValue(healthCheck.Host),
		Port:   basetypes.NewInt64Value(int64(healthCheck.Port)),
	}, nil
}

func adaptStickySession(
	ctx context.Context,
	stickySession domain.StickySession,
) (*model.StickySession, error) {

	return &model.StickySession{
		Enabled:     basetypes.NewBoolValue(stickySession.Enabled),
		MaxLifeTime: basetypes.NewInt64Value(int64(stickySession.MaxLifeTime)),
	}, nil
}

func adaptIp(
	ctx context.Context,
	ip domain.Ip,
) (*model.Ip, error) {

	ddos, diags := shared.AdaptNullableDomainEntityToResourceObject(
		ip.Ddos,
		model.Ddos{}.AttributeTypes(),
		ctx,
		adaptDdos,
	)

	if diags != nil {
		return nil, diags
	}

	return &model.Ip{
		Ip:            basetypes.NewStringValue(ip.Ip),
		PrefixLength:  basetypes.NewStringValue(ip.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(ip.Version)),
		NullRouted:    basetypes.NewBoolValue(ip.NullRouted),
		MainIp:        basetypes.NewBoolValue(ip.MainIp),
		NetworkType:   basetypes.NewStringValue(string(ip.NetworkType)),
		ReverseLookup: shared.AdaptNullableStringToStringValue(ip.ReverseLookup),
		Ddos:          ddos,
	}, nil
}

func adaptDdos(ctx context.Context, ddos domain.Ddos) (*model.Ddos, error) {
	return &model.Ddos{
		DetectionProfile: basetypes.NewStringValue(ddos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(ddos.ProtectionType),
	}, nil
}

func adaptVolume(
	ctx context.Context,
	volume domain.Volume,
) (*model.Volume, error) {
	return &model.Volume{
		Size: basetypes.NewFloat64Value(volume.Size),
		Unit: basetypes.NewStringValue(volume.Unit),
	}, nil
}

func adaptStorageSize(
	ctx context.Context,
	storageSize domain.StorageSize,
) (*model.StorageSize, error) {
	return &model.StorageSize{
		Size: basetypes.NewFloat64Value(storageSize.Size),
		Unit: basetypes.NewStringValue(storageSize.Unit),
	}, nil
}
