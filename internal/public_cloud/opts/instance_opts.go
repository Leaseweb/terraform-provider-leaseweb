package opts

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
)

type InstanceOpts struct {
	instance model.Instance
	ctx      context.Context
}

func (o *InstanceOpts) NewUpdateInstanceOpts() (*publicCloud.UpdateInstanceOpts, error) {
	opts := publicCloud.NewUpdateInstanceOpts()
	err := o.setOptionalUpdateInstanceOpts(opts)

	if err != nil {
		return nil, err
	}

	return opts, nil
}

func (o *InstanceOpts) setOptionalUpdateInstanceOpts(
	opts *publicCloud.UpdateInstanceOpts,
) *OptsError {
	contract := model.Contract{}
	o.instance.Contract.As(o.ctx, &contract, basetypes.ObjectAsOptions{})

	if !o.instance.Type.IsNull() && !o.instance.Type.IsUnknown() {
		instanceTypeName, err := publicCloud.NewInstanceTypeNameFromValue(
			o.instance.Type.ValueString(),
		)
		if err != nil {
			return cannotSetInstanceType(o.instance.Type.ValueString())
		}

		opts.SetType(*instanceTypeName)
	}

	if !o.instance.Reference.IsNull() && !o.instance.Reference.IsUnknown() {
		opts.SetReference(o.instance.Reference.ValueString())
	}

	if !o.instance.RootDiskSize.IsNull() && !o.instance.RootDiskSize.IsUnknown() {
		opts.SetRootDiskSize(int32(o.instance.RootDiskSize.ValueInt64()))
	}

	if !contract.Type.IsNull() && !contract.Type.IsUnknown() {
		contractType, err := publicCloud.NewContractTypeFromValue(
			contract.Type.ValueString(),
		)
		if err != nil {
			return cannotSetContractType(contract.Type.ValueString())
		}

		opts.SetContractType(*contractType)
	}

	if !contract.Term.IsNull() && !contract.Term.IsUnknown() {
		contractTerm, err := publicCloud.NewContractTermFromValue(
			int32(contract.Term.ValueInt64()),
		)
		if err != nil {
			return cannotSetContractTerm(contract.Term.ValueInt64())
		}

		opts.SetContractTerm(*contractTerm)
	}

	if !contract.BillingFrequency.IsNull() && !contract.BillingFrequency.IsUnknown() {
		billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
			int32(contract.BillingFrequency.ValueInt64()),
		)
		if err != nil {
			return cannotSetContractBillingFrequency(contract.BillingFrequency.ValueInt64())
		}

		opts.SetBillingFrequency(*billingFrequency)
	}

	return nil
}

func (o *InstanceOpts) NewLaunchInstanceOpts() (
	*publicCloud.LaunchInstanceOpts,
	*OptsError,
) {
	contract := model.Contract{}
	o.instance.Contract.As(o.ctx, &contract, basetypes.ObjectAsOptions{})

	imageId, err := publicCloud.NewImageIdFromValue(
		strings.Trim(
			o.instance.Image.Attributes()["id"].String(),
			"\"",
		),
	)
	if err != nil {
		return nil, cannotSetOperatingSystemId(
			o.instance.Image.Attributes()["id"].String(),
		)
	}

	instanceTypeName, err := publicCloud.NewInstanceTypeNameFromValue(
		o.instance.Type.ValueString(),
	)
	if err != nil {
		return nil, cannotSetInstanceType(o.instance.Type.ValueString())
	}

	rootDiskStorageType, err := publicCloud.NewRootDiskStorageTypeFromValue(
		o.instance.RootDiskStorageType.ValueString(),
	)
	if err != nil {
		return nil, cannotSetRootDiskStorageType(o.instance.RootDiskStorageType.ValueString())
	}

	term, err := publicCloud.NewContractTermFromValue(int32(contract.Term.ValueInt64()))
	if err != nil {
		return nil, cannotSetContractTerm(contract.Term.ValueInt64())
	}

	billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(int32(contract.BillingFrequency.ValueInt64()))
	if err != nil {
		return nil, cannotSetContractBillingFrequency(contract.BillingFrequency.ValueInt64())
	}

	contractType, err := publicCloud.NewContractTypeFromValue(contract.Type.ValueString())
	if err != nil {
		return nil, cannotSetContractType(contract.Type.ValueString())
	}

	opts := publicCloud.NewLaunchInstanceOpts(
		o.instance.Region.ValueString(),
		*instanceTypeName,
		*imageId,
		*contractType,
		*term,
		*billingFrequency,
		*rootDiskStorageType,
	)

	o.setOptionalLaunchInstanceOpts(opts)

	return opts, nil
}

func (o *InstanceOpts) setOptionalLaunchInstanceOpts(
	opts *publicCloud.LaunchInstanceOpts,
) {
	if !o.instance.MarketAppId.IsNull() && !o.instance.MarketAppId.IsUnknown() {
		opts.SetMarketAppId(o.instance.MarketAppId.ValueString())
	}
	if !o.instance.Reference.IsNull() && !o.instance.Reference.IsUnknown() {
		opts.SetReference(o.instance.Reference.ValueString())
	}
	if !o.instance.RootDiskSize.IsNull() && !o.instance.RootDiskSize.IsUnknown() {
		opts.SetRootDiskSize(int32(o.instance.RootDiskSize.ValueInt64()))
	}
	if !o.instance.SshKey.IsNull() && !o.instance.SshKey.IsUnknown() {
		opts.SetSshKey(o.instance.SshKey.ValueString())
	}
}

func NewInstanceOpts(
	instance model.Instance,
	ctx context.Context,
) InstanceOpts {
	return InstanceOpts{instance: instance, ctx: ctx}
}
