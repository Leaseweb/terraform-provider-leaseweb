package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/model"
)

type ReasonInstanceCannotBeTerminated string

type Instance struct {
	Id                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Image               types.Object `tfsdk:"image"`
	State               types.String `tfsdk:"state"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int64  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	Ips                 types.List   `tfsdk:"ips"`
	Contract            types.Object `tfsdk:"contract"`
	MarketAppId         types.String `tfsdk:"market_app_id"`
}

func (i Instance) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"region":    types.StringType,
		"reference": types.StringType,
		"image": types.ObjectType{
			AttrTypes: Image{}.AttributeTypes(),
		},
		"state":                  types.StringType,
		"type":                   types.StringType,
		"root_disk_size":         types.Int64Type,
		"root_disk_storage_type": types.StringType,
		"ips": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: Ip{}.AttributeTypes(),
			},
		},
		"contract": types.ObjectType{
			AttrTypes: Contract{}.AttributeTypes(),
		},
		"market_app_id": types.StringType,
	}
}

func (i Instance) GetLaunchInstanceOpts(ctx context.Context) (
	*publicCloud.LaunchInstanceOpts,
	error,
) {
	sdkRootDiskStorageType, err := publicCloud.NewStorageTypeFromValue(
		i.RootDiskStorageType.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	image := Image{}
	imageDiags := i.Image.As(
		ctx,
		&image,
		basetypes.ObjectAsOptions{},
	)
	if imageDiags != nil {
		return nil, model.ReturnError(
			"AdaptToCreateInstanceOpts",
			imageDiags,
		)
	}

	contract := Contract{}
	contractDiags := i.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if contractDiags != nil {
		return nil, model.ReturnError(
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
		i.Region.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	sdkInstanceType, err := publicCloud.NewTypeNameFromValue(
		i.Type.ValueString(),
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

	opts.MarketAppId = model.AdaptStringPointerValueToNullableString(
		i.MarketAppId,
	)
	opts.Reference = model.AdaptStringPointerValueToNullableString(
		i.Reference,
	)
	opts.RootDiskSize = model.AdaptInt64PointerValueToNullableInt32(
		i.RootDiskSize,
	)

	return opts, nil
}

func (i Instance) GetUpdateInstanceOpts(ctx context.Context) (
	*publicCloud.UpdateInstanceOpts,
	error,
) {

	opts := publicCloud.NewUpdateInstanceOpts()
	opts.Reference = model.AdaptStringPointerValueToNullableString(
		i.Reference,
	)
	opts.RootDiskSize = model.AdaptInt64PointerValueToNullableInt32(
		i.RootDiskSize,
	)

	contract := Contract{}
	diags := i.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if diags.HasError() {
		return nil, model.ReturnError(
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

	if i.Type.ValueString() != "" {
		instanceType, err := publicCloud.NewTypeNameFromValue(
			i.Type.ValueString(),
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

func (i Instance) CanBeTerminated(ctx context.Context) (bool, *ReasonInstanceCannotBeTerminated, error) {
	contract := Contract{}
	contractDiags := i.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if contractDiags != nil {
		return false, nil, model.ReturnError(
			"AdaptToCreateInstanceOpts",
			contractDiags,
		)
	}

	if i.State.ValueString() == string(publicCloud.STATE_CREATING) || i.State.ValueString() == string(publicCloud.STATE_DESTROYING) || i.State.ValueString() == string(publicCloud.STATE_DESTROYED) {
		reason := ReasonInstanceCannotBeTerminated(
			fmt.Sprintf("state is %q", i.State),
		)

		return false, &reason, nil
	}

	if !contract.EndsAt.IsNull() {
		reason := ReasonInstanceCannotBeTerminated(
			fmt.Sprintf("contract.endsAt is %q", contract.EndsAt.ValueString()),
		)

		return false, &reason, nil
	}

	return true, nil, nil
}

func NewFromInstance(
	sdkInstance publicCloud.Instance,
	ctx context.Context,
) (*Instance, error) {
	instance := Instance{
		Id:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           model.AdaptNullableStringToStringValue(sdkInstance.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		MarketAppId:         model.AdaptNullableStringToStringValue(sdkInstance.MarketAppId.Get()),
	}

	image, err := model.AdaptSdkModelToResourceObject(
		sdkInstance.Image,
		Image{}.AttributeTypes(),
		ctx,
		newImage,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Image = image

	ips, err := model.AdaptSdkModelsToListValue(
		sdkInstance.Ips,
		Ip{}.AttributeTypes(),
		ctx,
		newFromIp,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Ips = ips

	contract, err := model.AdaptSdkModelToResourceObject(
		sdkInstance.Contract,
		Contract{}.AttributeTypes(),
		ctx,
		newContract,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

func NewFromInstanceDetails(
	sdkInstanceDetails publicCloud.InstanceDetails,
	ctx context.Context,
) (*Instance, error) {
	instance := Instance{
		Id:                  basetypes.NewStringValue(sdkInstanceDetails.Id),
		Region:              basetypes.NewStringValue(string(sdkInstanceDetails.Region)),
		Reference:           model.AdaptNullableStringToStringValue(sdkInstanceDetails.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstanceDetails.State)),
		Type:                basetypes.NewStringValue(string(sdkInstanceDetails.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstanceDetails.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstanceDetails.RootDiskStorageType)),
		MarketAppId:         model.AdaptNullableStringToStringValue(sdkInstanceDetails.MarketAppId.Get()),
	}

	image, err := model.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Image,
		Image{}.AttributeTypes(),
		ctx,
		newImage,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Image = image

	ips, err := model.AdaptSdkModelsToListValue(
		sdkInstanceDetails.Ips,
		Ip{}.AttributeTypes(),
		ctx,
		newFromIpDetails,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Ips = ips

	contract, err := model.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Contract,
		Contract{}.AttributeTypes(),
		ctx,
		newContract,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}
