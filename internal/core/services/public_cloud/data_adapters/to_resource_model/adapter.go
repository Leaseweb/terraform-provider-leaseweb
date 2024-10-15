package to_resource_model

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func AdaptInstanceDetails(sdkInstance publicCloud.InstanceDetails, ctx context.Context) (*model.Instance, error) {
	instance := model.Instance{
		Id:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           shared.AdaptNullableStringToStringValue(sdkInstance.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		MarketAppId:         shared.AdaptNullableStringToStringValue(sdkInstance.MarketAppId.Get()),
	}

	image, err := shared.AdaptDomainEntityToResourceObject(
		sdkInstance.Image,
		model.Image{}.AttributeTypes(),
		ctx,
		adaptImage,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Image = image

	ips, err := shared.AdaptEntitiesToListValue(
		sdkInstance.Ips,
		model.Ip{}.AttributeTypes(),
		ctx,
		adaptIpDetails,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Ips = ips

	contract, err := shared.AdaptDomainEntityToResourceObject(
		sdkInstance.Contract,
		model.Contract{}.AttributeTypes(),
		ctx,
		adaptContract,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

func AdaptInstance(sdkInstance publicCloud.Instance, ctx context.Context) (*model.Instance, error) {
	instance := model.Instance{
		Id:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           shared.AdaptNullableStringToStringValue(sdkInstance.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		MarketAppId:         shared.AdaptNullableStringToStringValue(sdkInstance.MarketAppId.Get()),
	}

	image, err := shared.AdaptDomainEntityToResourceObject(
		sdkInstance.Image,
		model.Image{}.AttributeTypes(),
		ctx,
		adaptImage,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Image = image

	ips, err := shared.AdaptEntitiesToListValue(
		sdkInstance.Ips,
		model.Ip{}.AttributeTypes(),
		ctx,
		adaptIp,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptInstanceDetails: %w", err)
	}
	instance.Ips = ips

	contract, err := shared.AdaptDomainEntityToResourceObject(
		sdkInstance.Contract,
		model.Contract{}.AttributeTypes(),
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
) (*model.Image, error) {
	return &model.Image{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}, nil
}

func adaptIpDetails(
	ctx context.Context,
	sdkIp publicCloud.IpDetails,
) (*model.Ip, error) {
	return &model.Ip{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}, nil
}

func adaptIp(
	ctx context.Context,
	sdkIp publicCloud.Ip,
) (*model.Ip, error) {
	return &model.Ip{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}, nil
}

func adaptContract(
	ctx context.Context,
	sdkContract publicCloud.Contract,
) (*model.Contract, error) {
	return &model.Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(sdkContract.Term)),
		Type:             basetypes.NewStringValue(string(sdkContract.Type)),
		EndsAt:           shared.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.State)),
	}, nil
}
