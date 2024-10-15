package to_data_source_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
)

func AdaptInstances(sdkInstances []publicCloud.Instance) model.Instances {
	var instances model.Instances

	for _, sdkInstance := range sdkInstances {
		instance := adaptInstance(sdkInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func adaptInstance(sdkInstance publicCloud.Instance) model.Instance {
	return model.Instance{
		Id:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           shared.AdaptNullableStringToStringValue(sdkInstance.Reference.Get()),
		Image:               adaptImage(sdkInstance.Image),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		Ips:                 adaptIps(sdkInstance.Ips),
		Contract:            adaptContract(sdkInstance.Contract),
		MarketAppId:         shared.AdaptNullableStringToStringValue(sdkInstance.MarketAppId.Get()),
	}
}

func adaptImage(sdkImage publicCloud.Image) model.Image {
	return model.Image{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}
}

func adaptIps(sdkIps []publicCloud.Ip) []model.Ip {
	var ips []model.Ip
	for _, sdkIp := range sdkIps {
		ips = append(
			ips,
			adaptIp(sdkIp),
		)
	}

	return ips
}

func adaptIp(sdkIp publicCloud.Ip) model.Ip {
	return model.Ip{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}
}

func adaptContract(sdkContract publicCloud.Contract) model.Contract {
	return model.Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(sdkContract.Term)),
		Type:             basetypes.NewStringValue(string(sdkContract.Type)),
		EndsAt:           shared.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.State)),
	}
}
