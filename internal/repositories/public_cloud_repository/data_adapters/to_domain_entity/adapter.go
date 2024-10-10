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

	ips := adaptIps(sdkInstance.GetIps())

	contract, err := adaptContract(sdkInstance.GetContract())
	if err != nil {
		return nil, fmt.Errorf("AdaptInstance:  %w", err)
	}

	optionalValues := domainEntity.OptionalInstanceValues{
		Reference:   shared.AdaptNullableStringToValue(sdkInstance.Reference),
		MarketAppId: shared.AdaptNullableStringToValue(sdkInstance.MarketAppId),
		StartedAt:   shared.AdaptNullableTimeToValue(sdkInstance.StartedAt),
	}

	instance := domainEntity.NewInstance(
		sdkInstance.GetId(),
		string(sdkInstance.GetRegion()),
		adaptImage(sdkInstance.GetImage()),
		state,
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

	ips := adaptIpsDetails(sdkInstanceDetails.GetIps())

	contract, err := adaptContract(sdkInstanceDetails.GetContract())
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails:  %w", err)
	}

	optionalValues := domainEntity.OptionalInstanceValues{
		Reference: shared.AdaptNullableStringToValue(
			sdkInstanceDetails.Reference,
		),
		MarketAppId: shared.AdaptNullableStringToValue(
			sdkInstanceDetails.MarketAppId,
		),
		StartedAt: shared.AdaptNullableTimeToValue(sdkInstanceDetails.StartedAt),
	}

	instance := domainEntity.NewInstance(
		sdkInstanceDetails.GetId(),
		string(sdkInstanceDetails.GetRegion()),
		adaptImage(sdkInstanceDetails.GetImage()),
		state,
		*rootDiskSize,
		string(sdkInstanceDetails.GetType()),
		rootDiskStorageType,
		ips,
		*contract,
		optionalValues,
	)

	return &instance, nil
}

func adaptImage(sdkImage sdkModel.Image) domainEntity.Image {
	return domainEntity.NewImage(sdkImage.GetId())
}

func adaptIpsDetails(sdkIps []sdkModel.IpDetails) domainEntity.Ips {
	var ips domainEntity.Ips
	for _, sdkIp := range sdkIps {
		ips = append(ips, adaptIpDetails(sdkIp))
	}

	return ips
}

func adaptIps(sdkIps []sdkModel.Ip) domainEntity.Ips {
	var ips domainEntity.Ips
	for _, sdkIp := range sdkIps {
		ips = append(ips, adaptIp(sdkIp))
	}

	return ips
}

func adaptIpDetails(sdkIp sdkModel.IpDetails) domainEntity.Ip {
	return domainEntity.NewIp(sdkIp.GetIp())
}

func adaptIp(sdkIp sdkModel.Ip) domainEntity.Ip {
	return domainEntity.NewIp(sdkIp.GetIp())
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
		contractState,
		shared.AdaptNullableTimeToValue(sdkContract.EndsAt),
	)

	if err != nil {
		return nil, fmt.Errorf("adaptContract: %w", err)
	}

	return contract, nil
}
