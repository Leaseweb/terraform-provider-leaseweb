package to_resource_model

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/service/public_cloud/data_adapters/shared"
)

func AdaptInstanceDetails(sdkInstance publicCloud.InstanceDetails, ctx context.Context) (*resource.Instance, error) {
	instance := resource.Instance{
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
		resource.Image{}.AttributeTypes(),
		ctx,
		adaptImage,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Image = image

	ips, err := shared.AdaptSdkModelsToListValue(
		sdkInstance.Ips,
		resource.Ip{}.AttributeTypes(),
		ctx,
		adaptIpDetails,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Ips = ips

	contract, err := shared.AdaptSdkModelToResourceObject(
		sdkInstance.Contract,
		resource.Contract{}.AttributeTypes(),
		ctx,
		adaptContract,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

func AdaptInstance(sdkInstance publicCloud.Instance, ctx context.Context) (*resource.Instance, error) {
	instance := resource.Instance{
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
		resource.Image{}.AttributeTypes(),
		ctx,
		adaptImage,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Image = image

	ips, err := shared.AdaptSdkModelsToListValue(
		sdkInstance.Ips,
		resource.Ip{}.AttributeTypes(),
		ctx,
		adaptIp,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Ips = ips

	contract, err := shared.AdaptSdkModelToResourceObject(
		sdkInstance.Contract,
		resource.Contract{}.AttributeTypes(),
		ctx,
		adaptContract,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

func adaptImage(
	ctx context.Context,
	sdkImage publicCloud.Image,
) (*resource.Image, error) {
	return &resource.Image{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}, nil
}

func adaptIpDetails(
	ctx context.Context,
	sdkIp publicCloud.IpDetails,
) (*resource.Ip, error) {
	return &resource.Ip{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}, nil
}

func adaptIp(
	ctx context.Context,
	sdkIp publicCloud.Ip,
) (*resource.Ip, error) {
	return &resource.Ip{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}, nil
}

func adaptContract(
	ctx context.Context,
	sdkContract publicCloud.Contract,
) (*resource.Contract, error) {
	return &resource.Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(sdkContract.Term)),
		Type:             basetypes.NewStringValue(string(sdkContract.Type)),
		EndsAt:           shared.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.State)),
	}, nil
}
