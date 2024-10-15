package to_opts

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/dataadapters/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
)

func AdaptToLaunchInstanceOpts(
	instance resource.Instance,
	ctx context.Context,
) (*publicCloud.LaunchInstanceOpts, error) {
	sdkRootDiskStorageType, err := publicCloud.NewStorageTypeFromValue(
		instance.RootDiskStorageType.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	image := resource.Image{}
	imageDiags := instance.Image.As(
		ctx,
		&image,
		basetypes.ObjectAsOptions{},
	)
	if imageDiags != nil {
		return nil, shared.ReturnError(
			"AdaptToCreateInstanceOpts",
			imageDiags,
		)
	}

	contract := resource.Contract{}
	contractDiags := instance.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if contractDiags != nil {
		return nil, shared.ReturnError(
			"AdaptToCreateInstanceOpts",
			contractDiags,
		)
	}

	sdkContractType, err := publicCloud.NewContractTypeFromValue(
		contract.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	sdkContractTerm, err := publicCloud.NewContractTermFromValue(
		int32(contract.Term.ValueInt64()),
	)
	if err != nil {
		return nil, err
	}

	sdkBillingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
		int32(contract.BillingFrequency.ValueInt64()),
	)
	if err != nil {
		return nil, err
	}

	sdkRegionName, err := publicCloud.NewRegionNameFromValue(
		instance.Region.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	sdkInstanceType, err := publicCloud.NewTypeNameFromValue(
		instance.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	opts := publicCloud.NewLaunchInstanceOpts(
		*sdkRegionName,
		*sdkInstanceType,
		image.Id.ValueString(),
		*sdkContractType,
		*sdkContractTerm,
		*sdkBillingFrequency,
		*sdkRootDiskStorageType,
	)

	opts.MarketAppId = shared.AdaptStringPointerValueToNullableString(
		instance.MarketAppId,
	)
	opts.Reference = shared.AdaptStringPointerValueToNullableString(
		instance.Reference,
	)
	opts.RootDiskSize = shared.AdaptInt64PointerValueToNullableInt32(
		instance.RootDiskSize,
	)

	return opts, nil
}

func AdaptToUpdateInstanceOpts(
	instance resource.Instance,
	ctx context.Context,
) (*publicCloud.UpdateInstanceOpts, error) {
	opts := publicCloud.NewUpdateInstanceOpts()
	opts.Reference = shared.AdaptStringPointerValueToNullableString(
		instance.Reference,
	)
	opts.RootDiskSize = shared.AdaptInt64PointerValueToNullableInt32(
		instance.RootDiskSize,
	)

	contract := resource.Contract{}
	diags := instance.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if diags.HasError() {
		return nil, shared.ReturnError(
			"AdaptToUpdateInstanceOpts",
			diags,
		)
	}

	if contract.Type.ValueString() != "" {
		contractType, err := publicCloud.NewContractTypeFromValue(
			contract.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToUpdateInstanceOpts: %w",
				err,
			)
		}
		opts.ContractType = contractType
	}

	if contract.Term.ValueInt64() != 0 {
		contractTerm, err := publicCloud.NewContractTermFromValue(
			int32(contract.Term.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToUpdateInstanceOpts: %w",
				err,
			)
		}
		opts.ContractTerm = contractTerm
	}

	if contract.BillingFrequency.ValueInt64() != 0 {
		billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
			int32(contract.BillingFrequency.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToUpdateInstanceOpts: %w",
				err,
			)
		}
		opts.BillingFrequency = billingFrequency
	}

	if instance.Type.ValueString() != "" {
		instanceType, err := publicCloud.NewTypeNameFromValue(
			instance.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToUpdateInstanceOpts: %w",
				err,
			)
		}
		opts.Type = instanceType
	}

	return opts, nil
}
