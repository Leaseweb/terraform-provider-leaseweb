package opts

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
)

type InstanceOpts struct {
	instance model.Instance
	ctx      context.Context
}

func (o *InstanceOpts) NewUpdateInstanceOpts() *publicCloud.UpdateInstanceOpts {
	opts := publicCloud.NewUpdateInstanceOpts()
	o.setOptionalUpdateInstanceOpts(opts)

	return opts
}

func (o *InstanceOpts) setOptionalUpdateInstanceOpts(
	opts *publicCloud.UpdateInstanceOpts,
) {
	contract := model.Contract{}
	o.instance.Contract.As(o.ctx, &contract, basetypes.ObjectAsOptions{})

	if !o.instance.Type.IsUnknown() {
		opts.SetType(o.instance.Type.ValueString())
	}

	if !o.instance.Reference.IsUnknown() {
		opts.SetReference(o.instance.Reference.ValueString())
	}

	if !o.instance.RootDiskSize.IsUnknown() {
		opts.SetRootDiskSize(int32(o.instance.RootDiskSize.ValueInt64()))
	}

	if !contract.Type.IsUnknown() {
		opts.SetContractType(contract.Type.ValueString())
	}
	if !contract.Term.IsUnknown() {
		opts.SetContractTerm(int32(contract.Term.ValueInt64()))
	}
	if !contract.BillingFrequency.IsUnknown() {
		opts.SetBillingFrequency(int32(contract.BillingFrequency.ValueInt64()))
	}
}

func (o *InstanceOpts) NewLaunchInstanceOpts() *publicCloud.LaunchInstanceOpts {

	contract := model.Contract{}
	o.instance.Contract.As(o.ctx, &contract, basetypes.ObjectAsOptions{})

	opts := publicCloud.NewLaunchInstanceOpts(
		o.instance.Region.ValueString(),
		o.instance.Type.ValueString(),
		o.instance.OperatingSystem.Attributes()["id"].String(),
		contract.Type.ValueString(),
		int32(contract.Term.ValueInt64()),
		int32(contract.BillingFrequency.ValueInt64()),
		o.instance.RootDiskStorageType.ValueString(),
	)

	o.setOptionalLaunchInstanceOpts(opts)

	return opts
}

func (o *InstanceOpts) setOptionalLaunchInstanceOpts(
	opts *publicCloud.LaunchInstanceOpts,
) {
	if !o.instance.MarketAppId.IsUnknown() {
		opts.SetMarketAppId(o.instance.MarketAppId.ValueString())
	}
	if !o.instance.Reference.IsUnknown() {
		opts.SetReference(o.instance.Reference.ValueString())
	}
	if !o.instance.RootDiskSize.IsUnknown() {
		opts.SetRootDiskSize(int32(o.instance.RootDiskSize.ValueInt64()))
	}
	if !o.instance.SshKey.IsUnknown() {
		opts.SetSshKey(o.instance.SshKey.ValueString())
	}
}

func NewInstanceOpts(
	instance model.Instance,
	ctx context.Context,
) InstanceOpts {
	return InstanceOpts{instance: instance, ctx: ctx}
}
