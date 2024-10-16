package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/dataadapters/shared"
)

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
	// TODO Enable SSH key support
	//SshKey              types.String `tfsdk:"ssh_key"`
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
		// TODO Enable SSH key support
		//"ssh_key": types.StringType,
	}
}

func NewFromInstance(
	sdkInstance publicCloud.Instance,
	ctx context.Context,
) (*Instance, error) {
	instance := Instance{
		Id:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           shared.AdaptNullableStringToStringValue(sdkInstance.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		MarketAppId:         shared.AdaptNullableStringToStringValue(sdkInstance.MarketAppId.Get()),
	}

	image, err := shared.AdaptSdkModelToResourceObject(
		sdkInstance.Image,
		Image{}.AttributeTypes(),
		ctx,
		newImage,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Image = image

	ips, err := shared.AdaptSdkModelsToListValue(
		sdkInstance.Ips,
		Ip{}.AttributeTypes(),
		ctx,
		newFromIp,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Ips = ips

	contract, err := shared.AdaptSdkModelToResourceObject(
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
		Reference:           shared.AdaptNullableStringToStringValue(sdkInstanceDetails.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstanceDetails.State)),
		Type:                basetypes.NewStringValue(string(sdkInstanceDetails.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstanceDetails.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstanceDetails.RootDiskStorageType)),
		MarketAppId:         shared.AdaptNullableStringToStringValue(sdkInstanceDetails.MarketAppId.Get()),
	}

	image, err := shared.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Image,
		Image{}.AttributeTypes(),
		ctx,
		newImage,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Image = image

	ips, err := shared.AdaptSdkModelsToListValue(
		sdkInstanceDetails.Ips,
		Ip{}.AttributeTypes(),
		ctx,
		newFromIpDetails,
	)
	if err != nil {
		return nil, fmt.Errorf("NewFromInstance: %w", err)
	}
	instance.Ips = ips

	contract, err := shared.AdaptSdkModelToResourceObject(
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
