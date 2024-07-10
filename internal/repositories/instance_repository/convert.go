package instance_repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

var ErrCannotConvertStringToUUid = fmt.Errorf("cannot convert string to uuid")

var ErrNoLoadBalancerDetails = fmt.Errorf("loadBalancer details cannot be found")

func convertInstance(
	sdkInstance publicCloud.InstanceDetails,
	autoScalingGroup *entity.AutoScalingGroup,
) (*entity.Instance, error) {
	instanceId, err := convertStringToUuid(sdkInstance.GetId())
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing Instance id %q: %w",
			sdkInstance.GetId(),
			err,
		)
	}

	image, err := convertImage(sdkInstance.GetImage())
	if err != nil {
		return nil, fmt.Errorf(
			"error converting Image %q: %w",
			sdkInstance.GetImage().Id,
			err,
		)
	}

	state, err := enum.FindEnumForString(
		string(sdkInstance.GetState()),
		enum.StateValues,
		enum.StateUnknown,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing state %q: %w",
			string(sdkInstance.GetState()),
			err,
		)
	}

	rootDiskSize, err := value_object.NewRootDiskSize(int64(sdkInstance.GetRootDiskSize()))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing rootDiskSize %d: %w",
			sdkInstance.GetRootDiskSize(),
			err,
		)
	}

	rootDiskStorageType, err := enum.FindEnumForString(
		string(sdkInstance.GetRootDiskStorageType()),
		enum.RootDiskStorageTypeValues,
		enum.RootDiskStorageTypeCentral,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing rootDiskStorageType %q: %w",
			string(sdkInstance.GetRootDiskStorageType()),
			err,
		)
	}

	ips, err := convertIps(sdkInstance.GetIps())
	if err != nil {
		return nil, fmt.Errorf(
			"error converting ips:  %w",
			err,
		)
	}

	contract, err := convertContract(sdkInstance.GetContract())
	if err != nil {
		return nil, fmt.Errorf(
			"error converting contract:  %w",
			err,
		)
	}

	optionalValues := entity.OptionalInstanceValues{
		Reference:        convertNullableStringToValue(sdkInstance.Reference),
		MarketAppId:      convertNullableStringToValue(sdkInstance.MarketAppId),
		StartedAt:        convertNullableTimeToValue(sdkInstance.StartedAt),
		AutoScalingGroup: autoScalingGroup,
	}
	if sdkInstance.Iso.Get() != nil {
		iso := convertIso(*sdkInstance.Iso.Get())
		optionalValues.Iso = &iso
	}
	if sdkInstance.PrivateNetwork.Get() != nil {
		privateNetwork := convertPrivateNetwork(*sdkInstance.PrivateNetwork.Get())
		optionalValues.PrivateNetwork = &privateNetwork
	}

	instance := entity.NewInstance(
		*instanceId,
		sdkInstance.GetRegion(),
		convertResources(sdkInstance.GetResources()),
		*image,
		state,
		sdkInstance.GetProductType(),
		sdkInstance.GetHasPublicIpV4(),
		sdkInstance.GetIncludesPrivateNetwork(),
		*rootDiskSize,
		string(sdkInstance.GetType()),
		rootDiskStorageType,
		ips,
		*contract,
		optionalValues,
	)

	return &instance, nil
}

func convertResources(sdkResources publicCloud.Resources) entity.Resources {
	resources := entity.NewResources(
		convertCpu(sdkResources.GetCpu()),
		convertMemory(sdkResources.GetMemory()),
		convertNetworkSpeed(sdkResources.GetPublicNetworkSpeed()),
		convertNetworkSpeed(sdkResources.GetPrivateNetworkSpeed()),
	)

	return resources
}

func convertCpu(sdkCpu publicCloud.Cpu) entity.Cpu {
	return entity.NewCpu(int64(sdkCpu.GetValue()), sdkCpu.GetUnit())
}

func convertMemory(sdkMemory publicCloud.Memory) entity.Memory {
	return entity.NewMemory(float64(sdkMemory.GetValue()), sdkMemory.GetUnit())
}

func convertNetworkSpeed(sdkNetworkSpeed publicCloud.NetworkSpeed) entity.NetworkSpeed {
	return entity.NewNetworkSpeed(
		int64(sdkNetworkSpeed.GetValue()),
		sdkNetworkSpeed.GetUnit(),
	)
}

func convertImage(sdkImage publicCloud.ImageDetails) (*entity.Image, error) {
	imageId, err := enum.FindEnumForString(
		string(sdkImage.GetId()),
		enum.ImageIdValues,
		enum.Almalinux864Bit,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert Image %q: %w",
			sdkImage.GetId(),
			err,
		)
	}

	image := entity.NewImage(
		imageId,
		sdkImage.GetName(),
		sdkImage.GetVersion(),
		sdkImage.GetFamily(),
		sdkImage.GetFlavour(),
		sdkImage.GetArchitecture(),
		sdkImage.GetMarketApps(),
		sdkImage.GetStorageTypes(),
	)

	return &image, nil
}

func convertIps(sdkIps []publicCloud.IpDetails) (entity.Ips, error) {
	var ips entity.Ips
	for _, sdkIp := range sdkIps {
		ip, err := convertIp(sdkIp)
		if err != nil {
			return nil, fmt.Errorf("error parsing ip: %w", err)
		}
		ips = append(ips, *ip)
	}

	return ips, nil
}

func convertIp(sdkIp publicCloud.IpDetails) (*entity.Ip, error) {
	networkType, err := enum.FindEnumForString(
		string(sdkIp.GetNetworkType()),
		enum.NetworkTypeValues,
		enum.NetworkTypePublic,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing networkType %q: %w",
			string(sdkIp.GetNetworkType()),
			err,
		)
	}

	optionalIpValues := entity.OptionalIpValues{
		ReverseLookup: convertNullableStringToValue(sdkIp.ReverseLookup),
	}

	sdkDdos, _ := sdkIp.GetDdosOk()
	if sdkDdos != nil {
		ddos := convertDdos(*sdkDdos)
		optionalIpValues.Ddos = &ddos
	}

	ip := entity.NewIp(
		sdkIp.GetIp(),
		sdkIp.GetPrefixLength(),
		int64(sdkIp.GetVersion()),
		sdkIp.GetNullRouted(),
		sdkIp.GetMainIp(),
		networkType,
		optionalIpValues,
	)

	return &ip, nil
}

func convertDdos(sdkDdos publicCloud.Ddos) entity.Ddos {
	return entity.NewDdos(
		sdkDdos.GetDetectionProfile(),
		sdkDdos.GetProtectionType(),
	)
}

func convertContract(sdkContract publicCloud.Contract) (*entity.Contract, error) {
	billingFrequency, err := enum.FindEnumForInt(
		int64(sdkContract.GetBillingFrequency()),
		enum.ContractBillingFrequencyValues,
		enum.ContractBillingFrequencyZero,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing billingFrequency %d: %w",
			sdkContract.GetBillingFrequency(),
			err,
		)
	}

	contractTerm, err := enum.FindEnumForInt(
		int64(sdkContract.GetTerm()),
		enum.ContractTermValues,
		enum.ContractTermZero,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing contractTerm %d: %w",
			sdkContract.GetTerm(),
			err,
		)
	}

	contractType, err := enum.FindEnumForString(
		string(sdkContract.GetType()),
		enum.ContractTypeValues,
		enum.ContractTypeHourly,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing contractType %q: %w",
			string(sdkContract.GetType()),
			err,
		)
	}

	contractState, err := enum.FindEnumForString(
		string(sdkContract.GetState()),
		enum.ContractStateValues,
		enum.ContractStateActive,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing contractState %q: %w",
			string(sdkContract.GetState()),
			err,
		)
	}

	contract, err := entity.NewContract(
		billingFrequency,
		contractTerm,
		contractType,
		sdkContract.GetRenewalsAt(),
		sdkContract.GetCreatedAt(),
		contractState,
		convertNullableTimeToValue(sdkContract.EndsAt),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error parsing contract: %w",
			err,
		)
	}

	return contract, nil
}

func convertIso(sdkIso publicCloud.Iso) entity.Iso {
	return entity.NewIso(sdkIso.GetId(), sdkIso.GetName())
}

func convertPrivateNetwork(sdkPrivateNetwork publicCloud.PrivateNetwork) entity.PrivateNetwork {
	return entity.PrivateNetwork{
		Id:     sdkPrivateNetwork.GetPrivateNetworkId(),
		Status: sdkPrivateNetwork.GetStatus(),
		Subnet: sdkPrivateNetwork.GetSubnet(),
	}
}

func convertAutoScalingGroup(
	sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
	loadBalancer *entity.LoadBalancer,
) (
	*entity.AutoScalingGroup,
	error,
) {
	if sdkAutoScalingGroup.LoadBalancer.Get() != nil && loadBalancer == nil {
		return nil, ErrNoLoadBalancerDetails
	}

	autoScalingGroupId, err := convertStringToUuid(sdkAutoScalingGroup.GetId())
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing autoScalingGroup id %q: %w",
			sdkAutoScalingGroup.GetId(),
			err,
		)
	}

	autoScalingGroupType, err := enum.FindEnumForString(
		string(sdkAutoScalingGroup.GetType()),
		enum.AutoScalingGroupTypeValues,
		enum.AutoScalingCpuTypeManual,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing autoScalingGroupType %q: %w",
			sdkAutoScalingGroup.GetType(),
			err,
		)
	}

	state, err := enum.FindEnumForString(
		string(sdkAutoScalingGroup.GetState()),
		enum.StateValues,
		enum.StateUnknown,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing state %q: %w",
			sdkAutoScalingGroup.GetState(),
			err,
		)
	}

	reference, err := value_object.NewAutoScalingGroupReference(sdkAutoScalingGroup.GetReference())

	if err != nil {
		return nil, fmt.Errorf(
			"error parsing reference %q: %w",
			sdkAutoScalingGroup.GetReference(),
			err,
		)
	}

	options := entity.AutoScalingGroupOptions{
		DesiredAmount: convertNullableInt32ToValue(sdkAutoScalingGroup.DesiredAmount),
		MinimumAmount: convertNullableInt32ToValue(sdkAutoScalingGroup.MinimumAmount),
		MaximumAmount: convertNullableInt32ToValue(sdkAutoScalingGroup.MaximumAmount),
		CpuThreshold:  convertNullableInt32ToValue(sdkAutoScalingGroup.CpuThreshold),
		CoolDownTime:  convertNullableInt32ToValue(sdkAutoScalingGroup.CooldownTime),
		StartsAt:      convertNullableTimeToValue(sdkAutoScalingGroup.StartsAt),
		EndsAt:        convertNullableTimeToValue(sdkAutoScalingGroup.EndsAt),
		WarmupTime:    convertNullableInt32ToValue(sdkAutoScalingGroup.WarmupTime),
		LoadBalancer:  loadBalancer,
	}

	autoScalingGroup := entity.NewAutoScalingGroup(
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

func convertStringToUuid(id string) (*uuid.UUID, error) {
	convertedId, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrCannotConvertStringToUUid
	}

	return &convertedId, nil
}

func convertNullableStringToValue(nullableString publicCloud.NullableString) *string {
	return nullableString.Get()
}

func convertNullableTimeToValue(nullableTime publicCloud.NullableTime) *time.Time {
	return nullableTime.Get()
}

func convertNullableInt32ToValue(nullableInt publicCloud.NullableInt32) *int64 {
	if nullableInt.Get() == nil {
		return nil
	}

	value := int64(*nullableInt.Get())
	return &value
}

func convertLoadBalancer(sdkLoadBalancer publicCloud.LoadBalancerDetails) (
	*entity.LoadBalancer,
	error,
) {
	loadBalancerId, err := convertStringToUuid(sdkLoadBalancer.Id)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing loadBalancer id %q: %w",
			sdkLoadBalancer.GetId(),
			err,
		)
	}

	state, err := enum.FindEnumForString(
		string(sdkLoadBalancer.GetState()),
		enum.StateValues,
		enum.StateUnknown,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing load balancer state %q: %w",
			string(sdkLoadBalancer.GetState()),
			err,
		)
	}

	contract, err := convertContract(sdkLoadBalancer.GetContract())
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing contract: %w",
			err,
		)
	}

	ips, err := convertIps(sdkLoadBalancer.GetIps())
	if err != nil {
		return nil, fmt.Errorf(
			"error converting ips:  %w",
			err,
		)
	}

	options := entity.OptionalLoadBalancerValues{
		Reference: convertNullableStringToValue(sdkLoadBalancer.Reference),
		StartedAt: convertNullableTimeToValue(sdkLoadBalancer.StartedAt),
	}

	if sdkLoadBalancer.Configuration.Get() != nil {
		configuration, err := convertLoadBalancerConfiguration(sdkLoadBalancer.GetConfiguration())
		if err != nil {
			return nil, fmt.Errorf(
				"error converting configuration:  %w",
				err,
			)
		}
		options.Configuration = configuration
	}

	if sdkLoadBalancer.PrivateNetwork.Get() != nil {
		privateNetwork := convertPrivateNetwork(*sdkLoadBalancer.PrivateNetwork.Get())
		options.PrivateNetwork = &privateNetwork
	}

	loadBalancer := entity.NewLoadBalancer(
		*loadBalancerId,
		sdkLoadBalancer.GetType(),
		convertResources(sdkLoadBalancer.GetResources()),
		sdkLoadBalancer.GetRegion(),
		state,
		*contract,
		ips,
		options,
	)

	return &loadBalancer, nil
}

func convertLoadBalancerConfiguration(sdkLoadBalancerConfiguration publicCloud.LoadBalancerConfiguration) (
	*entity.LoadBalancerConfiguration,
	error,
) {
	balance, err := enum.FindEnumForString(
		sdkLoadBalancerConfiguration.GetBalance(),
		enum.BalanceValues,
		enum.BalanceRoundRobin,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing balance %q: %w",
			sdkLoadBalancerConfiguration.GetBalance(),
			err,
		)
	}

	options := entity.OptionalLoadBalancerConfigurationOptions{
		HealthCheck: nil,
	}
	if sdkLoadBalancerConfiguration.StickySession.Get() != nil {
		stickySession := convertStickySession(*sdkLoadBalancerConfiguration.StickySession.Get())
		options.StickySession = &stickySession
	}
	if sdkLoadBalancerConfiguration.HealthCheck.Get() != nil {
		healthCheck, err := convertHealthCheck(*sdkLoadBalancerConfiguration.HealthCheck.Get())
		if err != nil {
			return nil, fmt.Errorf(
				"error parsing healthCheck: %w",
				err,
			)
		}

		options.HealthCheck = healthCheck
	}

	configuration := entity.NewLoadBalancerConfiguration(
		balance,
		sdkLoadBalancerConfiguration.GetXForwardedFor(),
		int64(sdkLoadBalancerConfiguration.GetIdleTimeOut()),
		int64(sdkLoadBalancerConfiguration.GetTargetPort()),
		options,
	)

	return &configuration, nil
}

func convertStickySession(sdkStickySession publicCloud.StickySession) entity.StickySession {
	return entity.NewStickySession(
		sdkStickySession.GetEnabled(),
		int64(sdkStickySession.GetMaxLifeTime()),
	)
}

func convertHealthCheck(sdkHealthCheck publicCloud.HealthCheck) (*entity.HealthCheck, error) {
	method, err := enum.FindEnumForString(
		sdkHealthCheck.GetMethod(),
		enum.MethodValues,
		enum.MethodGet,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing method %q: %w",
			sdkHealthCheck.GetMethod(),
			err,
		)
	}

	healthCheck := entity.NewHealthCheck(
		method,
		sdkHealthCheck.GetUri(),
		int64(sdkHealthCheck.GetPort()),
		entity.OptionalHealthCheckValues{Host: convertNullableStringToValue(sdkHealthCheck.Host)},
	)

	return &healthCheck, nil
}
