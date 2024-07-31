package to_domain_entity

import (
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/repositories/shared"
)

// AdaptInstance adapts an instance domain entity to an sdk instance model.
func AdaptInstance(
	sdkInstance publicCloud.Instance,
) (
	*domain.Instance,
	error,
) {
	var autoScalingGroup *domain.AutoScalingGroup

	instanceId, err := value_object.NewUuid(sdkInstance.GetId())
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}

	state, err := enum.NewState(string(sdkInstance.GetState()))
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}

	rootDiskSize, err := value_object.NewRootDiskSize(
		int(sdkInstance.GetRootDiskSize()),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}

	rootDiskStorageType, err := enum.NewRootDiskStorageType(
		string(sdkInstance.GetRootDiskStorageType()),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance: %w", err)
	}

	ips, err := adaptIps(sdkInstance.GetIps())
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance:  %w", err)
	}

	contract, err := adaptContract(sdkInstance.GetContract())
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance:  %w", err)
	}

	sdkAutoScalingGroup, _ := sdkInstance.GetAutoScalingGroupOk()
	if sdkAutoScalingGroup != nil {
		autoScalingGroup, err = adaptAutoScalingGroup(*sdkAutoScalingGroup)
		if err != nil {
			return nil, fmt.Errorf("AdaptInstance:  %w", err)
		}
	}

	optionalValues := domain.OptionalInstanceValues{
		Reference:        shared.AdaptNullableStringToValue(sdkInstance.Reference),
		MarketAppId:      shared.AdaptNullableStringToValue(sdkInstance.MarketAppId),
		StartedAt:        shared.AdaptNullableTimeToValue(sdkInstance.StartedAt),
		AutoScalingGroup: autoScalingGroup,
	}

	instance := domain.NewInstance(
		*instanceId,
		sdkInstance.GetRegion(),
		adaptResources(sdkInstance.GetResources()),
		adaptImage(sdkInstance.GetImage()),
		state,
		sdkInstance.GetProductType(),
		sdkInstance.GetHasPublicIpV4(),
		sdkInstance.GetIncludesPrivateNetwork(),
		*rootDiskSize,
		value_object.NewUnvalidatedInstanceType(string(sdkInstance.GetType())),
		rootDiskStorageType,
		ips,
		*contract,
		optionalValues,
	)

	return &instance, nil
}

func AdaptInstanceDetails(
	sdkInstanceDetails publicCloud.InstanceDetails,
) (*domain.Instance, error) {
	var autoScalingGroup *domain.AutoScalingGroup

	instanceId, err := value_object.NewUuid(sdkInstanceDetails.GetId())
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}

	state, err := enum.NewState(string(sdkInstanceDetails.GetState()))
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}

	rootDiskSize, err := value_object.NewRootDiskSize(
		int(sdkInstanceDetails.GetRootDiskSize()),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}

	rootDiskStorageType, err := enum.NewRootDiskStorageType(
		string(sdkInstanceDetails.GetRootDiskStorageType()),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}

	ips, err := adaptIpsDetails(sdkInstanceDetails.GetIps())
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails:  %w", err)
	}

	contract, err := adaptContract(sdkInstanceDetails.GetContract())
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails:  %w", err)
	}

	sdkAutoScalingGroup, _ := sdkInstanceDetails.GetAutoScalingGroupOk()
	if sdkAutoScalingGroup != nil {
		autoScalingGroup, err = adaptAutoScalingGroup(*sdkAutoScalingGroup)
		if err != nil {
			return nil, fmt.Errorf("AdaptInstanceDetails:  %w", err)
		}
	}

	optionalValues := domain.OptionalInstanceValues{
		Reference: shared.AdaptNullableStringToValue(
			sdkInstanceDetails.Reference,
		),
		MarketAppId: shared.AdaptNullableStringToValue(
			sdkInstanceDetails.MarketAppId,
		),
		StartedAt:        shared.AdaptNullableTimeToValue(sdkInstanceDetails.StartedAt),
		AutoScalingGroup: autoScalingGroup,
	}
	if sdkInstanceDetails.Iso.Get() != nil {
		iso := adaptIso(*sdkInstanceDetails.Iso.Get())
		optionalValues.Iso = &iso
	}
	if sdkInstanceDetails.PrivateNetwork.Get() != nil {
		privateNetwork := adaptPrivateNetwork(
			*sdkInstanceDetails.PrivateNetwork.Get(),
		)
		optionalValues.PrivateNetwork = &privateNetwork
	}

	instance := domain.NewInstance(
		*instanceId,
		sdkInstanceDetails.GetRegion(),
		adaptResources(sdkInstanceDetails.GetResources()),
		adaptImageDetails(sdkInstanceDetails.GetImage()),
		state,
		sdkInstanceDetails.GetProductType(),
		sdkInstanceDetails.GetHasPublicIpV4(),
		sdkInstanceDetails.GetIncludesPrivateNetwork(),
		*rootDiskSize,
		value_object.NewUnvalidatedInstanceType(
			string(sdkInstanceDetails.GetType()),
		),
		rootDiskStorageType,
		ips,
		*contract,
		optionalValues,
	)

	return &instance, nil
}

func adaptResources(sdkResources publicCloud.Resources) domain.Resources {
	resources := domain.NewResources(
		adaptCpu(sdkResources.GetCpu()),
		adaptMemory(sdkResources.GetMemory()),
		adaptNetworkSpeed(sdkResources.GetPublicNetworkSpeed()),
		adaptNetworkSpeed(sdkResources.GetPrivateNetworkSpeed()),
	)

	return resources
}

func adaptCpu(sdkCpu publicCloud.Cpu) domain.Cpu {
	return domain.NewCpu(int(sdkCpu.GetValue()), sdkCpu.GetUnit())
}

func adaptMemory(sdkMemory publicCloud.Memory) domain.Memory {
	return domain.NewMemory(float64(sdkMemory.GetValue()), sdkMemory.GetUnit())
}

func adaptNetworkSpeed(sdkNetworkSpeed publicCloud.NetworkSpeed) domain.NetworkSpeed {
	return domain.NewNetworkSpeed(
		int(sdkNetworkSpeed.GetValue()),
		sdkNetworkSpeed.GetUnit(),
	)
}

func adaptImageDetails(sdkImage publicCloud.ImageDetails) domain.Image {
	return domain.NewImage(
		sdkImage.GetId(),
		sdkImage.GetName(),
		sdkImage.GetVersion(),
		sdkImage.GetFamily(),
		sdkImage.GetFlavour(),
		sdkImage.GetMarketApps(),
		sdkImage.GetStorageTypes(),
	)
}

func adaptImage(sdkImage publicCloud.Image) domain.Image {
	return domain.NewImage(
		sdkImage.GetId(),
		sdkImage.GetName(),
		sdkImage.GetVersion(),
		sdkImage.GetFamily(),
		sdkImage.GetFlavour(),
		[]string{},
		[]string{},
	)
}

func adaptIpsDetails(sdkIps []publicCloud.IpDetails) (domain.Ips, error) {
	var ips domain.Ips
	for _, sdkIp := range sdkIps {
		ip, err := adaptIpDetails(sdkIp)
		if err != nil {
			return nil, fmt.Errorf("adaptIpsDetails: %w", err)
		}
		ips = append(ips, *ip)
	}

	return ips, nil
}

func adaptIps(sdkIps []publicCloud.Ip) (domain.Ips, error) {
	var ips domain.Ips
	for _, sdkIp := range sdkIps {
		ip, err := adaptIp(sdkIp)
		if err != nil {
			return nil, fmt.Errorf("adaptIps: %w", err)
		}
		ips = append(ips, *ip)
	}

	return ips, nil
}

func adaptIpDetails(sdkIp publicCloud.IpDetails) (*domain.Ip, error) {
	networkType, err := enum.NewNetworkType(string(sdkIp.GetNetworkType()))
	if err != nil {
		return nil, fmt.Errorf("adaptIpDetails: %w", err)
	}

	optionalIpValues := domain.OptionalIpValues{
		ReverseLookup: shared.AdaptNullableStringToValue(sdkIp.ReverseLookup),
	}

	sdkDdos, _ := sdkIp.GetDdosOk()
	if sdkDdos != nil {
		ddos := adaptDdos(*sdkDdos)
		optionalIpValues.Ddos = &ddos
	}

	ip := domain.NewIp(
		sdkIp.GetIp(),
		sdkIp.GetPrefixLength(),
		int(sdkIp.GetVersion()),
		sdkIp.GetNullRouted(),
		sdkIp.GetMainIp(),
		networkType,
		optionalIpValues,
	)

	return &ip, nil
}

func adaptIp(sdkIp publicCloud.Ip) (*domain.Ip, error) {
	networkType, err := enum.NewNetworkType(string(sdkIp.GetNetworkType()))
	if err != nil {
		return nil, fmt.Errorf(
			"adaptIpDetails: %w",
			err,
		)
	}

	optionalIpValues := domain.OptionalIpValues{
		ReverseLookup: shared.AdaptNullableStringToValue(sdkIp.ReverseLookup),
	}

	ip := domain.NewIp(
		sdkIp.GetIp(),
		sdkIp.GetPrefixLength(),
		int(sdkIp.GetVersion()),
		sdkIp.GetNullRouted(),
		sdkIp.GetMainIp(),
		networkType,
		optionalIpValues,
	)

	return &ip, nil
}

func adaptDdos(sdkDdos publicCloud.Ddos) domain.Ddos {
	return domain.NewDdos(
		sdkDdos.GetDetectionProfile(),
		sdkDdos.GetProtectionType(),
	)
}

func adaptContract(sdkContract publicCloud.Contract) (*domain.Contract, error) {
	billingFrequency, err := enum.NewContractBillingFrequency(
		int(sdkContract.GetBillingFrequency()),
	)
	if err != nil {
		return nil, fmt.Errorf("adaptContract: %w", err)
	}

	contractTerm, err := enum.NewContractTerm(int(sdkContract.GetTerm()))
	if err != nil {
		return nil, fmt.Errorf("adaptContract: %w", err)
	}

	contractType, err := enum.NewContractType(string(sdkContract.GetType()))
	if err != nil {
		return nil, fmt.Errorf("adaptContract: %w", err)
	}

	contractState, err := enum.NewContractState(string(sdkContract.GetState()))
	if err != nil {
		return nil, fmt.Errorf("adaptContract: %w", err)
	}

	contract, err := domain.NewContract(
		billingFrequency,
		contractTerm,
		contractType,
		sdkContract.GetRenewalsAt(),
		sdkContract.GetCreatedAt(),
		contractState,
		shared.AdaptNullableTimeToValue(sdkContract.EndsAt),
	)

	if err != nil {
		return nil, fmt.Errorf("adaptContract: %w", err)
	}

	return contract, nil
}

func adaptIso(sdkIso publicCloud.Iso) domain.Iso {
	return domain.NewIso(sdkIso.GetId(), sdkIso.GetName())
}

func adaptPrivateNetwork(sdkPrivateNetwork publicCloud.PrivateNetwork) domain.PrivateNetwork {
	return domain.PrivateNetwork{
		Id:     sdkPrivateNetwork.GetPrivateNetworkId(),
		Status: sdkPrivateNetwork.GetStatus(),
		Subnet: sdkPrivateNetwork.GetSubnet(),
	}
}

func AdaptAutoScalingGroupDetails(
	sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
) (
	*domain.AutoScalingGroup,
	error,
) {
	var loadBalancer *domain.LoadBalancer

	autoScalingGroupId, err := value_object.NewUuid(sdkAutoScalingGroup.GetId())
	if err != nil {
		return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
	}

	autoScalingGroupType, err := enum.NewAutoScalingGroupType(
		string(sdkAutoScalingGroup.GetType()),
	)
	if err != nil {
		return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
	}

	state, err := enum.NewAutoScalingGroupState(
		string(sdkAutoScalingGroup.GetState()),
	)
	if err != nil {
		return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
	}

	reference, err := value_object.NewAutoScalingGroupReference(
		sdkAutoScalingGroup.GetReference(),
	)
	if err != nil {
		return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
	}

	sdkLoadBalancer, _ := sdkAutoScalingGroup.GetLoadBalancerOk()
	if sdkLoadBalancer != nil {
		loadBalancer, err = adaptLoadBalancer(*sdkLoadBalancer)
		if err != nil {
			return nil, fmt.Errorf("adaptAutoScalingGroupDetails: %w", err)
		}
	}

	options := domain.AutoScalingGroupOptions{
		DesiredAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.DesiredAmount),
		MinimumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MinimumAmount),
		MaximumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MaximumAmount),
		CpuThreshold:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CpuThreshold),
		CoolDownTime:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CooldownTime),
		StartsAt:      shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.StartsAt),
		EndsAt:        shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.EndsAt),
		WarmupTime:    shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.WarmupTime),
		LoadBalancer:  loadBalancer,
	}

	autoScalingGroup := domain.NewAutoScalingGroup(
		*autoScalingGroupId,
		autoScalingGroupType,
		state,
		sdkAutoScalingGroup.GetRegion(),
		*reference,
		sdkAutoScalingGroup.GetCreatedAt(),
		sdkAutoScalingGroup.GetUpdatedAt(),
		options,
	)

	return &autoScalingGroup, nil
}

func AdaptLoadBalancerDetails(
	sdkLoadBalancer publicCloud.LoadBalancerDetails,
) (
	*domain.LoadBalancer,
	error,
) {
	loadBalancerId, err := value_object.NewUuid(sdkLoadBalancer.Id)
	if err != nil {
		return nil, fmt.Errorf("AdaptLoadBalancerDetails: %w", err)
	}

	state, err := enum.NewState(string(sdkLoadBalancer.GetState()))
	if err != nil {
		return nil, fmt.Errorf("AdaptLoadBalancerDetails: %w", err)
	}

	contract, err := adaptContract(sdkLoadBalancer.GetContract())
	if err != nil {
		return nil, fmt.Errorf("AdaptLoadBalancerDetails: %w", err)
	}

	ips, err := adaptIpsDetails(sdkLoadBalancer.GetIps())
	if err != nil {
		return nil, fmt.Errorf("AdaptLoadBalancerDetails:  %w", err)
	}

	options := domain.OptionalLoadBalancerValues{
		Reference: shared.AdaptNullableStringToValue(sdkLoadBalancer.Reference),
		StartedAt: shared.AdaptNullableTimeToValue(sdkLoadBalancer.StartedAt),
	}

	if sdkLoadBalancer.Configuration.Get() != nil {
		configuration, err := adaptLoadBalancerConfiguration(sdkLoadBalancer.GetConfiguration())
		if err != nil {
			return nil, fmt.Errorf("AdaptLoadBalancerDetails:  %w", err)
		}
		options.Configuration = configuration
	}

	if sdkLoadBalancer.PrivateNetwork.Get() != nil {
		privateNetwork := adaptPrivateNetwork(*sdkLoadBalancer.PrivateNetwork.Get())
		options.PrivateNetwork = &privateNetwork
	}

	loadBalancer := domain.NewLoadBalancer(
		*loadBalancerId,
		value_object.NewUnvalidatedInstanceType(string(sdkLoadBalancer.GetType())),
		adaptResources(sdkLoadBalancer.GetResources()),
		sdkLoadBalancer.GetRegion(),
		state,
		*contract,
		ips,
		options,
	)

	return &loadBalancer, nil
}

func adaptLoadBalancer(sdkLoadBalancer publicCloud.LoadBalancer) (
	*domain.LoadBalancer,
	error,
) {
	loadBalancerId, err := value_object.NewUuid(sdkLoadBalancer.Id)
	if err != nil {
		return nil, fmt.Errorf("adaptLoadBalancer: %w", err)
	}

	state, err := enum.NewState(string(sdkLoadBalancer.GetState()))
	if err != nil {
		return nil, fmt.Errorf("adaptLoadBalancer: %w", err)
	}

	options := domain.OptionalLoadBalancerValues{
		Reference: shared.AdaptNullableStringToValue(sdkLoadBalancer.Reference),
		StartedAt: shared.AdaptNullableTimeToValue(sdkLoadBalancer.StartedAt),
	}

	loadBalancer := domain.NewLoadBalancer(
		*loadBalancerId,
		value_object.NewUnvalidatedInstanceType(string(sdkLoadBalancer.GetType())),
		adaptResources(sdkLoadBalancer.GetResources()),
		"",
		state,
		domain.Contract{},
		domain.Ips{},
		options,
	)

	return &loadBalancer, nil
}

func adaptLoadBalancerConfiguration(sdkLoadBalancerConfiguration publicCloud.LoadBalancerConfiguration) (
	*domain.LoadBalancerConfiguration,
	error,
) {
	balance, err := enum.NewBalance(string(sdkLoadBalancerConfiguration.GetBalance()))
	if err != nil {
		return nil, fmt.Errorf("adaptLoadBalancerConfiguration: %w", err)
	}

	options := domain.OptionalLoadBalancerConfigurationOptions{
		HealthCheck: nil,
	}
	if sdkLoadBalancerConfiguration.StickySession.Get() != nil {
		stickySession := adaptStickySession(*sdkLoadBalancerConfiguration.StickySession.Get())
		options.StickySession = &stickySession
	}
	if sdkLoadBalancerConfiguration.HealthCheck.Get() != nil {
		healthCheck, err := adaptHealthCheck(*sdkLoadBalancerConfiguration.HealthCheck.Get())
		if err != nil {
			return nil, fmt.Errorf("adaptLoadBalancerConfiguration: %w", err)
		}

		options.HealthCheck = healthCheck
	}

	configuration := domain.NewLoadBalancerConfiguration(
		balance,
		sdkLoadBalancerConfiguration.GetXForwardedFor(),
		int(sdkLoadBalancerConfiguration.GetIdleTimeOut()),
		int(sdkLoadBalancerConfiguration.GetTargetPort()),
		options,
	)

	return &configuration, nil
}

func adaptStickySession(sdkStickySession publicCloud.StickySession) domain.StickySession {
	return domain.NewStickySession(
		sdkStickySession.GetEnabled(),
		int(sdkStickySession.GetMaxLifeTime()),
	)
}

func adaptHealthCheck(sdkHealthCheck publicCloud.HealthCheck) (
	*domain.HealthCheck,
	error,
) {
	method, err := enum.NewMethod(sdkHealthCheck.GetMethod())
	if err != nil {
		return nil, fmt.Errorf("adaptHealthCheck: %w", err)
	}

	healthCheck := domain.NewHealthCheck(
		method,
		sdkHealthCheck.GetUri(),
		int(sdkHealthCheck.GetPort()),
		domain.OptionalHealthCheckValues{
			Host: shared.AdaptNullableStringToValue(sdkHealthCheck.Host),
		},
	)

	return &healthCheck, nil
}

func AdaptInstanceType(sdkInstanceType publicCloud.InstanceType) (
	*domain.InstanceType,
	error,
) {
	resources := adaptResources(sdkInstanceType.GetResources())
	prices := adaptPrices(sdkInstanceType.GetPrices())

	optional := domain.OptionalInstanceTypeValues{}

	sdkStorageTypes, _ := sdkInstanceType.GetStorageTypesOk()
	if sdkStorageTypes != nil {
		storageTypes, err := adaptStorageTypes(sdkStorageTypes)
		if err != nil {
			return nil, fmt.Errorf("AdaptInstanceType: %w", err)
		}
		optional.StorageTypes = storageTypes
	}

	instanceType := domain.NewInstanceType(
		sdkInstanceType.GetName(),
		resources,
		prices,
		optional,
	)

	return &instanceType, nil
}

func adaptPrices(sdkPrices publicCloud.Prices) domain.Prices {
	return domain.NewPrices(
		sdkPrices.GetCurrency(),
		sdkPrices.GetCurrencySymbol(),
		adaptPrice(sdkPrices.GetCompute()),
		adaptStorage(sdkPrices.GetStorage()),
	)
}

func adaptStorage(sdkStorage publicCloud.Storage) domain.Storage {
	return domain.NewStorage(
		adaptPrice(sdkStorage.Local),
		adaptPrice(sdkStorage.Central),
	)
}

func adaptPrice(sdkPrice publicCloud.Price) domain.Price {
	return domain.NewPrice(sdkPrice.GetHourlyPrice(), sdkPrice.GetMonthlyPrice())
}

func adaptStorageTypes(sdkStorageTypes []publicCloud.RootDiskStorageType) (
	*domain.StorageTypes,
	error,
) {
	var storageTypes domain.StorageTypes

	for _, sdkStorageType := range sdkStorageTypes {
		storageType, err := enum.NewRootDiskStorageType(string(sdkStorageType))
		if err != nil {
			return nil, fmt.Errorf("adaptStorageTypes: %w", err)
		}
		storageTypes = append(storageTypes, storageType)
	}

	return &storageTypes, nil
}

func AdaptRegion(sdkRegion publicCloud.Region) domain.Region {
	return domain.NewRegion(sdkRegion.GetName(), sdkRegion.GetLocation())
}

func adaptAutoScalingGroup(
	sdkAutoScalingGroup publicCloud.AutoScalingGroup,
) (
	*domain.AutoScalingGroup,
	error,
) {
	autoScalingGroupId, err := value_object.NewUuid(sdkAutoScalingGroup.GetId())
	if err != nil {
		return nil, fmt.Errorf("adaptAutoScalingGroup: %w", err)
	}

	autoScalingGroupType, err := enum.NewAutoScalingGroupType(
		string(sdkAutoScalingGroup.GetType()),
	)
	if err != nil {
		return nil, fmt.Errorf("adaptAutoScalingGroup: %w", err)
	}

	state, err := enum.NewAutoScalingGroupState(string(sdkAutoScalingGroup.GetState()))
	if err != nil {
		return nil, fmt.Errorf("adaptAutoScalingGroup: %w", err)
	}

	reference, err := value_object.NewAutoScalingGroupReference(sdkAutoScalingGroup.GetReference())

	if err != nil {
		return nil, fmt.Errorf("adaptAutoScalingGroup: %w", err)
	}

	options := domain.AutoScalingGroupOptions{
		DesiredAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.DesiredAmount),
		MinimumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MinimumAmount),
		MaximumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MaximumAmount),
		CpuThreshold:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CpuThreshold),
		CoolDownTime:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CooldownTime),
		StartsAt:      shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.StartsAt),
		EndsAt:        shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.EndsAt),
		WarmupTime:    shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.WarmupTime),
	}

	autoScalingGroup := domain.NewAutoScalingGroup(
		*autoScalingGroupId,
		autoScalingGroupType,
		state,
		sdkAutoScalingGroup.GetRegion(),
		*reference,
		sdkAutoScalingGroup.GetCreatedAt(),
		sdkAutoScalingGroup.GetUpdatedAt(),
		options,
	)

	return &autoScalingGroup, nil
}
