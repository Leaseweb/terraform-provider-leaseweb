package to_instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/handlers/shared"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func AdaptToCreateInstanceOpts(
	instanceResourceModel model.Instance,
	allowedInstancedTypes []string,
	ctx context.Context,
) (*domain.Instance, error) {
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
			imageDiags,
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

	imageId, err := enum.NewImageId(image.Id.ValueString())
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

	if instanceResourceModel.RootDiskSize.ValueInt64() != 0 {
		rootDiskSize, err = value_object.NewRootDiskSize(
			int(instanceResourceModel.RootDiskSize.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToCreateInstanceOpts: %w",
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
			"AdaptToCreateInstanceOpts: %w",
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
			MarketAppId: shared.AdaptStringPointerValueToNullableString(
				instanceResourceModel.MarketAppId,
			),
			Reference: shared.AdaptStringPointerValueToNullableString(
				instanceResourceModel.Reference,
			),
			SshKey:       sshKey,
			RootDiskSize: rootDiskSize,
		},
	)

	return &createInstanceOpts, nil
}

func AdaptToUpdateInstanceOpts(
	instanceResourceModel model.Instance,
	allowedInstanceTypes []string,
	ctx context.Context,
) (*domain.Instance, error) {

	id, err := value_object.NewUuid(instanceResourceModel.Id.ValueString())
	if err != nil {
		return nil, fmt.Errorf(
			"AdaptToUpdateInstanceOpts: %w",
			err,
		)
	}

	optionalValues := domain.OptionalUpdateInstanceValues{
		Reference: shared.AdaptStringPointerValueToNullableString(
			instanceResourceModel.Reference,
		),
	}

	if instanceResourceModel.RootDiskSize.ValueInt64() != 0 {
		rootDiskSize, err := value_object.NewRootDiskSize(
			int(instanceResourceModel.RootDiskSize.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToUpdateInstanceOpts: %w",
				err,
			)
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

	if instanceResourceModel.Type.ValueString() != "" {
		instanceType, err := value_object.NewInstanceType(
			instanceResourceModel.Type.ValueString(),
			allowedInstanceTypes,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"AdaptToUpdateInstanceOpts: %w",
				err,
			)
		}
		optionalValues.Type = instanceType
	}

	instance := domain.NewUpdateInstance(*id, optionalValues)

	return &instance, nil
}

func generateInstanceModel(
	rootDiskStorageType *string,
	imageId *string,
	contractType *string,
	contractTerm *int,
	billingFrequency *int,
	sshKey *string,
	rootDiskSize *int,
	instanceType *string,
) model.Instance {
	defaultRootDiskStorageType := "CENTRAL"
	defaultImageId := "UBUNTU_20_04_64BIT"
	defaultContractType := "MONTHLY"
	defaultContractTerm := 3
	defaultBillingFrequency := 1
	defaultRootDiskSize := 55
	defaultInstanceType := "lsw.m5a.4xlarge"

	if rootDiskStorageType == nil {
		rootDiskStorageType = &defaultRootDiskStorageType
	}
	if imageId == nil {
		imageId = &defaultImageId
	}
	if contractType == nil {
		contractType = &defaultContractType
	}
	if contractTerm == nil {
		contractTerm = &defaultContractTerm
	}
	if billingFrequency == nil {
		billingFrequency = &defaultBillingFrequency
	}
	if rootDiskSize == nil {
		rootDiskSize = &defaultRootDiskSize
	}
	if sshKey == nil {
		sshKey = &defaultSshKey
	}
	if instanceType == nil {
		instanceType = &defaultInstanceType
	}

	image, _ := types.ObjectValueFrom(
		context.TODO(),
		model.Image{}.AttributeTypes(),
		model.Image{
			Id:           basetypes.NewStringValue(*imageId),
			Name:         basetypes.NewStringUnknown(),
			Version:      basetypes.NewStringUnknown(),
			Family:       basetypes.NewStringUnknown(),
			Flavour:      basetypes.NewStringUnknown(),
			Architecture: basetypes.NewStringUnknown(),
			MarketApps:   basetypes.NewListUnknown(types.StringType),
			StorageTypes: basetypes.NewListUnknown(types.StringType),
		},
	)

	contract, _ := types.ObjectValueFrom(
		context.TODO(),
		model.Contract{}.AttributeTypes(),
		model.Contract{
			BillingFrequency: basetypes.NewInt64Value(int64(*billingFrequency)),
			Term:             basetypes.NewInt64Value(int64(*contractTerm)),
			Type:             basetypes.NewStringValue(*contractType),
			EndsAt:           basetypes.NewStringUnknown(),
			RenewalsAt:       basetypes.NewStringUnknown(),
			CreatedAt:        basetypes.NewStringUnknown(),
			State:            basetypes.NewStringUnknown(),
		},
	)

	instance := model.Instance{
		Id: basetypes.NewStringValue(
			value_object.NewGeneratedUuid().String(),
		),
		Region:              basetypes.NewStringValue("region"),
		Type:                basetypes.NewStringValue(*instanceType),
		RootDiskStorageType: basetypes.NewStringValue(*rootDiskStorageType),
		RootDiskSize:        basetypes.NewInt64Value(int64(*rootDiskSize)),
		Image:               image,
		Contract:            contract,
		MarketAppId:         basetypes.NewStringValue("marketAppId"),
		Reference:           basetypes.NewStringValue("reference"),
		SshKey:              basetypes.NewStringValue(*sshKey),
	}

	return instance
}
