package to_domain_entity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

// AdaptToCreateInstanceOpts transforms a resource model to an instance domain
// entity with all supported fields for creating an instance.
func AdaptToCreateInstanceOpts(
	instanceResourceModel model.Instance,
	allowedInstancedTypes []string,
	ctx context.Context,
) (*public_cloud.Instance, error) {
	var sshKey *value_object.SshKey
	var rootDiskSize *value_object.RootDiskSize

	image := model.Image{}
	imageDiags := instanceResourceModel.Image.As(
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

	contract := model.Contract{}
	contractDiags := instanceResourceModel.Contract.As(
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

	instanceType := model.InstanceType{}
	instanceTypeDiags := instanceResourceModel.Type.As(
		ctx,
		&instanceType,
		basetypes.ObjectAsOptions{},
	)
	if instanceTypeDiags != nil {
		return nil, shared.ReturnError(
			"AdaptToCreateInstanceOpts",
			instanceTypeDiags,
		)
	}

	region := model.Region{}
	regionDiags := instanceResourceModel.Region.As(
		ctx,
		&region,
		basetypes.ObjectAsOptions{},
	)
	if regionDiags != nil {
		return nil, shared.ReturnError(
			"AdaptToCreateInstanceOpts",
			regionDiags,
		)
	}

	rootDiskStorageType, err := enum.NewRootDiskStorageType(
		instanceResourceModel.RootDiskStorageType.ValueString(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"AdaptToCreateInstanceOpts: %w",
			err,
		)
	}

	contractType, err := enum.NewContractType(contract.Type.ValueString())
	if err != nil {
		return nil, fmt.Errorf(
			"AdaptToCreateInstanceOpts: %w",
			err,
		)
	}

	contractTerm, err := enum.NewContractTerm(int(contract.Term.ValueInt64()))
	if err != nil {
		return nil, fmt.Errorf(
			"AdaptToCreateInstanceOpts: %w",
			err,
		)
	}

	billingFrequency, err := enum.NewContractBillingFrequency(
		int(contract.BillingFrequency.ValueInt64()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"AdaptToCreateInstanceOpts: %w",
			err,
		)
	}

	// TODO Enable SSH key support
	/**
	  if instanceResourceModel.SshKey.ValueString() != "" {
	  	sshKey, err = value_object.NewSshKey(
	  		instanceResourceModel.SshKey.ValueString(),
	  	)
	  	if err != nil {
	  		return nil, fmt.Errorf(
	  			"AdaptToCreateInstanceOpts: %w",
	  			err,
	  		)
	  	}
	  }
	*/

	if instanceResourceModel.RootDiskSize.ValueInt64() != 0 {
		rootDiskSize, err = value_object.NewRootDiskSize(
			int(instanceResourceModel.RootDiskSize.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf("AdaptToCreateInstanceOpts: %w", err)
		}
	}

	createInstanceOpts, err := public_cloud.NewCreateInstance(
		region.Name.ValueString(),
		instanceType.Name.ValueString(),
		rootDiskStorageType,
		image.Id.ValueString(),
		contractType,
		contractTerm,
		billingFrequency,
		public_cloud.OptionalCreateInstanceValues{
			MarketAppId: shared.AdaptStringPointerValueToNullableString(
				instanceResourceModel.MarketAppId,
			),
			Reference: shared.AdaptStringPointerValueToNullableString(
				instanceResourceModel.Reference,
			),
			SshKey:       sshKey,
			RootDiskSize: rootDiskSize,
		},
		allowedInstancedTypes,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptToCreateInstanceOpts: %w", err)
	}

	return createInstanceOpts, nil
}

// AdaptToUpdateInstanceOpts transforms a resource model to an instance domain
// entity with all supported fields for updating an instance.
func AdaptToUpdateInstanceOpts(
	instanceResourceModel model.Instance,
	allowedInstanceTypes []string,
	currentInstanceType string,
	ctx context.Context,
) (*public_cloud.Instance, error) {
	optionalValues := public_cloud.OptionalUpdateInstanceValues{
		Reference: shared.AdaptStringPointerValueToNullableString(
			instanceResourceModel.Reference,
		),
	}

	if instanceResourceModel.RootDiskSize.ValueInt64() != 0 {
		rootDiskSize, err := value_object.NewRootDiskSize(
			int(instanceResourceModel.RootDiskSize.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf("AdaptToUpdateInstanceOpts: %w", err)
		}
		optionalValues.RootDiskSize = rootDiskSize
	}

	contract := model.Contract{}
	diags := instanceResourceModel.Contract.As(
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

	instanceType := model.InstanceType{}
	diags = instanceResourceModel.Type.As(
		ctx,
		&instanceType,
		basetypes.ObjectAsOptions{},
	)
	if diags.HasError() {
		return nil, shared.ReturnError(
			"AdaptToUpdateInstanceOpts",
			diags,
		)
	}

	if contract.Type.ValueString() != "" {
		contractType, err := enum.NewContractType(contract.Type.ValueString())
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToUpdateInstanceOpts: %w",
				err,
			)
		}
		optionalValues.ContractType = &contractType
	}

	if contract.Term.ValueInt64() != 0 {
		contractTerm, err := enum.NewContractTerm(int(contract.Term.ValueInt64()))
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToUpdateInstanceOpts: %w",
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
				"AdaptToUpdateInstanceOpts: %w",
				err,
			)
		}
		optionalValues.BillingFrequency = &billingFrequency
	}

	if instanceType.Name.ValueString() != "" {
		instanceTypeOpt := instanceType.Name.ValueString()
		optionalValues.Type = &instanceTypeOpt
	}

	instance, err := public_cloud.NewUpdateInstance(
		instanceResourceModel.Id.ValueString(),
		optionalValues,
		allowedInstanceTypes,
		currentInstanceType,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptToUpdateInstanceOpts: %w", err)
	}

	return instance, nil
}
