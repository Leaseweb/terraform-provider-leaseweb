package to_data_source_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/datasource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/service/public_cloud/data_adapters/shared"
)

func AdaptInstances(sdkInstances []publicCloud.Instance) datasource.Instances {
	var instances datasource.Instances

	for _, sdkInstance := range sdkInstances {
		instance := adaptInstance(sdkInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func adaptInstance(sdkInstance publicCloud.Instance) datasource.Instance {
	return datasource.Instance{
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

func adaptImage(sdkImage publicCloud.Image) datasource.Image {
	return datasource.Image{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}
}

func adaptIps(sdkIps []publicCloud.Ip) []datasource.Ip {
	var ips []datasource.Ip
	for _, sdkIp := range sdkIps {
		ips = append(
			ips,
			adaptIp(sdkIp),
		)
	}

	return ips
}

func adaptIp(sdkIp publicCloud.Ip) datasource.Ip {
	return datasource.Ip{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}
}

func adaptContract(sdkContract publicCloud.Contract) datasource.Contract {
	return datasource.Contract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(sdkContract.Term)),
		Type:             basetypes.NewStringValue(string(sdkContract.Type)),
		EndsAt:           shared.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.State)),
	}
}
