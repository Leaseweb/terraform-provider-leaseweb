package to_domain_entity

import (
	"fmt"

	sdkModel "github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	domainEntity "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// AdaptInstance adapts publicCloud.Instance to public_cloud.Instance.
func AdaptInstance(
	sdkInstance sdkModel.Instance,
) (
	*domainEntity.Instance,
	error,
) {
	var autoScalingGroup *domainEntity.AutoScalingGroup

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

	optionalValues := domainEntity.OptionalInstanceValues{
		Reference:        shared.AdaptNullableStringToValue(sdkInstance.Reference),
		MarketAppId:      shared.AdaptNullableStringToValue(sdkInstance.MarketAppId),
		StartedAt:        shared.AdaptNullableTimeToValue(sdkInstance.StartedAt),
		AutoScalingGroup: autoScalingGroup,
	}

	instance := domainEntity.NewInstance(
		sdkInstance.GetId(),
		domainEntity.Region{Name: sdkInstance.GetRegion()},
		adaptResources(sdkInstance.GetResources()),
		adaptImage(sdkInstance.GetImage()),
		state,
		sdkInstance.GetProductType(),
		sdkInstance.GetHasPublicIpV4(),
		sdkInstance.GetIncludesPrivateNetwork(),
		*rootDiskSize,
		domainEntity.InstanceType{Name: string(sdkInstance.GetType())},
		rootDiskStorageType,
		ips,
		*contract,
		optionalValues,
	)

	return &instance, nil
}

// AdaptInstanceDetails adapts publicCloud.InstanceDetails to public_cloud.Instance.
func AdaptInstanceDetails(sdkInstanceDetails sdkModel.InstanceDetails) (
	*domainEntity.Instance,
	error,
) {
	var autoScalingGroup *domainEntity.AutoScalingGroup

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

	optionalValues := domainEntity.OptionalInstanceValues{
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
	if sdkInstanceDetails.Volume.Get() != nil {
		volume := adaptVolume(
			*sdkInstanceDetails.Volume.Get(),
		)
		optionalValues.Volume = &volume
	}

	instance := domainEntity.NewInstance(
		sdkInstanceDetails.GetId(),
		domainEntity.Region{Name: sdkInstanceDetails.GetRegion()},
		adaptResources(sdkInstanceDetails.GetResources()),
		adaptImage(sdkInstanceDetails.GetImage()),
		state,
		sdkInstanceDetails.GetProductType(),
		sdkInstanceDetails.GetHasPublicIpV4(),
		sdkInstanceDetails.GetIncludesPrivateNetwork(),
		*rootDiskSize,
		domainEntity.InstanceType{Name: string(sdkInstanceDetails.GetType())},
		rootDiskStorageType,
		ips,
		*contract,
		optionalValues,
	)

	return &instance, nil
}

func adaptResources(sdkResources sdkModel.Resources) domainEntity.Resources {
	resources := domainEntity.NewResources(
		adaptCpu(sdkResources.GetCpu()),
		adaptMemory(sdkResources.GetMemory()),
		adaptNetworkSpeed(sdkResources.GetPublicNetworkSpeed()),
		adaptNetworkSpeed(sdkResources.GetPrivateNetworkSpeed()),
	)

	return resources
}

func adaptCpu(sdkCpu sdkModel.Cpu) domainEntity.Cpu {
	return domainEntity.NewCpu(int(sdkCpu.GetValue()), sdkCpu.GetUnit())
}

func adaptMemory(sdkMemory sdkModel.Memory) domainEntity.Memory {
	return domainEntity.NewMemory(float64(sdkMemory.GetValue()), sdkMemory.GetUnit())
}

func adaptNetworkSpeed(sdkNetworkSpeed sdkModel.NetworkSpeed) domainEntity.NetworkSpeed {
	return domainEntity.NewNetworkSpeed(
		int(sdkNetworkSpeed.GetValue()),
		sdkNetworkSpeed.GetUnit(),
	)
}

func adaptImage(sdkImage sdkModel.Image) domainEntity.Image {
	return domainEntity.NewImage(
		sdkImage.GetId(),
		sdkImage.GetName(),
		nil,
		sdkImage.GetFamily(),
		sdkImage.GetFlavour(),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		sdkImage.GetCustom(),
		nil,
		[]string{},
		[]string{},
	)
}

func adaptIpsDetails(sdkIps []sdkModel.IpDetails) (domainEntity.Ips, error) {
	var ips domainEntity.Ips
	for _, sdkIp := range sdkIps {
		ip, err := adaptIpDetails(sdkIp)
		if err != nil {
			return nil, fmt.Errorf("adaptIpsDetails: %w", err)
		}
		ips = append(ips, *ip)
	}

	return ips, nil
}

func adaptIps(sdkIps []sdkModel.Ip) (domainEntity.Ips, error) {
	var ips domainEntity.Ips
	for _, sdkIp := range sdkIps {
		ip, err := adaptIp(sdkIp)
		if err != nil {
			return nil, fmt.Errorf("adaptIps: %w", err)
		}
		ips = append(ips, *ip)
	}

	return ips, nil
}

func adaptIpDetails(sdkIp sdkModel.IpDetails) (*domainEntity.Ip, error) {
	networkType, err := enum.NewNetworkType(string(sdkIp.GetNetworkType()))
	if err != nil {
		return nil, fmt.Errorf("adaptIpDetails: %w", err)
	}

	optionalIpValues := domainEntity.OptionalIpValues{
		ReverseLookup: shared.AdaptNullableStringToValue(sdkIp.ReverseLookup),
	}

	sdkDdos, _ := sdkIp.GetDdosOk()
	if sdkDdos != nil {
		ddos := adaptDdos(*sdkDdos)
		optionalIpValues.Ddos = &ddos
	}

	ip := domainEntity.NewIp(
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

func adaptIp(sdkIp sdkModel.Ip) (*domainEntity.Ip, error) {
	networkType, err := enum.NewNetworkType(string(sdkIp.GetNetworkType()))
	if err != nil {
		return nil, fmt.Errorf(
			"adaptIpDetails: %w",
			err,
		)
	}

	optionalIpValues := domainEntity.OptionalIpValues{
		ReverseLookup: shared.AdaptNullableStringToValue(sdkIp.ReverseLookup),
	}

	ip := domainEntity.NewIp(
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

func adaptDdos(sdkDdos sdkModel.Ddos) domainEntity.Ddos {
	return domainEntity.NewDdos(
		sdkDdos.GetDetectionProfile(),
		sdkDdos.GetProtectionType(),
	)
}

func adaptContract(sdkContract sdkModel.Contract) (*domainEntity.Contract, error) {
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

	contract, err := domainEntity.NewContract(
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

func adaptIso(sdkIso sdkModel.Iso) domainEntity.Iso {
	return domainEntity.NewIso(sdkIso.GetId(), sdkIso.GetName())
}

func adaptPrivateNetwork(sdkPrivateNetwork sdkModel.PrivateNetwork) domainEntity.PrivateNetwork {
	return domainEntity.PrivateNetwork{
		Id:     sdkPrivateNetwork.GetPrivateNetworkId(),
		Status: sdkPrivateNetwork.GetStatus(),
		Subnet: sdkPrivateNetwork.GetSubnet(),
	}
}

// AdaptAutoScalingGroupDetails adapts publicCloud.AutoScalingGroupDetails to public_cloud.AutoScalingGroup.
func AdaptAutoScalingGroupDetails(sdkAutoScalingGroup sdkModel.AutoScalingGroupDetails) (
	*domainEntity.AutoScalingGroup,
	error,
) {
	var loadBalancer *domainEntity.LoadBalancer

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

	options := domainEntity.AutoScalingGroupOptions{
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

	autoScalingGroup := domainEntity.NewAutoScalingGroup(
		sdkAutoScalingGroup.GetId(),
		autoScalingGroupType,
		state,
		domainEntity.Region{Name: sdkAutoScalingGroup.GetRegion()},
		*reference,
		sdkAutoScalingGroup.GetCreatedAt(),
		sdkAutoScalingGroup.GetUpdatedAt(),
		options,
	)

	return &autoScalingGroup, nil
}

// AdaptLoadBalancerDetails adapts publicCloud.LoadBalancerDetails to public_cloud.LoadBalancer.
func AdaptLoadBalancerDetails(sdkLoadBalancer sdkModel.LoadBalancerDetails) (
	*domainEntity.LoadBalancer,
	error,
) {
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

	options := domainEntity.OptionalLoadBalancerValues{
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

	loadBalancer := domainEntity.NewLoadBalancer(
		sdkLoadBalancer.GetId(),
		domainEntity.InstanceType{Name: string(sdkLoadBalancer.GetType())},
		adaptResources(sdkLoadBalancer.GetResources()),
		domainEntity.Region{Name: sdkLoadBalancer.GetRegion()},
		state,
		*contract,
		ips,
		options,
	)

	return &loadBalancer, nil
}

func adaptLoadBalancer(sdkLoadBalancer sdkModel.LoadBalancer) (
	*domainEntity.LoadBalancer,
	error,
) {
	state, err := enum.NewState(string(sdkLoadBalancer.GetState()))
	if err != nil {
		return nil, fmt.Errorf("adaptLoadBalancer: %w", err)
	}

	options := domainEntity.OptionalLoadBalancerValues{
		Reference: shared.AdaptNullableStringToValue(sdkLoadBalancer.Reference),
		StartedAt: shared.AdaptNullableTimeToValue(sdkLoadBalancer.StartedAt),
	}

	loadBalancer := domainEntity.NewLoadBalancer(
		sdkLoadBalancer.GetId(),
		domainEntity.InstanceType{Name: string(sdkLoadBalancer.GetType())},
		adaptResources(sdkLoadBalancer.GetResources()),
		domainEntity.Region{},
		state,
		domainEntity.Contract{},
		domainEntity.Ips{},
		options,
	)

	return &loadBalancer, nil
}

func adaptLoadBalancerConfiguration(sdkLoadBalancerConfiguration sdkModel.LoadBalancerConfiguration) (
	*domainEntity.LoadBalancerConfiguration,
	error,
) {
	balance, err := enum.NewBalance(string(sdkLoadBalancerConfiguration.GetBalance()))
	if err != nil {
		return nil, fmt.Errorf("adaptLoadBalancerConfiguration: %w", err)
	}

	options := domainEntity.OptionalLoadBalancerConfigurationOptions{
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

	configuration := domainEntity.NewLoadBalancerConfiguration(
		balance,
		sdkLoadBalancerConfiguration.GetXForwardedFor(),
		int(sdkLoadBalancerConfiguration.GetIdleTimeOut()),
		int(sdkLoadBalancerConfiguration.GetTargetPort()),
		options,
	)

	return &configuration, nil
}

func adaptStickySession(sdkStickySession sdkModel.StickySession) domainEntity.StickySession {
	return domainEntity.NewStickySession(
		sdkStickySession.GetEnabled(),
		int(sdkStickySession.GetMaxLifeTime()),
	)
}

func adaptHealthCheck(sdkHealthCheck sdkModel.HealthCheck) (
	*domainEntity.HealthCheck,
	error,
) {
	method, err := enum.NewMethod(sdkHealthCheck.GetMethod())
	if err != nil {
		return nil, fmt.Errorf("adaptHealthCheck: %w", err)
	}

	healthCheck := domainEntity.NewHealthCheck(
		method,
		sdkHealthCheck.GetUri(),
		int(sdkHealthCheck.GetPort()),
		domainEntity.OptionalHealthCheckValues{
			Host: shared.AdaptNullableStringToValue(sdkHealthCheck.Host),
		},
	)

	return &healthCheck, nil
}

func AdaptInstanceType(sdkInstanceType sdkModel.InstanceType) (
	*domainEntity.InstanceType,
	error,
) {
	resources := adaptResources(sdkInstanceType.GetResources())
	prices := adaptPrices(sdkInstanceType.GetPrices())

	optional := domainEntity.OptionalInstanceTypeValues{}

	sdkStorageTypes, _ := sdkInstanceType.GetStorageTypesOk()
	if sdkStorageTypes != nil {
		storageTypes, err := adaptStorageTypes(sdkStorageTypes)
		if err != nil {
			return nil, fmt.Errorf("AdaptInstanceType: %w", err)
		}
		optional.StorageTypes = storageTypes
	}

	instanceType := domainEntity.NewInstanceType(
		sdkInstanceType.GetName(),
		resources,
		prices,
		optional,
	)

	return &instanceType, nil
}

func adaptPrices(sdkPrices sdkModel.Prices) domainEntity.Prices {
	return domainEntity.NewPrices(
		sdkPrices.GetCurrency(),
		sdkPrices.GetCurrencySymbol(),
		adaptPrice(sdkPrices.GetCompute()),
		adaptStorage(sdkPrices.GetStorage()),
	)
}

func adaptStorage(sdkStorage sdkModel.Storage) domainEntity.Storage {
	return domainEntity.NewStorage(
		adaptPrice(sdkStorage.Local),
		adaptPrice(sdkStorage.Central),
	)
}

func adaptPrice(sdkPrice sdkModel.Price) domainEntity.Price {
	return domainEntity.NewPrice(sdkPrice.GetHourlyPrice(), sdkPrice.GetMonthlyPrice())
}

func adaptStorageTypes(sdkStorageTypes []sdkModel.RootDiskStorageType) (
	*domainEntity.StorageTypes,
	error,
) {
	var storageTypes domainEntity.StorageTypes

	for _, sdkStorageType := range sdkStorageTypes {
		storageType, err := enum.NewRootDiskStorageType(string(sdkStorageType))
		if err != nil {
			return nil, fmt.Errorf("adaptStorageTypes: %w", err)
		}
		storageTypes = append(storageTypes, storageType)
	}

	return &storageTypes, nil
}

func AdaptRegion(sdkRegion sdkModel.Region) domainEntity.Region {
	return domainEntity.NewRegion(sdkRegion.GetName(), sdkRegion.GetLocation())
}

func adaptAutoScalingGroup(sdkAutoScalingGroup sdkModel.AutoScalingGroup) (
	*domainEntity.AutoScalingGroup,
	error,
) {
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

	options := domainEntity.AutoScalingGroupOptions{
		DesiredAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.DesiredAmount),
		MinimumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MinimumAmount),
		MaximumAmount: shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.MaximumAmount),
		CpuThreshold:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CpuThreshold),
		CoolDownTime:  shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.CooldownTime),
		StartsAt:      shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.StartsAt),
		EndsAt:        shared.AdaptNullableTimeToValue(sdkAutoScalingGroup.EndsAt),
		WarmupTime:    shared.AdaptNullableInt32ToValue(sdkAutoScalingGroup.WarmupTime),
	}

	autoScalingGroup := domainEntity.NewAutoScalingGroup(
		sdkAutoScalingGroup.GetId(),
		autoScalingGroupType,
		state,
		domainEntity.Region{Name: sdkAutoScalingGroup.GetRegion()},
		*reference,
		sdkAutoScalingGroup.GetCreatedAt(),
		sdkAutoScalingGroup.GetUpdatedAt(),
		options,
	)

	return &autoScalingGroup, nil
}

func adaptVolume(sdkVolume sdkModel.Volume) domainEntity.Volume {
	return domainEntity.NewVolume(float64(sdkVolume.GetSize()), sdkVolume.GetUnit())
}

// AdaptImageDetails adapts publicCloud.ImageDetails to public_cloud.Image.
func AdaptImageDetails(sdkImageDetails sdkModel.ImageDetails) domainEntity.Image {
	var region *domainEntity.Region

	state, _ := sdkImageDetails.GetStateOk()
	stateReason, _ := sdkImageDetails.GetStateReasonOk()
	regionName, _ := sdkImageDetails.GetRegionOk()
	createdAt, _ := sdkImageDetails.GetCreatedAtOk()
	updatedAt, _ := sdkImageDetails.GetUpdatedAtOk()
	storageSize := adaptStorageSize(sdkImageDetails.GetStorageSize())
	version, _ := sdkImageDetails.GetVersionOk()
	architecture, _ := sdkImageDetails.GetArchitectureOk()

	if regionName != nil {
		region = &domainEntity.Region{}
		region.Name = *regionName
	}

	return domainEntity.NewImage(
		sdkImageDetails.GetId(),
		sdkImageDetails.GetName(),
		version,
		sdkImageDetails.GetFamily(),
		sdkImageDetails.GetFlavour(),
		architecture,
		state,
		stateReason,
		region,
		createdAt,
		updatedAt,
		sdkImageDetails.GetCustom(),
		&storageSize,
		sdkImageDetails.GetMarketApps(),
		sdkImageDetails.GetStorageTypes(),
	)
}

func adaptStorageSize(sdkStorageSize sdkModel.StorageSize) domainEntity.StorageSize {
	return domainEntity.NewStorageSize(
		float64(sdkStorageSize.GetSize()),
		sdkStorageSize.GetUnit(),
	)
}
