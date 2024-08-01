package to_sdk_model

import (
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain"
)

// AdaptToLaunchInstanceOpts adapts an instance domain entity to supported launch instance opts.
func AdaptToLaunchInstanceOpts(instance domain.Instance) (
	*publicCloud.LaunchInstanceOpts,
	error,
) {
	instanceTypeName, err := publicCloud.NewTypeNameFromValue(
		instance.Type.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptToLaunchInstanceOpts: %w", err)
	}

	rootDiskStorageType, err := publicCloud.NewRootDiskStorageTypeFromValue(
		instance.RootDiskStorageType.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptToLaunchInstanceOpts: %w", err)
	}

	contractType, err := publicCloud.NewContractTypeFromValue(
		instance.Contract.Type.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptToLaunchInstanceOpts: %w", err)
	}

	contractTerm, err := publicCloud.NewContractTermFromValue(
		int32(instance.Contract.Term.Value()),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptToLaunchInstanceOpts: %w", err)
	}

	billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
		int32(instance.Contract.BillingFrequency.Value()),
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptToLaunchInstanceOpts: %w", err)
	}

	launchInstanceOpts := publicCloud.NewLaunchInstanceOpts(
		instance.Region,
		*instanceTypeName,
		instance.Image.Id,
		*contractType,
		*contractTerm,
		*billingFrequency,
		*rootDiskStorageType,
	)
	launchInstanceOpts.MarketAppId = instance.MarketAppId
	launchInstanceOpts.Reference = instance.Reference

	if instance.SshKey != nil {
		sshKey := instance.SshKey.String()
		launchInstanceOpts.SshKey = &sshKey
	}

	return launchInstanceOpts, nil
}

// AdaptToUpdateInstanceOpts adapts an instance domain entity to supported update instance opts.
func AdaptToUpdateInstanceOpts(instance domain.Instance) (
	*publicCloud.UpdateInstanceOpts,
	error,
) {
	updateInstanceOpts := publicCloud.NewUpdateInstanceOpts()
	updateInstanceOpts.Reference = instance.Reference

	if instance.RootDiskSize.Value != 0 {
		rootDiskSize := int32(instance.RootDiskSize.Value)
		updateInstanceOpts.RootDiskSize = &rootDiskSize
	}

	if instance.Type.String() != "" {
		instanceTypeName, err := publicCloud.NewTypeNameFromValue(
			instance.Type.String(),
		)
		if err != nil {
			return nil, fmt.Errorf("AdaptToUpdateInstanceOpts: %w", err)
		}
		updateInstanceOpts.Type = instanceTypeName
	}

	if instance.Contract.Type != "" {
		contractType, err := publicCloud.NewContractTypeFromValue(instance.Contract.Type.String())
		if err != nil {
			return nil, fmt.Errorf("AdaptToUpdateInstanceOpts: %w", err)
		}
		updateInstanceOpts.ContractType = contractType
	}

	if instance.Contract.Term != 0 {
		contractTerm, err := publicCloud.NewContractTermFromValue(
			int32(instance.Contract.Term.Value()),
		)
		if err != nil {
			return nil, fmt.Errorf("AdaptToUpdateInstanceOpts: %w", err)
		}
		updateInstanceOpts.ContractTerm = contractTerm
	}

	if instance.Contract.BillingFrequency != 0 {
		billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
			int32(instance.Contract.BillingFrequency.Value()),
		)
		if err != nil {
			return nil, fmt.Errorf("AdaptToUpdateInstanceOpts: %w", err)
		}
		updateInstanceOpts.BillingFrequency = billingFrequency
	}

	return updateInstanceOpts, nil
}
