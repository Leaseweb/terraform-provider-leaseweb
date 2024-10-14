// Package to_domain_entity implements adapters to convert public_cloud sdk models to domain entities.
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

	rootDiskStorageType, err := enum.NewStorageType(
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
		string(sdkInstance.GetRegion()),
		adaptResources(sdkInstance.GetResources()),
		adaptImage(sdkInstance.GetImage()),
		state,
		sdkInstance.GetProductType(),
		sdkInstance.GetHasPublicIpV4(),
		sdkInstance.GetIncludesPrivateNetwork(),
		sdkInstance.GetHasUserData(),
		*rootDiskSize,
		string(sdkInstance.GetType()),
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

	rootDiskStorageType, err := enum.NewStorageType(
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

	instance := domainEntity.NewInstance(
		sdkInstanceDetails.GetId(),
		string(sdkInstanceDetails.GetRegion()),
		adaptResources(sdkInstanceDetails.GetResources()),
		adaptImage(sdkInstanceDetails.GetImage()),
		state,
		sdkInstanceDetails.GetProductType(),
		sdkInstanceDetails.GetHasPublicIpV4(),
		sdkInstanceDetails.GetIncludesPrivateNetwork(),
		sdkInstanceDetails.GetHasUserData(),
		*rootDiskSize,
		string(sdkInstanceDetails.GetType()),
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
		sdkImage.GetFamily(),
		sdkImage.GetFlavour(),
		sdkImage.GetCustom(),
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
		string(sdkAutoScalingGroup.GetRegion()),
		*reference,
		sdkAutoScalingGroup.GetCreatedAt(),
		sdkAutoScalingGroup.GetUpdatedAt(),
		options,
	)

	return &autoScalingGroup, nil
}
