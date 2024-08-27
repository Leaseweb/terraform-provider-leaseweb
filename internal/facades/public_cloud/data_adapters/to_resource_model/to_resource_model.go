package to_resource_model

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func AdaptInstance(
	instance public_cloud.Instance,
	ctx context.Context,
) (*model.Instance, error) {
	plan := model.Instance{}

	plan.Id = basetypes.NewStringValue(instance.Id)
	plan.Reference = shared.AdaptNullableStringToStringValue(instance.Reference)
	plan.State = basetypes.NewStringValue(string(instance.State))
	plan.ProductType = basetypes.NewStringValue(instance.ProductType)
	plan.HasPublicIpv4 = basetypes.NewBoolValue(instance.HasPublicIpv4)
	plan.HasPrivateNetwork = basetypes.NewBoolValue(instance.HasPrivateNetwork)
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

	// TODO Enable SSH key support
	/**
	  if instance.SshKey != nil {
	  	plan.SshKey = basetypes.NewStringValue(instance.SshKey.String())
	  }
	*/

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

	instanceType, err := shared.AdaptDomainEntityToResourceObject(
		instance.Type,
		model.InstanceType{}.AttributeTypes(),
		ctx,
		adaptInstanceType,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Type = instanceType

	region, err := shared.AdaptDomainEntityToResourceObject(
		instance.Region,
		model.Region{}.AttributeTypes(),
		ctx,
		adaptRegion,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}
	plan.Region = region

	return &plan, nil
}

func adaptImage(
	ctx context.Context,
	image public_cloud.Image,
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

	region, err := shared.AdaptNullableDomainEntityToResourceObject(
		image.Region,
		model.Region{}.AttributeTypes(),
		ctx,
		adaptRegion,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptImage: %w", err)
	}
	plan.Region = region

	plan.Id = basetypes.NewStringValue(image.Id)
	plan.Name = basetypes.NewStringValue(image.Name)
	plan.Version = shared.AdaptNullableStringToStringValue(image.Version)
	plan.Family = basetypes.NewStringValue(image.Family)
	plan.Flavour = basetypes.NewStringValue(image.Flavour)
	plan.Architecture = shared.AdaptNullableStringToStringValue(image.Architecture)
	plan.State = shared.AdaptNullableStringToStringValue(image.State)
	plan.StateReason = shared.AdaptNullableStringToStringValue(image.StateReason)
	plan.CreatedAt = shared.AdaptNullableTimeToStringValue(image.CreatedAt)
	plan.UpdatedAt = shared.AdaptNullableTimeToStringValue(image.UpdatedAt)
	plan.Custom = shared.AdaptBoolToBoolValue(image.Custom)

	plan.MarketApps = marketApps
	plan.StorageTypes = storageTypes

	return plan, nil
}

func adaptContract(
	ctx context.Context,
	contract public_cloud.Contract,
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
	iso public_cloud.Iso,
) (*model.Iso, error) {
	return &model.Iso{
		Id:   basetypes.NewStringValue(iso.Id),
		Name: basetypes.NewStringValue(iso.Name),
	}, nil
}

func adaptPrivateNetwork(
	ctx context.Context,
	privateNetwork public_cloud.PrivateNetwork,
) (*model.PrivateNetwork, error) {
	return &model.PrivateNetwork{
		Id:     basetypes.NewStringValue(privateNetwork.Id),
		Status: basetypes.NewStringValue(privateNetwork.Status),
		Subnet: basetypes.NewStringValue(privateNetwork.Subnet),
	}, nil
}

func adaptResources(
	ctx context.Context,
	domainResources public_cloud.Resources,
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
	cpu public_cloud.Cpu,
) (*model.Cpu, error) {
	return &model.Cpu{
		Value: basetypes.NewInt64Value(int64(cpu.Value)),
		Unit:  basetypes.NewStringValue(cpu.Unit),
	}, nil
}

func adaptMemory(
	ctx context.Context,
	memory public_cloud.Memory,
) (*model.Memory, error) {
	return &model.Memory{
		Value: basetypes.NewFloat64Value(memory.Value),
		Unit:  basetypes.NewStringValue(memory.Unit),
	}, nil
}

func adaptNetworkSpeed(
	ctx context.Context,
	networkSpeed public_cloud.NetworkSpeed,
) (*model.NetworkSpeed, error) {

	return &model.NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(networkSpeed.Value)),
		Unit:  basetypes.NewStringValue(networkSpeed.Unit),
	}, nil
}

func adaptAutoScalingGroup(
	ctx context.Context,
	autoScalingGroup public_cloud.AutoScalingGroup,
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

	region, diags := shared.AdaptDomainEntityToResourceObject(
		autoScalingGroup.Region,
		model.Region{}.AttributeTypes(),
		ctx,
		adaptRegion,
	)
	if diags != nil {
		return nil, diags
	}

	return &model.AutoScalingGroup{
		Id:    basetypes.NewStringValue(autoScalingGroup.Id),
		Type:  basetypes.NewStringValue(string(autoScalingGroup.Type)),
		State: basetypes.NewStringValue(string(autoScalingGroup.State)),
		DesiredAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.DesiredAmount,
		),
		Region: region,
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
	loadBalancer public_cloud.LoadBalancer,
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

	instanceType, diags := shared.AdaptDomainEntityToResourceObject(
		loadBalancer.Type,
		model.InstanceType{}.AttributeTypes(),
		ctx,
		adaptInstanceType,
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

	region, diags := shared.AdaptDomainEntityToResourceObject(
		loadBalancer.Region,
		model.Region{}.AttributeTypes(),
		ctx,
		adaptRegion,
	)
	if diags != nil {
		return nil, diags
	}

	return &model.LoadBalancer{
		Id:        basetypes.NewStringValue(loadBalancer.Id),
		Type:      instanceType,
		Resources: resources,
		Region:    region,
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
	configuration public_cloud.LoadBalancerConfiguration,
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
	healthCheck public_cloud.HealthCheck,
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
	stickySession public_cloud.StickySession,
) (*model.StickySession, error) {

	return &model.StickySession{
		Enabled:     basetypes.NewBoolValue(stickySession.Enabled),
		MaxLifeTime: basetypes.NewInt64Value(int64(stickySession.MaxLifeTime)),
	}, nil
}

func adaptIp(
	ctx context.Context,
	ip public_cloud.Ip,
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

func adaptDdos(ctx context.Context, ddos public_cloud.Ddos) (*model.Ddos, error) {
	return &model.Ddos{
		DetectionProfile: basetypes.NewStringValue(ddos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(ddos.ProtectionType),
	}, nil
}

func adaptVolume(
	ctx context.Context,
	volume public_cloud.Volume,
) (*model.Volume, error) {
	return &model.Volume{
		Size: basetypes.NewFloat64Value(volume.Size),
		Unit: basetypes.NewStringValue(volume.Unit),
	}, nil
}

func adaptStorageSize(
	ctx context.Context,
	storageSize public_cloud.StorageSize,
) (*model.StorageSize, error) {
	return &model.StorageSize{
		Size: basetypes.NewFloat64Value(storageSize.Size),
		Unit: basetypes.NewStringValue(storageSize.Unit),
	}, nil
}

func adaptInstanceType(
	ctx context.Context,
	instanceType public_cloud.InstanceType,
) (*model.InstanceType, error) {
	resources, err := shared.AdaptDomainEntityToResourceObject(
		instanceType.Resources,
		model.Resources{}.AttributeTypes(),
		ctx,
		adaptResources,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceType: %w", err)
	}

	prices, err := shared.AdaptDomainEntityToResourceObject(
		instanceType.Prices,
		model.Prices{}.AttributeTypes(),
		ctx,
		adaptPrices,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceType: %w", err)
	}

	storageTypes, storageTypesDiags := basetypes.NewListValueFrom(
		ctx,
		types.StringType,
		instanceType.StorageTypes,
	)
	if storageTypesDiags != nil {
		return nil, shared.ReturnError(
			"adaptInstanceType",
			storageTypesDiags,
		)
	}

	return &model.InstanceType{
		Name:         basetypes.NewStringValue(instanceType.Name),
		Resources:    resources,
		Prices:       prices,
		StorageTypes: storageTypes,
	}, nil
}

func adaptPrices(
	ctx context.Context,
	prices public_cloud.Prices,
) (*model.Prices, error) {
	compute, err := shared.AdaptDomainEntityToResourceObject(
		prices.Compute,
		model.Price{}.AttributeTypes(),
		ctx,
		adaptPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptPrices: %w", err)
	}

	storage, err := shared.AdaptDomainEntityToResourceObject(
		prices.Storage,
		model.Storage{}.AttributeTypes(),
		ctx,
		adaptStorage,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptPrices: %w", err)
	}

	return &model.Prices{
		Currency:       basetypes.NewStringValue(prices.Currency),
		CurrencySymbol: basetypes.NewStringValue(prices.CurrencySymbol),
		Compute:        compute,
		Storage:        storage,
	}, nil
}

func adaptPrice(ctx context.Context, price public_cloud.Price) (*model.Price, error) {
	return &model.Price{
		HourlyPrice:  basetypes.NewStringValue(price.HourlyPrice),
		MonthlyPrice: basetypes.NewStringValue(price.MonthlyPrice),
	}, nil
}

func adaptStorage(
	ctx context.Context,
	storage public_cloud.Storage,
) (*model.Storage, error) {
	local, err := shared.AdaptDomainEntityToResourceObject(
		storage.Local,
		model.Price{}.AttributeTypes(),
		ctx,
		adaptPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptStorage: %w", err)
	}

	central, err := shared.AdaptDomainEntityToResourceObject(
		storage.Central,
		model.Price{}.AttributeTypes(),
		ctx,
		adaptPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptStorage: %w", err)
	}

	return &model.Storage{
		Local:   local,
		Central: central,
	}, nil
}

func adaptRegion(
	ctx context.Context,
	region public_cloud.Region,
) (*model.Region, error) {
	return &model.Region{
		Name:     basetypes.NewStringValue(region.Name),
		Location: basetypes.NewStringValue(region.Location),
	}, nil
}
