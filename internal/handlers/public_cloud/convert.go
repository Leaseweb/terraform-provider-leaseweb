package public_cloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/handlers/shared"
	dataSourceModel "terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	resourcesModel "terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func returnError(functionName string, diags diag.Diagnostics) error {
	for _, diagError := range diags {
		return fmt.Errorf(
			"%s: %q %q",
			functionName,
			diagError.Summary(),
			diagError.Detail(),
		)
	}

	return nil
}

func convertInstanceToResourceModel(
	instance domain.Instance,
	ctx context.Context,
) (*resourcesModel.Instance, error) {
	plan := resourcesModel.Instance{}

	plan.Id = basetypes.NewStringValue(instance.Id.String())
	plan.Region = basetypes.NewStringValue(instance.Region)
	plan.Reference = shared.ConvertNullableStringToStringValue(instance.Reference)
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
	plan.StartedAt = shared.ConvertNullableTimeToStringValue(instance.StartedAt)
	plan.MarketAppId = shared.ConvertNullableStringToStringValue(
		instance.MarketAppId,
	)

	if instance.SshKey != nil {
		plan.SshKey = basetypes.NewStringValue(instance.SshKey.String())
	}

	image, err := shared.ConvertDomainEntityToResourceObject(
		instance.Image,
		resourcesModel.Image{}.AttributeTypes(),
		ctx,
		convertImageToResourceModel,
	)
	if err != nil {
		return nil, fmt.Errorf("convertInstanceToResourceModel: %w", err)
	}
	plan.Image = image

	contract, err := shared.ConvertDomainEntityToResourceObject(
		instance.Contract,
		resourcesModel.Contract{}.AttributeTypes(),
		ctx,
		convertContractToResourceModel,
	)
	if err != nil {
		return nil, fmt.Errorf("convertInstanceToResourceModel: %w", err)
	}
	plan.Contract = contract

	iso, err := shared.ConvertNullableDomainEntityToResourceObject(
		instance.Iso,
		resourcesModel.Iso{}.AttributeTypes(),
		ctx,
		convertIsoToResourceModel,
	)
	if err != nil {
		return nil, fmt.Errorf("convertInstanceToResourceModel: %w", err)
	}
	plan.Iso = iso

	privateNetwork, err := shared.ConvertNullableDomainEntityToResourceObject(
		instance.PrivateNetwork,
		resourcesModel.PrivateNetwork{}.AttributeTypes(),
		ctx,
		convertPrivateNetworkToResourceModel,
	)
	if err != nil {
		return nil, fmt.Errorf("convertInstanceToResourceModel: %w", err)
	}
	plan.PrivateNetwork = privateNetwork

	resources, err := shared.ConvertDomainEntityToResourceObject(
		instance.Resources,
		resourcesModel.Resources{}.AttributeTypes(),
		ctx,
		convertResourcesToResourceModel,
	)
	if err != nil {
		return nil, fmt.Errorf("convertInstanceToResourceModel: %w", err)
	}
	plan.Resources = resources

	autoScalingGroup, err := shared.ConvertNullableDomainEntityToResourceObject(
		instance.AutoScalingGroup,
		resourcesModel.AutoScalingGroup{}.AttributeTypes(),
		ctx,
		convertAutoScalingGroupToResourceModel,
	)
	if err != nil {
		return nil, fmt.Errorf("convertInstanceToResourceModel: %w", err)
	}
	plan.AutoScalingGroup = autoScalingGroup

	ips, err := shared.ConvertEntitiesToListValue(
		instance.Ips,
		resourcesModel.Ip{}.AttributeTypes(),
		ctx,
		convertIpToResourceModel,
	)
	if err != nil {
		return nil, fmt.Errorf("convertInstanceToResourceModel: %w", err)
	}
	plan.Ips = ips

	return &plan, nil
}

func convertImageToResourceModel(
	ctx context.Context,
	image domain.Image,
) (*resourcesModel.Image, error) {
	plan := &resourcesModel.Image{}

	marketApps, marketAppsDiags := basetypes.NewListValueFrom(
		ctx,
		types.StringType,
		image.MarketApps,
	)
	if marketAppsDiags != nil {
		return nil, returnError("convertImageToResourceModel", marketAppsDiags)
	}

	storageTypes, storageTypesDiags := basetypes.NewListValueFrom(
		ctx,
		types.StringType,
		image.StorageTypes,
	)
	if storageTypesDiags != nil {
		return nil, returnError("convertImageToResourceModel", storageTypesDiags)
	}

	plan.Id = basetypes.NewStringValue(string(image.Id))
	plan.Name = basetypes.NewStringValue(image.Name)
	plan.Version = basetypes.NewStringValue(image.Version)
	plan.Family = basetypes.NewStringValue(image.Family)
	plan.Flavour = basetypes.NewStringValue(image.Flavour)
	plan.Architecture = basetypes.NewStringValue(image.Architecture)
	plan.MarketApps = marketApps
	plan.StorageTypes = storageTypes

	return plan, nil
}

func convertContractToResourceModel(
	ctx context.Context,
	contract domain.Contract,
) (*resourcesModel.Contract, error) {

	return &resourcesModel.Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(contract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(contract.Term)),
		Type:             basetypes.NewStringValue(string(contract.Type)),
		EndsAt:           shared.ConvertNullableTimeToStringValue(contract.EndsAt),
		RenewalsAt:       basetypes.NewStringValue(contract.RenewalsAt.String()),
		CreatedAt:        basetypes.NewStringValue(contract.CreatedAt.String()),
		State:            basetypes.NewStringValue(string(contract.State)),
	}, nil
}

func convertIsoToResourceModel(
	ctx context.Context,
	iso domain.Iso,
) (*resourcesModel.Iso, error) {
	return &resourcesModel.Iso{
		Id:   basetypes.NewStringValue(iso.Id),
		Name: basetypes.NewStringValue(iso.Name),
	}, nil
}

func convertPrivateNetworkToResourceModel(
	ctx context.Context,
	privateNetwork domain.PrivateNetwork,
) (*resourcesModel.PrivateNetwork, error) {
	return &resourcesModel.PrivateNetwork{
		Id:     basetypes.NewStringValue(privateNetwork.Id),
		Status: basetypes.NewStringValue(privateNetwork.Status),
		Subnet: basetypes.NewStringValue(privateNetwork.Subnet),
	}, nil
}

func convertResourcesToResourceModel(
	ctx context.Context,
	domainResources domain.Resources,
) (*resourcesModel.Resources, error) {
	var resources resourcesModel.Resources

	cpu, cpuDiags := shared.ConvertDomainEntityToResourceObject(
		domainResources.Cpu,
		resourcesModel.Cpu{}.AttributeTypes(),
		ctx,
		convertCpuToResourceModel,
	)
	if cpuDiags != nil {
		return &resources, cpuDiags
	}
	resources.Cpu = cpu

	memory, memoryDiags := shared.ConvertDomainEntityToResourceObject(
		domainResources.Memory,
		resourcesModel.Memory{}.AttributeTypes(),
		ctx,
		convertMemoryToResourceModel,
	)
	if memoryDiags != nil {
		return &resources, memoryDiags
	}
	resources.Memory = memory

	publicNetworkSpeed, publicNetworkSpeedDiags := shared.ConvertDomainEntityToResourceObject(
		domainResources.PublicNetworkSpeed,
		resourcesModel.NetworkSpeed{}.AttributeTypes(),
		ctx,
		convertNetworkSpeedToResourceModel,
	)
	if publicNetworkSpeedDiags != nil {
		return &resources, publicNetworkSpeedDiags
	}
	resources.PublicNetworkSpeed = publicNetworkSpeed

	privateNetworkSpeed, privateNetworkSpeedDiags := shared.ConvertDomainEntityToResourceObject(
		domainResources.PrivateNetworkSpeed,
		resourcesModel.NetworkSpeed{}.AttributeTypes(),
		ctx,
		convertNetworkSpeedToResourceModel,
	)
	if privateNetworkSpeedDiags != nil {
		return &resources, privateNetworkSpeedDiags
	}
	resources.PrivateNetworkSpeed = privateNetworkSpeed

	return &resources, nil
}

func convertCpuToResourceModel(
	ctx context.Context,
	cpu domain.Cpu,
) (*resourcesModel.Cpu, error) {
	return &resourcesModel.Cpu{
		Value: basetypes.NewInt64Value(int64(cpu.Value)),
		Unit:  basetypes.NewStringValue(cpu.Unit),
	}, nil
}

func convertMemoryToResourceModel(
	ctx context.Context,
	memory domain.Memory,
) (*resourcesModel.Memory, error) {
	return &resourcesModel.Memory{
		Value: basetypes.NewFloat64Value(memory.Value),
		Unit:  basetypes.NewStringValue(memory.Unit),
	}, nil
}

func convertNetworkSpeedToResourceModel(
	ctx context.Context,
	networkSpeed domain.NetworkSpeed,
) (*resourcesModel.NetworkSpeed, error) {

	return &resourcesModel.NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(networkSpeed.Value)),
		Unit:  basetypes.NewStringValue(networkSpeed.Unit),
	}, nil
}

func convertAutoScalingGroupToResourceModel(
	ctx context.Context,
	autoScalingGroup domain.AutoScalingGroup,
) (*resourcesModel.AutoScalingGroup, error) {

	loadBalancer, loadBalancerDiags := shared.ConvertNullableDomainEntityToResourceObject(
		autoScalingGroup.LoadBalancer,
		resourcesModel.LoadBalancer{}.AttributeTypes(),
		ctx,
		convertLoadBalancerToResourceModel,
	)
	if loadBalancerDiags != nil {
		return nil, loadBalancerDiags
	}
	return &resourcesModel.AutoScalingGroup{
		Id:    basetypes.NewStringValue(autoScalingGroup.Id.String()),
		Type:  basetypes.NewStringValue(string(autoScalingGroup.Type)),
		State: basetypes.NewStringValue(string(autoScalingGroup.State)),
		DesiredAmount: shared.ConvertNullableIntToInt64Value(
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
		StartsAt: shared.ConvertNullableTimeToStringValue(
			autoScalingGroup.StartsAt,
		),
		EndsAt: shared.ConvertNullableTimeToStringValue(
			autoScalingGroup.EndsAt,
		),
		MinimumAmount: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.MinimumAmount,
		),
		MaximumAmount: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.MaximumAmount,
		),
		CpuThreshold: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.CpuThreshold,
		),
		WarmupTime: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.WarmupTime,
		),
		CooldownTime: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.CooldownTime,
		),
		LoadBalancer: loadBalancer,
	}, nil
}

func convertLoadBalancerToResourceModel(
	ctx context.Context,
	loadBalancer domain.LoadBalancer,
) (*resourcesModel.LoadBalancer, error) {

	resources, diags := shared.ConvertDomainEntityToResourceObject(
		loadBalancer.Resources,
		resourcesModel.Resources{}.AttributeTypes(),
		ctx,
		convertResourcesToResourceModel,
	)
	if diags != nil {
		return nil, diags
	}

	contract, diags := shared.ConvertDomainEntityToResourceObject(
		loadBalancer.Contract,
		resourcesModel.Contract{}.AttributeTypes(),
		ctx,
		convertContractToResourceModel,
	)
	if diags != nil {
		return nil, diags
	}

	configuration, diags := shared.ConvertNullableDomainEntityToResourceObject(
		loadBalancer.Configuration,
		resourcesModel.LoadBalancerConfiguration{}.AttributeTypes(),
		ctx,
		convertLoadBalancerConfigurationToResourceModel,
	)
	if diags != nil {
		return nil, diags
	}

	privateNetwork, diags := shared.ConvertNullableDomainEntityToResourceObject(
		loadBalancer.PrivateNetwork,
		resourcesModel.PrivateNetwork{}.AttributeTypes(),
		ctx,
		convertPrivateNetworkToResourceModel,
	)
	if diags != nil {
		return nil, diags
	}

	ips, diags := shared.ConvertEntitiesToListValue(
		loadBalancer.Ips,
		resourcesModel.Ip{}.AttributeTypes(),
		ctx,
		convertIpToResourceModel,
	)
	if diags != nil {
		return nil, diags
	}

	return &resourcesModel.LoadBalancer{
		Id: basetypes.NewStringValue(loadBalancer.Id.String()),
		Type: basetypes.NewStringValue(
			loadBalancer.Type.String(),
		),
		Resources: resources,
		Region:    basetypes.NewStringValue(loadBalancer.Region),
		Reference: shared.ConvertNullableStringToStringValue(
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

func convertLoadBalancerConfigurationToResourceModel(
	ctx context.Context,
	configuration domain.LoadBalancerConfiguration,
) (*resourcesModel.LoadBalancerConfiguration, error) {

	healthCheckObject, diags := shared.ConvertNullableDomainEntityToResourceObject(
		configuration.HealthCheck,
		resourcesModel.HealthCheck{}.AttributeTypes(),
		ctx,
		convertHealthCheckToResourceModel,
	)
	if diags != nil {
		return nil, diags
	}

	stickySessionObject, diags := shared.ConvertNullableDomainEntityToResourceObject(
		configuration.StickySession,
		resourcesModel.StickySession{}.AttributeTypes(),
		ctx,
		convertStickySessionToResourceModel,
	)
	if diags != nil {
		return nil, diags
	}

	return &resourcesModel.LoadBalancerConfiguration{
		Balance:       basetypes.NewStringValue(configuration.Balance.String()),
		HealthCheck:   healthCheckObject,
		StickySession: stickySessionObject,
		XForwardedFor: basetypes.NewBoolValue(configuration.XForwardedFor),
		IdleTimeout:   basetypes.NewInt64Value(int64(configuration.IdleTimeout)),
		TargetPort:    basetypes.NewInt64Value(int64(configuration.TargetPort)),
	}, nil
}

func convertHealthCheckToResourceModel(
	ctx context.Context,
	healthCheck domain.HealthCheck,
) (*resourcesModel.HealthCheck, error) {

	return &resourcesModel.HealthCheck{
		Method: basetypes.NewStringValue(string(healthCheck.Method)),
		Uri:    basetypes.NewStringValue(healthCheck.Uri),
		Host:   shared.ConvertNullableStringToStringValue(healthCheck.Host),
		Port:   basetypes.NewInt64Value(int64(healthCheck.Port)),
	}, nil
}

func convertStickySessionToResourceModel(
	ctx context.Context,
	stickySession domain.StickySession,
) (*resourcesModel.StickySession, error) {

	return &resourcesModel.StickySession{
		Enabled:     basetypes.NewBoolValue(stickySession.Enabled),
		MaxLifeTime: basetypes.NewInt64Value(int64(stickySession.MaxLifeTime)),
	}, nil
}

func convertIpToResourceModel(
	ctx context.Context,
	ip domain.Ip,
) (*resourcesModel.Ip, error) {

	ddos, diags := shared.ConvertNullableDomainEntityToResourceObject(
		ip.Ddos,
		resourcesModel.Ddos{}.AttributeTypes(),
		ctx,
		convertDdosToResourceModel,
	)

	if diags != nil {
		return nil, diags
	}

	return &resourcesModel.Ip{
		Ip:            basetypes.NewStringValue(ip.Ip),
		PrefixLength:  basetypes.NewStringValue(ip.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(ip.Version)),
		NullRouted:    basetypes.NewBoolValue(ip.NullRouted),
		MainIp:        basetypes.NewBoolValue(ip.MainIp),
		NetworkType:   basetypes.NewStringValue(string(ip.NetworkType)),
		ReverseLookup: shared.ConvertNullableStringToStringValue(ip.ReverseLookup),
		Ddos:          ddos,
	}, nil
}

func convertDdosToResourceModel(
	ctx context.Context,
	ddos domain.Ddos,
) (*resourcesModel.Ddos, error) {

	return &resourcesModel.Ddos{
		DetectionProfile: basetypes.NewStringValue(ddos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(ddos.ProtectionType),
	}, nil
}

func convertInstanceResourceModelToCreateInstanceOpts(
	instanceResourceModel resourcesModel.Instance,
	allowedInstancedTypes []string,
	ctx context.Context,
) (*domain.Instance, error) {
	var sshKey *value_object.SshKey
	var rootDiskSize *value_object.RootDiskSize

	image := resourcesModel.Image{}
	imageDiags := instanceResourceModel.Image.As(ctx, &image, basetypes.ObjectAsOptions{})
	if imageDiags != nil {
		return nil, returnError(
			"convertInstanceResourceModelToCreateInstanceOpts",
			imageDiags,
		)
	}

	contract := resourcesModel.Contract{}
	contractDiags := instanceResourceModel.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		return nil, returnError(
			"convertInstanceResourceModelToCreateInstanceOpts",
			imageDiags,
		)
	}

	rootDiskStorageType, err := enum.NewRootDiskStorageType(instanceResourceModel.RootDiskStorageType.ValueString())
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstanceResourceModelToCreateInstanceOpts: %w",
			err,
		)
	}

	imageId, err := enum.NewImageId(image.Id.ValueString())
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstanceResourceModelToCreateInstanceOpts: %w",
			err,
		)
	}

	contractType, err := enum.NewContractType(contract.Type.ValueString())
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstanceResourceModelToCreateInstanceOpts: %w",
			err,
		)
	}

	contractTerm, err := enum.NewContractTerm(int(contract.Term.ValueInt64()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstanceResourceModelToCreateInstanceOpts: %w",
			err,
		)
	}

	billingFrequency, err := enum.NewContractBillingFrequency(int(contract.BillingFrequency.ValueInt64()))
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstanceResourceModelToCreateInstanceOpts: %w",
			err,
		)
	}

	if instanceResourceModel.SshKey.ValueString() != "" {
		sshKey, err = value_object.NewSshKey(instanceResourceModel.SshKey.ValueString())
		if err != nil {
			return nil, fmt.Errorf(
				"convertInstanceResourceModelToCreateInstanceOpts: %w",
				err,
			)
		}
	}

	if instanceResourceModel.RootDiskSize.ValueInt64() != 0 {
		rootDiskSize, err = value_object.NewRootDiskSize(int(instanceResourceModel.RootDiskSize.ValueInt64()))
		if err != nil {
			return nil, fmt.Errorf(
				"convertInstanceResourceModelToCreateInstanceOpts: %w",
				err,
			)
		}
	}

	instanceType, err := value_object.NewInstanceType(
		instanceResourceModel.Type.ValueString(),
		allowedInstancedTypes,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"convertInstanceResourceModelToCreateInstanceOpts: %w",
			err,
		)
	}

	createInstanceOpts := domain.NewCreateInstance(
		instanceResourceModel.Region.ValueString(),
		*instanceType,
		rootDiskStorageType,
		imageId,
		contractType,
		contractTerm,
		billingFrequency,
		domain.OptionalCreateInstanceValues{
			MarketAppId: shared.ConvertStringPointerValueToNullableString(
				instanceResourceModel.MarketAppId,
			),
			Reference: shared.ConvertStringPointerValueToNullableString(
				instanceResourceModel.Reference,
			),
			SshKey:       sshKey,
			RootDiskSize: rootDiskSize,
		},
	)

	return &createInstanceOpts, nil
}

func convertInstancesToDataSourceModel(domainInstances domain.Instances) dataSourceModel.Instances {
	var instances dataSourceModel.Instances

	for _, domainInstance := range domainInstances {
		instance := convertInstanceToDataSourceModel(domainInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func convertInstanceToDataSourceModel(domainInstance domain.Instance) dataSourceModel.Instance {
	instance := dataSourceModel.Instance{
		Id:     basetypes.NewStringValue(domainInstance.Id.String()),
		Region: basetypes.NewStringValue(domainInstance.Region),
		Reference: shared.ConvertNullableStringToStringValue(
			domainInstance.Reference,
		),
		Resources: convertResourcesToDataSourceModel(
			domainInstance.Resources,
		),
		Image:         convertImageToDataSourceModel(domainInstance.Image),
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
		StartedAt: shared.ConvertNullableTimeToStringValue(
			domainInstance.StartedAt,
		),
		Contract: convertContractToDataSourceModel(
			domainInstance.Contract,
		),
		MarketAppId: shared.ConvertNullableStringToStringValue(
			domainInstance.MarketAppId,
		),
		AutoScalingGroup: shared.ConvertNullableDomainEntityToDatasourceModel(
			domainInstance.AutoScalingGroup,
			convertAutoScalingGroupToDataSourceModel,
		),
		Iso: shared.ConvertNullableDomainEntityToDatasourceModel(
			domainInstance.Iso,
			convertIsoToDataSourceModel,
		),
		PrivateNetwork: shared.ConvertNullableDomainEntityToDatasourceModel(
			domainInstance.PrivateNetwork,
			convertPrivateNetworkToDataSourceModel,
		),
	}

	for _, autoScalingGroupIp := range domainInstance.Ips {
		ip := convertIpToDataSourceModel(autoScalingGroupIp)
		instance.Ips = append(instance.Ips, ip)
	}

	return instance
}

func convertResourcesToDataSourceModel(resources domain.Resources) dataSourceModel.Resources {
	return dataSourceModel.Resources{
		Cpu:    convertCpuToDataSourceModel(resources.Cpu),
		Memory: convertMemoryToDataSourceModel(resources.Memory),
		PublicNetworkSpeed: convertNetworkSpeedToDataSourceModel(
			resources.PublicNetworkSpeed,
		),
		PrivateNetworkSpeed: convertNetworkSpeedToDataSourceModel(
			resources.PrivateNetworkSpeed,
		),
	}
}

func convertCpuToDataSourceModel(cpu domain.Cpu) dataSourceModel.Cpu {
	return dataSourceModel.Cpu{
		Value: basetypes.NewInt64Value(int64(cpu.Value)),
		Unit:  basetypes.NewStringValue(cpu.Unit),
	}
}

func convertMemoryToDataSourceModel(memory domain.Memory) dataSourceModel.Memory {
	return dataSourceModel.Memory{
		Value: basetypes.NewFloat64Value(memory.Value),
		Unit:  basetypes.NewStringValue(memory.Unit),
	}
}

func convertNetworkSpeedToDataSourceModel(networkSpeed domain.NetworkSpeed) dataSourceModel.NetworkSpeed {
	return dataSourceModel.NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(networkSpeed.Value)),
		Unit:  basetypes.NewStringValue(networkSpeed.Unit),
	}
}

func convertImageToDataSourceModel(domainImage domain.Image) dataSourceModel.Image {
	image := dataSourceModel.Image{
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

func convertContractToDataSourceModel(contract domain.Contract) dataSourceModel.Contract {
	return dataSourceModel.Contract{
		BillingFrequency: basetypes.NewInt64Value(
			int64(contract.BillingFrequency),
		),
		Term:       basetypes.NewInt64Value(int64(contract.Term)),
		Type:       basetypes.NewStringValue(string(contract.Type)),
		EndsAt:     shared.ConvertNullableTimeToStringValue(contract.EndsAt),
		RenewalsAt: basetypes.NewStringValue(contract.RenewalsAt.String()),
		CreatedAt:  basetypes.NewStringValue(contract.CreatedAt.String()),
		State:      basetypes.NewStringValue(string(contract.State)),
	}
}

func convertAutoScalingGroupToDataSourceModel(autoScalingGroup domain.AutoScalingGroup) *dataSourceModel.AutoScalingGroup {
	return &dataSourceModel.AutoScalingGroup{
		Id:    basetypes.NewStringValue(autoScalingGroup.Id.String()),
		Type:  basetypes.NewStringValue(string(autoScalingGroup.Type)),
		State: basetypes.NewStringValue(string(autoScalingGroup.State)),
		DesiredAmount: shared.ConvertNullableIntToInt64Value(
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
		StartsAt: shared.ConvertNullableTimeToStringValue(
			autoScalingGroup.StartsAt,
		),
		EndsAt: shared.ConvertNullableTimeToStringValue(
			autoScalingGroup.EndsAt,
		),
		MinimumAmount: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.MinimumAmount,
		),
		MaximumAmount: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.MaximumAmount,
		),
		CpuThreshold: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.CpuThreshold,
		),
		WarmupTime: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.WarmupTime,
		),
		CooldownTime: shared.ConvertNullableIntToInt64Value(
			autoScalingGroup.CooldownTime,
		),
		LoadBalancer: shared.ConvertNullableDomainEntityToDatasourceModel(
			autoScalingGroup.LoadBalancer,
			convertLoadBalancerToDataSourceModel,
		),
	}
}

func convertLoadBalancerToDataSourceModel(loadBalancer domain.LoadBalancer) *dataSourceModel.LoadBalancer {
	var ips []dataSourceModel.Ip
	for _, ip := range loadBalancer.Ips {
		ips = append(ips, convertIpToDataSourceModel(ip))
	}

	return &dataSourceModel.LoadBalancer{
		Id:        basetypes.NewStringValue(loadBalancer.Id.String()),
		Type:      basetypes.NewStringValue(loadBalancer.Type.String()),
		Resources: convertResourcesToDataSourceModel(loadBalancer.Resources),
		Region:    basetypes.NewStringValue(loadBalancer.Region),
		Reference: shared.ConvertNullableStringToStringValue(loadBalancer.Reference),
		State:     basetypes.NewStringValue(string(loadBalancer.State)),
		Contract:  convertContractToDataSourceModel(loadBalancer.Contract),
		StartedAt: shared.ConvertNullableTimeToStringValue(loadBalancer.StartedAt),
		Ips:       ips,
		LoadBalancerConfiguration: shared.ConvertNullableDomainEntityToDatasourceModel(
			loadBalancer.Configuration,
			convertLoadBalancerConfigurationToDataSourceModel,
		),
		PrivateNetwork: shared.ConvertNullableDomainEntityToDatasourceModel(
			loadBalancer.PrivateNetwork,
			convertPrivateNetworkToDataSourceModel,
		),
	}
}

func convertLoadBalancerConfigurationToDataSourceModel(configuration domain.LoadBalancerConfiguration) *dataSourceModel.LoadBalancerConfiguration {
	return &dataSourceModel.LoadBalancerConfiguration{
		Balance: basetypes.NewStringValue(configuration.Balance.String()),
		HealthCheck: shared.ConvertNullableDomainEntityToDatasourceModel(
			configuration.HealthCheck,
			convertHealthCheckToDataSourceModel,
		),
		StickySession: shared.ConvertNullableDomainEntityToDatasourceModel(
			configuration.StickySession,
			convertStickySessionToDataSourceModel,
		),
		XForwardedFor: basetypes.NewBoolValue(configuration.XForwardedFor),
		IdleTimeout:   basetypes.NewInt64Value(int64(configuration.IdleTimeout)),
		TargetPort:    basetypes.NewInt64Value(int64(configuration.TargetPort)),
	}
}

func convertHealthCheckToDataSourceModel(healthCheck domain.HealthCheck) *dataSourceModel.HealthCheck {
	return &dataSourceModel.HealthCheck{
		Method: basetypes.NewStringValue(healthCheck.Method.String()),
		Uri:    basetypes.NewStringValue(healthCheck.Uri),
		Host:   shared.ConvertNullableStringToStringValue(healthCheck.Host),
		Port:   basetypes.NewInt64Value(int64(healthCheck.Port)),
	}
}

func convertStickySessionToDataSourceModel(stickySession domain.StickySession) *dataSourceModel.StickySession {
	return &dataSourceModel.StickySession{
		Enabled:     basetypes.NewBoolValue(stickySession.Enabled),
		MaxLifeTime: basetypes.NewInt64Value(int64(stickySession.MaxLifeTime)),
	}
}

func convertPrivateNetworkToDataSourceModel(privateNetwork domain.PrivateNetwork) *dataSourceModel.PrivateNetwork {
	return &dataSourceModel.PrivateNetwork{
		Id:     basetypes.NewStringValue(privateNetwork.Id),
		Status: basetypes.NewStringValue(privateNetwork.Status),
		Subnet: basetypes.NewStringValue(privateNetwork.Subnet),
	}
}

func convertIsoToDataSourceModel(iso domain.Iso) *dataSourceModel.Iso {
	return &dataSourceModel.Iso{
		Id:   basetypes.NewStringValue(iso.Id),
		Name: basetypes.NewStringValue(iso.Name),
	}
}

func convertIpToDataSourceModel(ip domain.Ip) dataSourceModel.Ip {
	return dataSourceModel.Ip{
		Ip:            basetypes.NewStringValue(ip.Ip),
		PrefixLength:  basetypes.NewStringValue(ip.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(ip.Version)),
		NullRouted:    basetypes.NewBoolValue(ip.NullRouted),
		MainIp:        basetypes.NewBoolValue(ip.MainIp),
		NetworkType:   basetypes.NewStringValue(string(ip.NetworkType)),
		ReverseLookup: shared.ConvertNullableStringToStringValue(ip.ReverseLookup),
		Ddos: shared.ConvertNullableDomainEntityToDatasourceModel(
			ip.Ddos,
			convertDdosToDataSourceModel,
		),
	}
}

func convertDdosToDataSourceModel(ddos domain.Ddos) *dataSourceModel.Ddos {
	return &dataSourceModel.Ddos{
		DetectionProfile: basetypes.NewStringValue(ddos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(ddos.ProtectionType),
	}
}

func convertInstanceResourceModelToUpdateInstanceOpts(
	instanceResourceModel resourcesModel.Instance,
	allowedInstanceTypes []string,
	ctx context.Context,
) (*domain.Instance, error) {

	id, err := value_object.NewUuid(instanceResourceModel.Id.ValueString())
	if err != nil {
		return nil, fmt.Errorf("convertInstanceResourceModelToUpdateInstanceOpts: %w", err)
	}

	optionalValues := domain.OptionalUpdateInstanceValues{
		Reference: shared.ConvertStringPointerValueToNullableString(instanceResourceModel.Reference),
	}

	if instanceResourceModel.RootDiskSize.ValueInt64() != 0 {
		rootDiskSize, err := value_object.NewRootDiskSize(int(instanceResourceModel.RootDiskSize.ValueInt64()))
		if err != nil {
			return nil, fmt.Errorf(
				"convertInstanceResourceModelToUpdateInstanceOpts: %w",
				err,
			)
		}
		optionalValues.RootDiskSize = rootDiskSize
	}

	contract := resourcesModel.Contract{}
	diags := instanceResourceModel.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, returnError(
			"convertInstanceResourceModelToUpdateInstanceOpts",
			diags,
		)
	}

	if contract.Type.ValueString() != "" {
		contractType, err := enum.NewContractType(contract.Type.ValueString())
		if err != nil {
			return nil, fmt.Errorf(
				"convertInstanceResourceModelToUpdateInstanceOpts: %w",
				err,
			)
		}
		optionalValues.ContractType = &contractType
	}

	if contract.Term.ValueInt64() != 0 {
		contractTerm, err := enum.NewContractTerm(int(contract.Term.ValueInt64()))
		if err != nil {
			return nil, fmt.Errorf(
				"convertInstanceResourceModelToUpdateInstanceOpts: %w",
				err,
			)
		}
		optionalValues.Term = &contractTerm
	}

	if contract.BillingFrequency.ValueInt64() != 0 {
		billingFrequency, err := enum.NewContractBillingFrequency(
			int(contract.BillingFrequency.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"convertInstanceResourceModelToUpdateInstanceOpts: %w",
				err,
			)
		}
		optionalValues.BillingFrequency = &billingFrequency
	}

	if instanceResourceModel.Type.ValueString() != "" {
		instanceType, err := value_object.NewInstanceType(
			instanceResourceModel.Type.ValueString(),
			allowedInstanceTypes,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"convertInstanceResourceModelToUpdateInstanceOpts: %w",
				err,
			)
		}
		optionalValues.Type = instanceType
	}

	instance := domain.NewUpdateInstance(*id, optionalValues)

	return &instance, nil
}

func convertIntArrayToInt64(items []int) []int64 {
	var convertedItems []int64

	for _, item := range items {
		convertedItems = append(
			convertedItems,
			int64(item),
		)
	}

	return convertedItems
}
