package model

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
	"terraform-provider-leaseweb/internal/utils"
)

type Instance struct {
	Id                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Resources           types.Object `tfsdk:"resources"`
	Image               types.Object `tfsdk:"image"`
	State               types.String `tfsdk:"state"`
	ProductType         types.String `tfsdk:"product_type"`
	HasPublicIpv4       types.Bool   `tfsdk:"has_public_ipv4"`
	HasPrivateNetwork   types.Bool   `tfsdk:"has_private_network"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int64  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	Ips                 types.List   `tfsdk:"ips"`
	StartedAt           types.String `tfsdk:"started_at"`
	Contract            types.Object `tfsdk:"contract"`
	MarketAppId         types.String `tfsdk:"market_app_id"`
	AutoScalingGroup    types.Object `tfsdk:"auto_scaling_group"`
	Iso                 types.Object `tfsdk:"iso"`
	PrivateNetwork      types.Object `tfsdk:"private_network"`
	SshKey              types.String `tfsdk:"ssh_key"`
}

func (i *Instance) Populate(instance entity.Instance,
	ctx context.Context,
) diag.Diagnostics {
	i.Id = basetypes.NewStringValue(instance.Id.String())
	i.Region = basetypes.NewStringValue(instance.Region)
	i.Reference = utils.ConvertNullableStringToStringValue(instance.Reference)
	i.State = basetypes.NewStringValue(string(instance.State))
	i.ProductType = basetypes.NewStringValue(instance.ProductType)
	i.HasPublicIpv4 = basetypes.NewBoolValue(instance.HasPublicIpv4)
	i.HasPrivateNetwork = basetypes.NewBoolValue(instance.HasPrivateNetwork)
	i.Type = basetypes.NewStringValue(instance.Type)
	i.RootDiskSize = basetypes.NewInt64Value(int64(instance.RootDiskSize.Value))
	i.RootDiskStorageType = basetypes.NewStringValue(string(instance.RootDiskStorageType))
	i.StartedAt = utils.ConvertNullableTimeToStringValue(instance.StartedAt)
	i.MarketAppId = utils.ConvertNullableStringToStringValue(instance.MarketAppId)

	if instance.SshKey != nil {
		i.SshKey = basetypes.NewStringValue(instance.SshKey.String())
	}

	image, diags := utils.ConvertDomainEntityToResourceObject(
		instance.Image,
		Image{}.AttributeTypes(),
		ctx,
		newImage,
	)
	if diags.HasError() {
		return diags
	}
	i.Image = image

	contract, diags := utils.ConvertDomainEntityToResourceObject(
		instance.Contract,
		Contract{}.AttributeTypes(),
		ctx,
		newContract,
	)
	if diags.HasError() {
		return diags
	}
	i.Contract = contract

	iso, diags := utils.ConvertNullableDomainEntityToResourceObject(
		instance.Iso,
		Iso{}.AttributeTypes(),
		ctx,
		newIso,
	)
	if diags.HasError() {
		return diags
	}
	i.Iso = iso

	privateNetwork, diags := utils.ConvertNullableDomainEntityToResourceObject(
		instance.PrivateNetwork,
		PrivateNetwork{}.AttributeTypes(),
		ctx,
		newPrivateNetwork,
	)
	if diags.HasError() {
		return diags
	}
	i.PrivateNetwork = privateNetwork

	resources, diags := utils.ConvertDomainEntityToResourceObject(
		instance.Resources,
		Resources{}.AttributeTypes(),
		ctx,
		newResources,
	)
	if diags.HasError() {
		return diags
	}
	i.Resources = resources

	autoScalingGroup, diags := utils.ConvertNullableDomainEntityToResourceObject(
		instance.AutoScalingGroup,
		AutoScalingGroup{}.AttributeTypes(),
		ctx,
		newAutoScalingGroup,
	)
	if diags.HasError() {
		return diags
	}
	i.AutoScalingGroup = autoScalingGroup

	ips, diags := utils.ConvertEntitiesToListValue(
		instance.Ips,
		Ip{}.AttributeTypes(),
		ctx,
		newIp,
	)
	if diags.HasError() {
		return diags
	}
	i.Ips = ips

	return nil
}

func (i *Instance) GenerateCreateInstanceEntity(ctx context.Context) (
	*entity.Instance,
	diag.Diagnostics,
) {
	var sshKey *value_object.SshKey
	var rootDiskSize *value_object.RootDiskSize
	var diags diag.Diagnostics

	image := Image{}
	diags = i.Image.As(ctx, &image, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	contract := Contract{}
	diags = i.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	rootDiskStorageType, err := enum.NewRootDiskStorageType(i.RootDiskStorageType.ValueString())
	if err != nil {
		diags.AddError(
			fmt.Sprintf(
				"cannot parse rootDisStorageType %q",
				i.RootDiskStorageType.ValueString(),
			),
			err.Error(),
		)
		return nil, diags
	}

	imageId, err := enum.NewImageId(image.Id.ValueString())
	if err != nil {
		diags.AddError(
			fmt.Sprintf(
				"cannot parse imageId %q",
				image.Id.ValueString(),
			),
			err.Error(),
		)
		return nil, diags
	}

	contractType, err := enum.NewContractType(contract.Type.ValueString())
	if err != nil {
		diags.AddError(
			fmt.Sprintf(
				"cannot parse contractType %s",
				contract.Type.ValueString(),
			),
			err.Error(),
		)
		return nil, diags
	}

	contractTerm, err := enum.NewContractTerm(int(contract.Term.ValueInt64()))
	if err != nil {
		diags.AddError(
			fmt.Sprintf(
				"cannot parse contractTerm %d",
				contract.Term.ValueInt64(),
			),
			err.Error(),
		)
		return nil, diags
	}

	billingFrequency, err := enum.NewContractBillingFrequency(int(contract.BillingFrequency.ValueInt64()))
	if err != nil {
		diags.AddError(
			fmt.Sprintf(
				"cannot parse billingFrequency %d",
				contract.BillingFrequency.ValueInt64(),
			),
			err.Error(),
		)
		return nil, diags
	}

	if i.SshKey.ValueString() != "" {
		sshKey, err = value_object.NewSshKey(i.SshKey.ValueString())
		if err != nil {
			diags.AddError(
				fmt.Sprintf(
					"invalid sshKey %q",
					i.SshKey.ValueString(),
				),
				err.Error(),
			)
			return nil, diags
		}
	}

	if i.RootDiskSize.ValueInt64() != 0 {
		rootDiskSize, err = value_object.NewRootDiskSize(int(i.RootDiskSize.ValueInt64()))
		if err != nil {
			diags.AddError(
				fmt.Sprintf(
					"invalid rootDiskSize %d",
					i.RootDiskSize.ValueInt64(),
				),
				err.Error(),
			)
			return nil, diags
		}
	}

	instance := entity.NewCreateInstance(
		i.Region.ValueString(),
		i.Type.ValueString(),
		rootDiskStorageType,
		imageId,
		contractType,
		contractTerm,
		billingFrequency,
		entity.OptionalCreateInstanceValues{
			MarketAppId:  i.MarketAppId.ValueStringPointer(),
			Reference:    i.Reference.ValueStringPointer(),
			SshKey:       sshKey,
			RootDiskSize: rootDiskSize,
		},
	)

	return &instance, nil
}

func (i *Instance) GenerateUpdateInstanceEntity(ctx context.Context) (
	*entity.Instance,
	diag.Diagnostics,
) {
	diags := diag.Diagnostics{}

	id, err := value_object.NewUuid(i.Id.ValueString())
	if err != nil {
		diags.AddError(
			fmt.Sprintf(
				"invalid id %q",
				i.Id.ValueString(),
			),
			err.Error(),
		)
		return nil, diags
	}

	optionalValues := entity.OptionalUpdateInstanceValues{
		Type:      i.Type.ValueStringPointer(),
		Reference: i.Reference.ValueStringPointer(),
	}

	if i.RootDiskSize.ValueInt64() != 0 {
		rootDiskSize, err := value_object.NewRootDiskSize(int(i.RootDiskSize.ValueInt64()))
		if err != nil {
			diags.AddError(
				fmt.Sprintf(
					"invalid rootDiskSize %d",
					i.RootDiskSize.ValueInt64(),
				),
				err.Error(),
			)
			return nil, diags
		}
		optionalValues.RootDiskSize = rootDiskSize
	}

	contract := Contract{}
	diags = i.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	if contract.Type.ValueString() != "" {
		contractType, err := enum.NewContractType(contract.Type.ValueString())
		if err != nil {
			diags.AddError(
				fmt.Sprintf(
					"cannot convert parse contractType %s",
					contract.Type.ValueString(),
				),
				err.Error(),
			)
			return nil, diags
		}
		optionalValues.ContractType = &contractType
	}

	if contract.Term.ValueInt64() != 0 {
		contractTerm, err := enum.NewContractTerm(int(contract.Term.ValueInt64()))
		if err != nil {
			diags.AddError(
				fmt.Sprintf(
					"cannot convert parse contractTerm %d",
					contract.Term.ValueInt64(),
				),
				err.Error(),
			)
			return nil, diags
		}
		optionalValues.Term = &contractTerm
	}

	if contract.BillingFrequency.ValueInt64() != 0 {
		billingFrequency, err := enum.NewContractBillingFrequency(int(contract.BillingFrequency.ValueInt64()))
		if err != nil {
			diags.AddError(
				fmt.Sprintf(
					"cannot convert parse billingFrequency %d",
					contract.BillingFrequency.ValueInt64(),
				),
				err.Error(),
			)
			return nil, diags
		}
		optionalValues.BillingFrequency = &billingFrequency
	}

	instance := entity.NewUpdateInstance(*id, optionalValues)

	return &instance, nil
}
